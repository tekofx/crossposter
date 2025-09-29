package services

import (
	"context"
	"fmt"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/tweet/managetweet"
	"github.com/michimani/gotwi/tweet/managetweet/types"
	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/model"
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

func PostToTwitter(post *model.Post) error {
	var err error

	if post.HasImages {

	} else {
		err = postTextToTwitter(post)
	}
	post.PublishedOnTwitter = err == nil
	return nil
}

func postTextToTwitter(post *model.Post) error {
	p := &types.CreateInput{
		Text: gotwi.String(post.Text),
	}
	res, err := managetweet.Create(context.Background(), twitterClient, p)
	if err != nil {
		logger.Error(err)
		return err
	}
	post.TwitterLink = fmt.Sprintf("https://x.com/%s/status/%s", config.Conf.TwitterUsername, *res.Data.ID)
	return nil
}
