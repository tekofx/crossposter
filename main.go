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

	logger.Log("Started program")
	// Twitter client
	// config := oauth1.NewConfig(twitterConsumerKey, twitterConsumerSecret)
	// token := oauth1.NewToken(twitterAccessToken, twitterAccessSecret)
	// httpClient := config.Client(oauth1.NoContext, token)
	// twClient := twitter.NewClient(httpClient)

	for {
		fmt.Println("Checking bsky posts")
		feed, err := services.GetBlueskyFeed()
		if err != nil {
			logger.Error("Bluesky error:", err)
			time.Sleep(config.Conf.PollInterval)
			continue
		}
		last := getLastPostedURI()
		var newPosts []model.BskyPost
		for _, post := range feed.Posts {
			if post.Post.Uri == last {
				break
			}
			if post.Reason != nil {
				break
			}

			newPosts = append(newPosts, post)
		}
		if len(newPosts) == 0 {
			logger.Log("No new posts")
			time.Sleep(config.Conf.PollInterval)
			continue
		}
		for _, post := range newPosts {
			logger.Log("Posting post", post.Post.Uri)
			err = services.PostToTelegram(post)
			if err != nil {
				logger.Error(err)
			}
			//_ = postToTwitter(twClient, txt)
			setLastPostedURI(post.Post.Uri)
		}
		fmt.Println("Waiting", config.Conf.PollInterval)
		time.Sleep(config.Conf.PollInterval)
	}
}
