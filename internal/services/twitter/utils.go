package twitter

import (
	"context"
	"strings"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/media/upload"
	"github.com/michimani/gotwi/media/upload/types"
	"github.com/michimani/gotwi/tweet/managetweet"
	mtTypes "github.com/michimani/gotwi/tweet/managetweet/types"
	merrors "github.com/tekofx/crossposter/internal/errors"
)

func initializeMediaUpload(c *gotwi.Client, p *types.InitializeInput) (*types.InitializeOutput, *merrors.MError) {
	res, err := upload.Initialize(context.Background(), c, p)

	if err != nil {
		if strings.Contains(err.Error(), "503 Service Unavailable") {
			return nil, merrors.New(merrors.TwitterServiceUnavailableErrorCode, err.Error())
		}
		return nil, merrors.New(merrors.TwitterInitializeMediaErrorCode, err.Error())
	}

	return res, nil
}
func appendMediaUpload(c *gotwi.Client, p *types.AppendInput) (*types.AppendOutput, *merrors.MError) {
	res, err := upload.Append(context.Background(), c, p)
	if err != nil {
		return nil, merrors.New(merrors.TwitterCannotAppendMediaUploadErrorCode, err.Error())
	}

	return res, nil
}

func finalizeInput(c *gotwi.Client, p *types.FinalizeInput) (*types.FinalizeOutput, *merrors.MError) {
	res, err := upload.Finalize(context.Background(), c, p)
	if err != nil {
		return nil, merrors.New(merrors.TwitterCannotFinalizeInputErrorCode, err.Error())
	}

	return res, nil
}

func postTweetWithMedia(c *gotwi.Client, text string, mediaIds []string) (string, *merrors.MError) {
	p := &mtTypes.CreateInput{
		Text: gotwi.String(text),
		Media: &mtTypes.CreateInputMedia{
			MediaIDs: mediaIds,
		},
	}

	res, err := managetweet.Create(context.Background(), c, p)
	if err != nil {
		return "", merrors.New(merrors.TwitterCannotCreatePostErrorCode, err.Error())
	}

	return gotwi.StringValue(res.Data.ID), nil
}
