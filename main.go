package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	config "github.com/tekofx/crossposter/internal/config"
)

type BskyFeedResp struct {
	Feed []struct {
		Post struct {
			URI    string `json:"uri"`
			CID    string `json:"cid"`
			Record struct {
				Text string `json:"text"`
			} `json:"record"`
			CreatedAt string `json:"createdAt"`
		} `json:"post"`
	} `json:"feed"`
}

type BskyPost struct {
	Uri    string
	Author BskyAuthor
}

type BskyAuthor struct {
	did    string
	handle string
}

func getBlueskyFeed(handle string) ([]string, map[string]string, error) {
	// Bluesky doesn't have an official Go SDK, so we'll call the feed generator REST API
	url := fmt.Sprintf("https://public.api.bsky.app/xrpc/app.bsky.feed.getAuthorFeed?actor=%s&limit=5", handle)
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var feed BskyFeedResp
	err = json.Unmarshal(body, &feed)
	if err != nil {
		return nil, nil, err
	}
	var uris []string
	uriText := map[string]string{}
	for _, item := range feed.Feed {
		uris = append(uris, item.Post.URI)
		uriText[item.Post.URI] = item.Post.Record.Text
	}
	return uris, uriText, nil
}

func getLastPostedURI() string {
	data, err := ioutil.ReadFile(config.Conf.StateFile)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func setLastPostedURI(uri string) {
	ioutil.WriteFile(config.Conf.StateFile, []byte(uri), 0644)
}

func postToTelegram(botToken string, chatID int64, text string) error {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(chatID, text)
	_, err = bot.Send(msg)
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
	// Twitter client
	// config := oauth1.NewConfig(twitterConsumerKey, twitterConsumerSecret)
	// token := oauth1.NewToken(twitterAccessToken, twitterAccessSecret)
	// httpClient := config.Client(oauth1.NoContext, token)
	// twClient := twitter.NewClient(httpClient)

	for {
		fmt.Println("started program")
		uris, uriText, err := getBlueskyFeed(config.Conf.BskyHandle)
		fmt.Println(uris)
		if err != nil {
			log.Println("Bluesky error:", err)
			time.Sleep(config.Conf.PollInterval)
			continue
		}
		last := getLastPostedURI()
		var newUris []string
		for _, uri := range uris {
			if uri == last {
				break
			}
			newUris = append(newUris, uri)
		}
		if len(newUris) == 0 {
			time.Sleep(config.Conf.PollInterval)
			continue
		}
		// Reverse to post oldest first
		for i, j := 0, len(newUris)-1; i < j; i, j = i+1, j-1 {
			newUris[i], newUris[j] = newUris[j], newUris[i]
		}
		for _, uri := range newUris {
			txt := uriText[uri]
			log.Println("Reposting to Telegram and Twitter:", txt)
			//_ = postToTelegram(telegramBotToken, telegramChatID, txt)
			//_ = postToTwitter(twClient, txt)
			setLastPostedURI(uri)
		}
		time.Sleep(config.Conf.PollInterval)
		fmt.Println("Waiting", config.Conf.PollInterval)
	}
}
