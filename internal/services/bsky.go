package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/model"
)

type BlueskyClient struct {
	Handle   string
	Password string
	JWT      string
	DID      string
}

var BskyClient *BlueskyClient

func InitializeBluesky() error {
	BskyClient = &BlueskyClient{Handle: config.Conf.BskyHandle, Password: config.Conf.BskyAppPassword}
	if err := authenticate(); err != nil {
		return err
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

func PostToBsky(post *model.Post) error {
	var err error
	if len(post.Images) > 0 {
	} else {
		err = postText(post.Text)
	}

	return err
}

// SendBskyTextPost sends a post to Bluesky with the given text
func postText(text string) error {
	postUrl := "https://bsky.social/xrpc/com.atproto.repo.createRecord"
	postPayload := map[string]interface{}{
		"repo":       BskyClient.DID,
		"collection": "app.bsky.feed.post",
		"record": map[string]interface{}{
			"$type":     "app.bsky.feed.post",
			"text":      text,
			"createdAt": time.Now().UTC().Format(time.RFC3339),
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
