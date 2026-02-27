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
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/model"
	"github.com/tekofx/crossposter/internal/utils"
)

func postImages(post *model.Post) (*string, error) {
	var uploadImages []ImageItem

	for _, image := range post.Images {
		blob, err := uploadBlob(&image)
		if err != nil {
			logger.Error("Error uploading blob", err)
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
		return nil, fmt.Errorf("post request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body once
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var publishReponse PublishResponse

	// If response is JSON, parse it
	if err := json.Unmarshal(body, &publishReponse); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("post failed: %s %s", resp.Status, resp.Body)
	}

	postLink := fmt.Sprintf("https://bsky.app/profile/%s/post/%s", config.Conf.BskyHandle, utils.LastSplit(publishReponse.Uri, "/"))

	return &postLink, nil
}

func postText(post *model.Post) (*string, error) {
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
		return nil, fmt.Errorf("post request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("post failed: %s", resp.Status)
	}
	// Read response body once
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var publishReponse PublishResponse

	// If response is JSON, parse it
	if err := json.Unmarshal(body, &publishReponse); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}
	post.BskyLink = fmt.Sprintf("https://bsky.app/profile/%s/post/%s", config.Conf.BskyHandle, utils.LastSplit(publishReponse.Uri, "/"))
	return &post.BskyLink, nil
}
func uploadBlob(image *model.Image) (*Blob, error) {

	file, err := os.ReadFile(image.Filename)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", uploadBlobUrl, bytes.NewReader(file))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+BskyClient.JWT)
	req.Header.Set("Content-Type", image.MimeType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upload failed with status: %d", resp.StatusCode)
	}

	var blobResp BlobResponse

	// If response is JSON, parse it
	if err := json.Unmarshal(body, &blobResp); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return &blobResp.Blob, nil
}
