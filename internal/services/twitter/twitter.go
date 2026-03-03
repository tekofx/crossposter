package twitter

import (
	"bytes"
	"context"
	"fmt"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/media/upload/types"
	mediaTypes "github.com/michimani/gotwi/media/upload/types"
	"github.com/michimani/gotwi/tweet/managetweet"
	mtTypes "github.com/michimani/gotwi/tweet/managetweet/types"
	"github.com/tekofx/crossposter/internal/config"
	merrors "github.com/tekofx/crossposter/internal/errors"
	"github.com/tekofx/crossposter/internal/model"
	"github.com/tekofx/crossposter/internal/services"
	"github.com/tekofx/crossposter/internal/utils"
)

var twitterClient *gotwi.Client

func Initialize() *merrors.MError {
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
		return merrors.New(merrors.TwitterClientCreationErrorCode, err.Error())
	}
	return nil
}

func PostToTwitter(post *model.Post) (*string, *merrors.MError) {
	var err *merrors.MError
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
	services.UpdatePost(post)
	return postLink, nil
}

func postTextToTwitter(post *model.Post) (*string, *merrors.MError) {
	p := &mtTypes.CreateInput{
		Text: gotwi.String(post.Text),
	}
	res, err := managetweet.Create(context.Background(), twitterClient, p)
	if err != nil {
		return nil, merrors.New(merrors.TwitterCannotPostTextErrorCode, err.Error())
	}
	post.TwitterLink = fmt.Sprintf("https://x.com/%s/status/%s", config.Conf.TwitterUsername, *res.Data.ID)
	return &post.TwitterLink, nil
}

func postImagesToTwitter(post *model.Post) (*string, *merrors.MError) {

	var mediaIds []string
	for _, image := range post.Images {
		fileBytes, err := utils.ReadFile(image.Filename)
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
