package bsky

import (
	"github.com/tekofx/crossposter/internal/config"
	merrors "github.com/tekofx/crossposter/internal/errors"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/model"
	"github.com/tekofx/crossposter/internal/services"
)

var BskyClient *BlueskyClient

const (
	createRecordUrl = "https://bsky.social/xrpc/com.atproto.repo.createRecord"
	uploadBlobUrl   = "https://bsky.social/xrpc/com.atproto.repo.uploadBlob"
	loginUrl        = "https://bsky.social/xrpc/com.atproto.server.createSession"
)

func InitializeBluesky() *merrors.MError {
	BskyClient = &BlueskyClient{Handle: config.Conf.BskyHandle, Password: config.Conf.BskyAppPassword}
	if err := authenticate(); err != nil {
		return err
	}
	return nil
}

func PostToBsky(post *model.Post) (*string, *merrors.MError) {
	var err *merrors.MError
	var postLink *string
	if post.HasImages {
		postLink, err = postImages(post)
	} else {
		postLink, err = postText(post)
	}

	if err != nil {
		return nil, err
	}

	post.BskyLink = *postLink
	post.PublishedOnBsky = true
	err = services.UpdatePost(post)
	if err != nil {
		logger.Error(err)
	}

	return postLink, err
}
