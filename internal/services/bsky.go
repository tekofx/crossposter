package services

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
)

type BlueskyClient struct {
	Handle   string
	Password string
	JWT      string
	DID      string
}

type BlobResponse struct {
	Blob Blob `json:"blob"`
}

type Post struct {
	Type      string       `json:"$type"`
	Text      string       `json:"text"`
	CreatedAt string       `json:"createdAt"`
	Embed     *EmbedImages `json:"embed,omitempty"`
}

type EmbedImages struct {
	Type   string      `json:"$type"`
	Images []ImageItem `json:"images"`
}

type ImageItem struct {
	Alt   string `json:"alt"`
	Image Blob   `json:"image"`
}

type Blob struct {
	Type     string `json:"$type"`
	Ref      Ref    `json:"ref"`
	MimeType string `json:"mimeType"`
	Size     int64  `json:"size"`
}

type Ref struct {
	Link string `json:"$link"`
}

type PostRequest struct {
	Repo       string     `json:"repo"`
	Collection string     `json:"collection"`
	Record     PostRecord `json:"record"`
}

// PostRecord represents the actual post content (the "record")
type PostRecord struct {
	Type      string       `json:"$type"`
	Text      string       `json:"text"`
	CreatedAt string       `json:"createdAt"`
	Embed     *EmbedImages `json:"embed,omitempty"` // Optional: only if embedding images
}

var BskyClient *BlueskyClient

func InitializeBluesky() error {
	BskyClient = &BlueskyClient{Handle: config.Conf.BskyHandle, Password: config.Conf.BskyAppPassword}
	if err := authenticate(); err != nil {
		return err
	}
	return nil
}

func PostToBsky(post *model.Post) error {
	var err error
	if post.HasImages {
		err = postImages(post)
	} else {
		err = postText(post.Text)
	}

	post.PublishedOnBsky = err == nil

	return err
}

func uploadBlob(image *model.Image) (*Blob, error) {
	url := "https://bsky.social/xrpc/com.atproto.repo.uploadBlob"

	file, err := os.ReadFile(image.Filename)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(file))
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

func postImages(post *model.Post) error {
	url := "https://bsky.social/xrpc/com.atproto.repo.createRecord"
	var uploadImages []ImageItem

	for _, image := range post.Images {
		blob, err := uploadBlob(&image)
		if err != nil {
			logger.Error("Error uploading blob", err)
			return err
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
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+BskyClient.JWT)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("post request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("post failed: %s %s", resp.Status, resp.Body)
	}

	return nil
}

func authenticate() error {
	loginUrl := "https://bsky.social/xrpc/com.atproto.server.createSession"
	loginPayload := map[string]string{
		"identifier": BskyClient.Handle,
		"password":   BskyClient.Password,
	}
	body, _ := json.Marshal(loginPayload)
	resp, err := http.Post(loginUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth failed: %s", resp.Status)
	}
	var loginData struct {
		AccessJwt string `json:"accessJwt"`
		Did       string `json:"did"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginData); err != nil {
		return fmt.Errorf("auth decode failed: %w", err)
	}

	BskyClient.JWT = loginData.AccessJwt
	BskyClient.DID = loginData.Did
	return nil
}

func postText(text string) error {
	postUrl := "https://bsky.social/xrpc/com.atproto.repo.createRecord"
	postPayload := PostRequest{
		Repo:       BskyClient.DID,
		Collection: "app.bsky.feed.post",
		Record: PostRecord{
			Type:      "app.bsky.feed.post",
			Text:      text,
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
		},
	}
	postBody, _ := json.Marshal(postPayload)
	req, _ := http.NewRequest("POST", postUrl, bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+BskyClient.JWT)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("post request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("post failed: %s", resp.Status)
	}
	return nil
}
