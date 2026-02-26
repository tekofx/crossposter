package services

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/media/upload"
	"github.com/michimani/gotwi/media/upload/types"
	mediaTypes "github.com/michimani/gotwi/media/upload/types"
	"github.com/michimani/gotwi/tweet/managetweet"
	mtTypes "github.com/michimani/gotwi/tweet/managetweet/types"
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

func PostToTwitter(post *model.Post) (*string, error) {
	var err error
	var postLink *string

	if post.HasImages {
		postLink, err = postImagesToTwitter(post)
	} else {
		postLink, err = postTextToTwitter(post)
	}

	if err != nil {
		return nil, err
	}

	post.TwitterLink = *postLink
	post.PublishedOnTwitter = true
	UpdatePost(post)
	return postLink, nil
}

func postTextToTwitter(post *model.Post) (*string, error) {
	p := &mtTypes.CreateInput{
		Text: gotwi.String(post.Text),
	}
	res, err := managetweet.Create(context.Background(), twitterClient, p)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	post.TwitterLink = fmt.Sprintf("https://x.com/%s/status/%s", config.Conf.TwitterUsername, *res.Data.ID)
	return &post.TwitterLink, nil
}

func postImagesToTwitter(post *model.Post) (*string, error) {

	var mediaIds []string
	for _, image := range post.Images {
		fileBytes, err := os.ReadFile(image.Filename)
		if err != nil {
			return nil, err
		}
		res, err := initialize(twitterClient, &mediaTypes.InitializeInput{
			MediaType:     mediaTypes.MediaType(image.MimeType),
			TotalBytes:    len(fileBytes),
			Shared:        false,
			MediaCategory: mediaTypes.MediaCategoryTweetImage,
		})
		mediaIds = append(mediaIds, res.Data.MediaID)
		_, err = appendMediaUpload(twitterClient, &types.AppendInput{
			MediaID:      res.Data.MediaID,
			Media:        bytes.NewReader(fileBytes),
			SegmentIndex: 0,
		})
		if err != nil {
			return nil, err
		}
		_, err = finalizeInput(twitterClient, &types.FinalizeInput{
			MediaID: res.Data.MediaID,
		})
		if err != nil {
			return nil, err
		}
	}
	postedID, err := postTweetWithMedia(twitterClient, post.Text, mediaIds)
	if err != nil {
		return nil, err
	}

	post.TwitterLink = fmt.Sprintf("https://x.com/%s/status/%s", config.Conf.TwitterUsername, postedID)
	return &post.TwitterLink, nil
}
func initialize(c *gotwi.Client, p *types.InitializeInput) (*types.InitializeOutput, error) {
	res, err := upload.Initialize(context.Background(), c, p)
	if err != nil {
		return nil, err
	}

	return res, nil
}
func appendMediaUpload(c *gotwi.Client, p *types.AppendInput) (*types.AppendOutput, error) {
	res, err := upload.Append(context.Background(), c, p)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func finalizeInput(c *gotwi.Client, p *types.FinalizeInput) (*types.FinalizeOutput, error) {
	res, err := upload.Finalize(context.Background(), c, p)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func postTweetWithMedia(c *gotwi.Client, text string, mediaIds []string) (string, error) {
	p := &mtTypes.CreateInput{
		Text: gotwi.String(text),
		Media: &mtTypes.CreateInputMedia{
			MediaIDs: mediaIds,
		},
	}

	res, err := managetweet.Create(context.Background(), c, p)
	if err != nil {
		return "", err
	}

	return gotwi.StringValue(res.Data.ID), nil
}
