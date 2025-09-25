package main

import (
	"fmt"
	"time"

	config "github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/database"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/model"
	"github.com/tekofx/crossposter/internal/services"
)

func main() {
	config.InitializeConfig()
	services.InitializeTelegram()
	database.InitializeDb()
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
		var newPosts []model.Post
		for _, post := range posts {

			if services.PostExistsInDatabase(post.BskyId) {
				break
			}

			if post.IsReply || post.IsRepost || post.IsQuote {
				continue
			}

			newPosts = append(newPosts, post)
		}
		if len(newPosts) == 0 {
			logger.Log("No new posts")
			logger.Log("Waiting", config.Conf.PollInterval)
			time.Sleep(config.Conf.PollInterval)
			continue
		}
		for _, post := range newPosts {
			logger.Log("Posting", post.BskyId)
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
			services.InsertPost(&post)
		}
		logger.Log("Waiting", config.Conf.PollInterval)
		time.Sleep(config.Conf.PollInterval)
	}
}
