package instagram

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/tekofx/crossposter/internal/config"
	merrors "github.com/tekofx/crossposter/internal/errors"
	"github.com/tekofx/crossposter/internal/model"
)

func PostToInstagram(post *model.Post) *merrors.MError {
	creationId, err := uploadImage(post)
	if err != nil {
		return err
	}
	err = createPost(*creationId)
	if err != nil {
		return err
	}

	return nil

}

const (
	instagramBaseURL = "https://graph.facebook.com/v22.0"
)

func uploadImage(post *model.Post) (*string, *merrors.MError) {
	containerURL := fmt.Sprintf("%s/%s/media", instagramBaseURL, config.Conf.InstagramUserId)
	data := url.Values{}
	data.Set("image_url", imageURL)
	data.Set("caption", post.Text)
	data.Set("access_token", config.Conf.InstagramAccessToken)

	resp, err := http.PostForm(containerURL, data)
	if err != nil {
		return nil, merrors.New(merrors.InvalidRequestErrorCode, err.Error())
	}
	defer resp.Body.Close()

	var containerResp map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&containerResp); err != nil {
		return nil, merrors.New(merrors.ParseJSONErrorCode, err.Error())
	}

	creationID, ok := containerResp["id"]
	if !ok {
		return nil, merrors.New(merrors.InstagramUploadImageErrorCode, "Failed uploading image to Instagram")
	}

	return &creationID, nil

}

func createPost(creationId string) *merrors.MError {
	data := url.Values{}
	publishURL := fmt.Sprintf("%s/%s/media_publish", instagramBaseURL, config.Conf.InstagramUserId)
	data = url.Values{}
	data.Set("creation_id", creationId)
	data.Set("access_token", config.Conf.InstagramAccessToken)

	resp, err := http.PostForm(publishURL, data)
	if err != nil {
		return merrors.New(merrors.InvalidRequestErrorCode, err.Error())
	}
	defer resp.Body.Close()

	var publishResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&publishResp); err != nil {
		return merrors.New(merrors.ParseJSONErrorCode, err.Error())
	}

	return nil

}
