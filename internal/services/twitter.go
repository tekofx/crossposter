package services

import (
	"context"
	"fmt"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"
	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/logger"
)

var twitterClient *gotwi.Client

func InitializeTwitter() {
	in := &gotwi.NewClientInput{
		AuthenticationMethod: gotwi.AuthenMethodOAuth1UserContext,
		OAuthToken:           config.Conf.TwitterAccessToken,
		OAuthTokenSecret:     config.Conf.TwitterAccessSecret,
		APIKey:               config.Conf.TwitterConsumerKey,
		APIKeySecret:         config.Conf.TwitterConsumerSecret,
	}
	var err error
	twitterClient, err = gotwi.NewClient(in)
	if err != nil {
		logger.Fatal(err)
	}
}

func PostToTwitter(text string) error {
	p := &types.CreateInput{
		Text: gotwi.String("Uwu"),
	}
	res, err := managetweet.Create(context.Background(), twitterClient, p)
	if err != nil {
		logger.Error(err)
		return err
	}
	fmt.Printf("[%s] %s\n", gotwi.StringValue(res.Data.ID), gotwi.StringValue(res.Data.Text))

	return nil
}
