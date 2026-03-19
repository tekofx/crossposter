package bsky

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/tekofx/crossposter/internal/config"
	merrors "github.com/tekofx/crossposter/internal/errors"
	"github.com/tekofx/crossposter/internal/model"
	"github.com/tekofx/crossposter/internal/utils"
)

func postImages(post *model.Post) (*string, *merrors.MError) {
	var uploadImages []ImageItem

	for _, image := range post.Images {
		blob, err := uploadBlob(&image)
		if err != nil {
			return nil, err
		}

		uploadImages = append(uploadImages, ImageItem{
			Alt:   "",
			Image: *blob})
	}

	embed := EmbedImages{
		Type:   "app.bsky.embed.images",
		Images: uploadImages,
	}

	postPayload := PostRequest{
		Repo:       BskyClient.DID,
		Collection: "app.bsky.feed.post",
		Record: PostRecord{
			Type:      "app.bsky.feed.post",
			Text:      post.Text,
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
			Embed:     &embed,
		},
	}

	postBody, _ := json.Marshal(postPayload)
	req, _ := http.NewRequest("POST", createRecordUrl, bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+BskyClient.JWT)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, merrors.New(merrors.BskyPostRequestErrorCode, err.Error())
	}
	defer resp.Body.Close()

	// Read response body once
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, merrors.New(merrors.ReadResponseBodyErrorCode, err.Error())
	}

	var publishReponse PublishResponse

	// If response is JSON, parse it
	if err := json.Unmarshal(body, &publishReponse); err != nil {
		return nil, merrors.New(merrors.ParseJSONErrorCode, err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, merrors.New(merrors.BskyPostErrorCode, fmt.Sprintf("%s", resp.Body))
	}

	postLink := fmt.Sprintf("https://bsky.app/profile/%s/post/%s", config.Conf.BskyHandle, utils.LastSplit(publishReponse.Uri, "/"))

	return &postLink, nil
}

func postText(post *model.Post) (*string, *merrors.MError) {
	postPayload := PostRequest{
		Repo:       BskyClient.DID,
		Collection: "app.bsky.feed.post",
		Record: PostRecord{
			Type:      "app.bsky.feed.post",
			Text:      post.Text,
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
		},
	}
	postBody, _ := json.Marshal(postPayload)
	req, _ := http.NewRequest("POST", createRecordUrl, bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+BskyClient.JWT)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, merrors.New(merrors.BskyPostRequestErrorCode, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, merrors.New(merrors.BskyPostErrorCode, fmt.Sprintf("%s", resp.Body))
	}
	// Read response body once
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, merrors.New(merrors.ReadResponseBodyErrorCode, err.Error())
	}

	var publishReponse PublishResponse

	// If response is JSON, parse it
	if err := json.Unmarshal(body, &publishReponse); err != nil {
		return nil, merrors.New(merrors.ParseJSONErrorCode, err.Error())
	}
	post.BskyLink = fmt.Sprintf("https://bsky.app/profile/%s/post/%s", config.Conf.BskyHandle, utils.LastSplit(publishReponse.Uri, "/"))
	return &post.BskyLink, nil
}
func uploadBlob(image *model.Image) (*Blob, *merrors.MError) {

	file, err := os.ReadFile(image.Filename)
	if err != nil {
		return nil, merrors.New(merrors.CannotReadFileErrorCode, err.Error())
	}
	req, err := http.NewRequest("POST", uploadBlobUrl, bytes.NewReader(file))
	if err != nil {
		return nil, merrors.New(merrors.CannotCreateRequestErrorCode, err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+BskyClient.JWT)
	req.Header.Set("Content-Type", image.MimeType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, merrors.New(merrors.DoRequestErrorCode, err.Error())
	}
	defer resp.Body.Close()

	// Read the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, merrors.New(merrors.ReadResponseBodyErrorCode, err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, merrors.New(merrors.BskyUploadBlobErrorCode, fmt.Sprintf("%s", resp.Body))
	}

	var blobResp BlobResponse

	// If response is JSON, parse it
	if err := json.Unmarshal(body, &blobResp); err != nil {
		return nil, merrors.New(merrors.ParseJSONErrorCode, err.Error())
	}

	return &blobResp.Blob, nil
}
