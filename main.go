package main

import (
	"context"
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

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/tekofx/crossposter/internal/commands"
	config "github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/database"
	"github.com/tekofx/crossposter/internal/handlers"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/services/bsky"
	"github.com/tekofx/crossposter/internal/services/twitter"
	"github.com/tekofx/crossposter/internal/tasks"
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

func postTweet(text string) error {
	// Endpoint
	apiURL := "https://api.x.com/2/tweets"

	// OAuth 1.0a parameters
	params := map[string]string{
		"oauth_consumer_key":     config.Conf.TwitterConsumerKey,
		"oauth_nonce":            generateNonce(),
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        fmt.Sprintf("%d", http.TimeFormat),
		"oauth_token":            config.Conf.TwitterAccessSecret,
		"oauth_version":          "1.0",
		"status":                 text,
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

	// JSON payload
	payload := fmt.Sprintf(`{"text":"%s"}`, text)

	// Make request
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 201 {
		return fmt.Errorf("failed to post tweet: %d %s", resp.StatusCode, string(body))
	}

	fmt.Printf("Tweet posted: %s\n", string(body))
	return nil
}
func main() {
	config.InitializeConfig()
	terr := postTweet("test")
	if terr != nil {
		logger.Error(terr)
	}

	return

	tasks.Initialize()
	database.InitializeDb()

	err := bsky.Initialize()
	if err != nil {
		logger.Fatal("Bluesky", err)
	}
	err = twitter.Initialize()
	if err != nil {
		logger.Fatal("Twitter", err)
	}
	bot, botErr := telego.NewBot(config.Conf.TelegramBotToken)

	if botErr != nil {
		logger.Fatal(botErr)
	}

	// Get updates channel
	updates, botErr := bot.UpdatesViaLongPolling(context.Background(), nil)
	if botErr != nil {
		logger.Fatal(botErr)
	}

	// Create bot handler and specify from where to get updates
	bh, botErr := th.NewBotHandler(bot, updates)
	if botErr != nil {
		logger.Fatal(err)
	}

	// Add commands
	commands.AddCommands(bh, bot)
	handlers.AddHandlers(bh, bot)

	// Stop handling updates
	defer func() { _ = bh.Stop() }()
	logger.Log("Bot started as", bot.Username())
	tasks.CheckUnpostedPosts(bot)
	botErr = bh.Start()
	if botErr != nil {
		logger.Fatal(err)
	}
}
