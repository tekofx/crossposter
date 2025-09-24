package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	config "github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/model"
)

func getBlueskyFeed(handle string) (*model.BskyFeedResp, error) {
	// Bluesky doesn't have an official Go SDK, so we'll call the feed generator REST API
	url := fmt.Sprintf("https://public.api.bsky.app/xrpc/app.bsky.feed.getAuthorFeed?actor=%s&limit=5", handle)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var feed model.BskyFeedResp
	err = json.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}

	return &feed, nil
}

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

func postToTelegram(botToken string, text string) error {
	bot, botErr := telego.NewBot(config.Conf.TelegramBotToken)
	logger.Log("Logged in Telegram as", bot.Username())

	if botErr != nil {
		logger.Fatal(botErr)
	}

	_, err := bot.SendMessage(context.Background(), &telego.SendMessageParams{
		ChatID: tu.ID(int64(config.GetConfig().TelegramChatId)),
		Text:   text,
	})

	return err
}

func postToTwitter(client *twitter.Client, text string) error {
	// X/Twitter limit: 280 chars
	if len(text) > 280 {
		text = text[:277] + "..."
	}
	_, _, err := client.Statuses.Update(text, nil)
	return err
}

func main() {
	config.InitializeConfig()
	logger.Log("Started program")
	// Twitter client
	// config := oauth1.NewConfig(twitterConsumerKey, twitterConsumerSecret)
	// token := oauth1.NewToken(twitterAccessToken, twitterAccessSecret)
	// httpClient := config.Client(oauth1.NoContext, token)
	// twClient := twitter.NewClient(httpClient)

	for {
		fmt.Println("Checking bsky posts")
		feed, err := getBlueskyFeed(config.Conf.BskyHandle)
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
			newPosts = append(newPosts, post)
		}
		if len(newPosts) == 0 {
			logger.Log("No new posts")
			time.Sleep(config.Conf.PollInterval)
			continue
		}
		// Reverse to post oldest first
		for i, j := 0, len(newPosts)-1; i < j; i, j = i+1, j-1 {
			newPosts[i], newPosts[j] = newPosts[j], newPosts[i]
		}
		for _, post := range newPosts {
			logger.Log("Posting post", post.Post.Uri)
			err = postToTelegram(config.Conf.TelegramBotToken, post.Post.Record.Text)
			if err != nil {
				logger.Error(err)
			}
			//_ = postToTwitter(twClient, txt)
			setLastPostedURI(post.Post.Uri)
		}
		time.Sleep(config.Conf.PollInterval)
		fmt.Println("Waiting", config.Conf.PollInterval)
	}
}
