package services

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/tekofx/crossposter/internal/config"
)

var twitterClient *twitter.Client

func InitializeTwitter() {
	twitterConfig := oauth1.NewConfig(config.Conf.TwitterConsumerKey, config.Conf.TwitterConsumerSecret)
	token := oauth1.NewToken(config.Conf.TwitterAccessToken, config.Conf.TwitterAccessSecret)
	httpClient := twitterConfig.Client(oauth1.NoContext, token)
	twitterClient = twitter.NewClient(httpClient)
}

func PostToTwitter(text string) error {
	// X/Twitter limit: 280 chars
	if len(text) > 280 {
		text = text[:277] + "..."
	}
	_, _, err := twitterClient.Statuses.Update(text, nil)
	return err
}
