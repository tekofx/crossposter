package instagram

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/model"
)

func PostToInstagram(post *model.Post) {
	creationId, err := uploadImage(post)
	if err != nil {
		logger.Error("Instagram", "Error uploading image", err)
	}
	err = createPost(*creationId)
	if err != nil {
		logger.Error("Instagram", "Error creating post", err)
	}

}

const (
	instagramBaseURL = "https://graph.facebook.com/v22.0"
)

func uploadImage(post *model.Post) (*string, error) {
	containerURL := fmt.Sprintf("%s/%s/media", instagramBaseURL, config.Conf.InstagramUserId)
	data := url.Values{}
	data.Set("image_url", imageURL)
	data.Set("caption", post.Text)
	data.Set("access_token", config.Conf.InstagramAccessToken)

	resp, err := http.PostForm(containerURL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var containerResp map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&containerResp); err != nil {
		return nil, err
	}

	creationID, ok := containerResp["id"]
	if !ok {
		return nil, err
	}

	return &creationID, nil

}

func createPost(creationId string) error {
	data := url.Values{}
	publishURL := fmt.Sprintf("%s/%s/media_publish", instagramBaseURL, config.Conf.InstagramUserId)
	data = url.Values{}
	data.Set("creation_id", creationId)
	data.Set("access_token", config.Conf.InstagramAccessToken)

	resp, err := http.PostForm(publishURL, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var publishResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&publishResp); err != nil {
		return err
	}

	return nil

}
