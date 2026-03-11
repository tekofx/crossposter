package twitter

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/tekofx/crossposter/internal/config"
	merrors "github.com/tekofx/crossposter/internal/errors"
	"github.com/tekofx/crossposter/internal/logger"
)

type SignatureData struct {
	HttpMethod          string
	Url                 string
	OauthConsumerKey    string
	OauthConsumerSecret string
	OauthNonce          string
	OauthTimestamp      int
	OauthToken          *string
	OauthVersion        string
	OtherData           map[string]string
}

func (s SignatureData) ToMap() map[string]string {
	data := map[string]string{
		"oauth_consumer_key":     s.OauthConsumerKey,
		"oauth_nonce":            s.OauthNonce,
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        string(s.OauthTimestamp),
		"oauth_version":          s.OauthVersion,
	}

	if s.OauthToken != nil {
		data["oauth_token"] = *s.OauthToken
	}

	for k, v := range s.OtherData {
		data[k] = v
	}

	return data

}

func generateNonce() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	rand.Read(b)
	for i := range b {
		b[i] = chars[b[i]%byte(len(chars))]
	}
	return string(b)
}

func sign(key, data string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func auth1a() *merrors.MError {

	requestTokenURL := "https://api.x.com/oauth/request_token"
	oauth_callback := "http://localhost:3000"
	params := map[string]string{
		"oauth_callback":     oauth_callback,
		"oauth_consumer_key": config.Conf.TwitterConsumerKey,
	}
	// Make request
	req, err := http.NewRequest(
		"POST",
		formUrl(requestTokenURL, params),
		nil,
	)
	logger.Log("Auth1a URL", req.URL)
	if err != nil {
		logger.Error(err)
		return merrors.New(merrors.DoRequestErrorCode, err.Error())
	}

	authHeader, err := createAuthHeader(requestTokenURL, "", params)
	req.Header.Set("Authorization", *authHeader)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return merrors.New(merrors.DoRequestErrorCode, err.Error())
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 201 {

		return merrors.New(merrors.UnexpectedErrorCode, string(body))
	}
	return nil
}

// Signing Key is composed of <consumer_secret>&<oauth_token_secret>. When <oauth_token_secret> is yet unknown,
// the signing key is only <consumer_secret>&
func formSigningKey(twitterConsumerSecret string, oauthTokenSecret *string) string {
	if oauthTokenSecret == nil {
		return fmt.Sprintf("%s&", twitterConsumerSecret)
	}

	return fmt.Sprintf("%s&%s", twitterConsumerSecret, *oauthTokenSecret)
}

func CreateSignature(signatureData SignatureData) string {
	// Build base string
	baseStr := url.QueryEscape(formUrl(signatureData.Url, signatureData.ToMap()))
	logger.Log("Basestr", baseStr)

	// Signing key
	signingKey := formSigningKey(signatureData.OauthConsumerSecret, nil)
	logger.Log("Signing Key", signingKey)

	// Generate signature
	signature := sign(signingKey, baseStr)

	return signature
}

func createAuthHeader(apiUrl string, twitterConsumerKey string, pathParameters map[string]string) (*string, error) {
	// Endpoint

	unixTimestamp := time.Now().Unix()

	// OAuth 1.0a parameters
	params := map[string]string{
		"oauth_consumer_key":     twitterConsumerKey,
		"oauth_nonce":            generateNonce(),
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        fmt.Sprintf("%d", unixTimestamp),
		"oauth_version":          "1.0",
	}

	maps.Copy(params, pathParameters)

	// Build base string
	baseStr := url.QueryEscape(formUrl(apiUrl, params))
	logger.Log("Basestr", baseStr)

	// Signing key
	signingKey := formSigningKey(twitterConsumerKey, nil)
	logger.Log("Signing Key", signingKey)

	// Generate signature
	signature := sign(signingKey, baseStr)
	params["oauth_signature"] = signature

	logger.Log("Signature", signature)

	// Build Authorization header
	var authParts []string
	for k, v := range params {
		authParts = append(authParts, fmt.Sprintf(`%s="%s"`, k, url.QueryEscape(v)))
	}
	sort.Strings(authParts)
	authHeader := "OAuth " + strings.Join(authParts, ", ")
	logger.Log(authHeader)

	return &authHeader, nil
}

func PostTweet(text string) {
	// Endpoint
	apiURL := "https://api.x.com/1.1/statuses/update.json?include_entities=true"

	petitionParams := map[string]string{
		"status": "test",
	}

	err2 := auth1a()
	if err2 != nil {
		logger.Fatal(err2)
	}

	return

	authHeader, err := createAuthHeader(apiURL, "", petitionParams)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Log(authHeader)

	// JSON payload
	payload := fmt.Sprintf(`{"text":"%s"}`, text)

	// Make request
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(payload))
	if err != nil {
		logger.Error(err)
		return
	}

	req.Header.Set("Authorization", *authHeader)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 201 {
		logger.Error(fmt.Errorf("failed to post tweet: %d %s", resp.StatusCode, string(body)))
		return
	}

	fmt.Printf("Tweet posted: %s\n", string(body))

}
