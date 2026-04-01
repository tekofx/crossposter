package instagram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/database"
	merrors "github.com/tekofx/crossposter/internal/errors"
)

type TokenData struct {
	AccessToken string `json:"access_token"`
	UserId      string `json:"user_id"`
	Permissions int    `json:"permissions"`
}

type LongLivedTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type TokenResponse struct {
	Data []TokenData `json:"data"`
}

func GetLoginUrl() string {
	output := fmt.Sprintf(
		"https://www.instagram.com/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=instagram_business_basic,instagram_business_manage_messages,instagram_business_manage_comments,instagram_business_content_publish",
		config.Conf.InstagramClientId,
		config.Conf.InstagramLoginRedirectUrl,
	)
	return output
}

// Processes the code returned by Instagram API when an account authorizes the app
// 1. Requests token from code
// 2. With the token, requests a long-lived token and stores in db
// 3. Creates a task to renew the long-lived token
func processCode(code string) *merrors.MError {

	token, merr := requestTokenFromCode(code)
	if merr != nil {
		return merr
	}

	merr = requestLongLivedToken(*token)
	if merr != nil {
		return merr
	}
	return nil
}

func requestTokenFromCode(code string) (*string, *merrors.MError) {
	resp, err := http.PostForm("https://api.instagram.com/oauth/access_token", url.Values{
		"client_id":     {config.Conf.InstagramClientId},
		"client_secret": {config.Conf.InstagramClientSecret},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {config.Conf.InstagramLoginRedirectUrl},
		"code":          {code},
	})
	if err != nil {
		return nil, merrors.New(merrors.InstagramInvalidGetTokenFromCodeErrorCode, err.Error())
	}
	defer resp.Body.Close()

	var tokenResponse TokenResponse

	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return nil, merrors.New(merrors.ParseJSONErrorCode, err.Error())
	}

	return &tokenResponse.Data[0].AccessToken, nil

}

func requestLongLivedToken(token string) *merrors.MError {
	resp, err := http.PostForm("https://graph.instagram.com/access_token", url.Values{
		"grant_type":    {"ig_exchange_token"},
		"client_secret": {config.Conf.InstagramClientSecret},
		"access_token":  {token},
	})
	if err != nil {
		return merrors.New(merrors.InstagramInvalidGetTokenFromCodeErrorCode, err.Error())
	}
	defer resp.Body.Close()

	var longLivedTokenResponse LongLivedTokenResponse

	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &longLivedTokenResponse)
	if err != nil {
		return merrors.New(merrors.ParseJSONErrorCode, err.Error())
	}

	database.CreateInstagramLogin(
		longLivedTokenResponse.AccessToken,
		longLivedTokenResponse.TokenType,
		longLivedTokenResponse.ExpiresIn,
	)

	return nil
}
