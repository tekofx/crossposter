package instagram

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/tekofx/crossposter/internal/config"
	merrors "github.com/tekofx/crossposter/internal/errors"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/model"
	"github.com/tekofx/crossposter/internal/services"
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
	services.StartServer()
	containerURL := fmt.Sprintf("%s/%s/media", instagramBaseURL, config.Conf.InstagramUserId)
	imageUrl := fmt.Sprintf("%s/%s", config.Conf.FileServerUrl, post.Images[0].Filename)
	imageUrl = "https://skyleriearts.tekofx.duckdns.org/AQADhgxrG3XFIFJ-.jpg"
	logger.Log(imageUrl)
	data := url.Values{}
	data.Set("image_url", imageUrl)
	data.Set("caption", post.Text)
	data.Set("access_token", config.Conf.InstagramAccessToken)

	resp, err := http.PostForm(containerURL, data)
	if err != nil {
		return nil, merrors.New(merrors.InvalidRequestErrorCode, err.Error())
	}
	defer resp.Body.Close()

	var containerResp map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&containerResp); err != nil {
		return nil, merrors.New(merrors.ParseJSONErrorCode, err.Error())
	}

	if errorObj, ok := containerResp["error"].(map[string]interface{}); ok {
		if code, ok := errorObj["code"].(float64); ok && int(code) == 190 {
			return nil, merrors.New(merrors.InstagramInvalidAccessTokenErrorCode, "Invalid Instagram Access Token")
		}
	}

	creationID, ok := containerResp["id"].(string)
	if !ok {
		return nil, merrors.New(merrors.InstagramUploadImageErrorCode, "Failed uploading image to Instagram")
	}

	services.StopServer()

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
