package instagram

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/tekofx/crossposter/internal/config"
)

func CreateLoginURL() string {

	output := fmt.Sprintf(
		"https://www.instagram.com/oauth/authorize?client_id=%d&redirect_uri=%s&response_type=code&scope=instagram_business_basic,instagram_business_manage_messages,instagram_business_manage_comments,instagram_business_content_publish",
		config.Conf.InstagramClientId,
		config.Conf.InstagramLoginRedirectUrl,
	)
	return output

}

func GetTokenFromCode() {
	resp, err := http.PostForm("https://api.instagram.com/oauth/access_token", url.Values{
		"client_id":     {"990602627938098"},
		"client_secret": {"a1b2C3D4"},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {"https://my.m.redirect.net/"},
		"code":          {"AQBx-hBsH3..."},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}
