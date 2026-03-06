package twitter

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/media/upload"
	"github.com/michimani/gotwi/media/upload/types"
	"github.com/michimani/gotwi/tweet/managetweet"
	mtTypes "github.com/michimani/gotwi/tweet/managetweet/types"
	merrors "github.com/tekofx/crossposter/internal/errors"
)

func formUrl(baseUrl string, pathParameters map[string]string) string {
	// Sort keys
	keys := make([]string, 0, len(pathParameters))
	for k := range pathParameters {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build parameter string
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, url.QueryEscape(pathParameters[k])))
	}
	paramStr := strings.Join(parts, "&")

	return fmt.Sprintf("%s?%s", baseUrl, paramStr)
}

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
