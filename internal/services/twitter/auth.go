package twitter

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/logger"
)

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

func auth(petitionParameters map[string]string) (*string, error) {
	// Endpoint
	apiURL := "https://api.x.com/2/tweets"

	parsedTime, err := time.Parse(http.TimeFormat, http.TimeFormat)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return nil, err
	}

	unixTimestamp := parsedTime.Unix()

	// OAuth 1.0a parameters
	params := map[string]string{
		"oauth_consumer_key":     config.Conf.TwitterConsumerKey,
		"oauth_nonce":            generateNonce(),
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        fmt.Sprintf("%d", unixTimestamp),
		"oauth_token":            config.Conf.TwitterAccessSecret,
		"oauth_version":          "1.0",
	}

	for k, v := range petitionParameters {
		params[k] = v
	}

	// Sort keys
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build parameter string
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, url.QueryEscape(params[k])))
	}
	paramStr := strings.Join(parts, "&")

	// Build base string
	baseStr := fmt.Sprintf("POST&%s&%s",
		url.QueryEscape(apiURL),
		url.QueryEscape(paramStr),
	)
	logger.Log("Basestr", baseStr)

	// Signing key
	signingKey := url.QueryEscape(config.Conf.TwitterConsumerKey) + "&" + url.QueryEscape(config.Conf.TwitterAccessSecret)
	logger.Log("Signing Key", signingKey)

	// Generate signature
	signature := sign("signingKey", baseStr)
	params["oauth_signature"] = signature

	logger.Log("Signature", signature)

	// Build Authorization header
	var authParts []string
	for k, v := range params {
		authParts = append(authParts, fmt.Sprintf(`%s="%s"`, k, url.QueryEscape(v)))
	}
	sort.Strings(authParts)
	authHeader := "OAuth " + strings.Join(authParts, ", ")
	logger.Log("OAuth", authHeader)

	return &authHeader, nil
}

func postTweet(text string) {
	// Endpoint
	apiURL := "https://api.x.com/2/tweets"

	petitionParams := map[string]string{
		"text": "test",
	}

	authHeader, err := auth(petitionParams)
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
