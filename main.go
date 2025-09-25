package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	config "github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/model"
	"github.com/tekofx/crossposter/internal/services"
)

func getLastPostedURI() string {
	data, err := os.ReadFile(config.Conf.StateFile)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func setLastPostedURI(uri string) {
	os.WriteFile(config.Conf.StateFile, []byte(uri), 0644)
}

func main() {
	config.InitializeConfig()
	services.InitializeTelegram()
	//database.InitializeDb()
	//services.InitializeTwitter()

	logger.Log("Started program")

	for {
		logger.Log("Checking bsky posts")
		posts, err := services.GetBlueskyPosts()
		if err != nil {
			logger.Error("Bluesky error:", err)
			time.Sleep(config.Conf.PollInterval)
			continue
		}
		last := getLastPostedURI()
		var newPosts []model.Post
		for _, post := range posts {

			if post.BskyId == last {
				continue
			}

			if post.IsReply || post.IsRepost || post.IsQuote {
				continue
			}

			newPosts = append(newPosts, post)
		}
		if len(newPosts) == 0 {
			logger.Log("No new posts")
			time.Sleep(config.Conf.PollInterval)
			continue
		}
		for _, post := range newPosts {

			err = services.PostToTelegram(&post)
			if err != nil {
				logger.Error(err)
				services.NotifyOwner(fmt.Sprintf("Could not post to telegram channel: %s", err))
			}
			post.PublishedOnTelegram = true
			// err := services.PostToTwitter("test")
			// if err != nil {
			// 	logger.Error(err)
			// 	services.NotifyOwner(fmt.Sprintf("Could not post to Twitter account: %s", err))
			// }

			//setLastPostedURI(post.BskyId)
		}
		logger.Log("Waiting", config.Conf.PollInterval)
		time.Sleep(config.Conf.PollInterval)
	}
}
