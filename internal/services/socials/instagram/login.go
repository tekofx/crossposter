package instagram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/tekofx/crossposter/internal/config"
	merrors "github.com/tekofx/crossposter/internal/errors"
)

type TokenData struct {
	AccessToken string `json:"access_token"`
	UserID      string `json:"user_id"`
	Permissions string `json:"permissions"`
}

type Response struct {
	Data []TokenData `json:"data"`
}

func GetLoginUrl() string {

	output := fmt.Sprintf(
		"https://www.instagram.com/oauth/authorize?client_id=%d&redirect_uri=%s&response_type=code&scope=instagram_business_basic,instagram_business_manage_messages,instagram_business_manage_comments,instagram_business_content_publish",
		config.Conf.InstagramClientId,
		config.Conf.InstagramLoginRedirectUrl,
	)
	return output

}

func GetTokenFromCode(code string) *merrors.MError {
	resp, err := http.PostForm("https://api.instagram.com/oauth/access_token", url.Values{
		"client_id":     {config.Conf.InstagramClientId},
		"client_secret": {config.Conf.InstagramClientSecret},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {config.Conf.InstagramLoginRedirectUrl},
		"code":          {code},
	})
	if err != nil {
		return merrors.New(merrors.InstagramInvalidGetTokenFromCodeErrorCode, err.Error())
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return merrors.New(merrors.ParseJSONErrorCode, err.Error())
	}

	return nil

}
