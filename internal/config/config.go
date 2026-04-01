package config

import (
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/tekofx/crossposter/internal/logger"
)

type Config struct {
	// Bluesky
	BskyEnabled     bool
	BskyHandle      string
	BskyAppPassword string
	BskyPostHour    int

	// Telegram
	TelegramEnabled   bool
	TelegramBotToken  string
	TelegramChannelId int
	TelegramOwner     int
	TelegramPostHour  int

	// Instagram
	InstagramEnabled          bool
	InstagramUserId           string
	InstagramClientId         int
	InstagramClientSecret     string
	InstagramLoginRedirectUrl string
	InstagramPostHour         int

	// WebServer
	WebServerUrl  string
	WebServerPort int
}

var lock = &sync.Mutex{}

var Conf *Config

func Initialize() {
	if Conf == nil {
		lock.Lock()
		defer lock.Unlock()
		if Conf == nil {
			Conf = GetConfig()
		}
	}

	err := os.MkdirAll("data", 0755)
	if err != nil {
		logger.Fatal(err)
	}

	err = os.MkdirAll("data/images", 0755)
	if err != nil {
		logger.Fatal(err)
	}
}

func getIntEnvVariable(name string) int {
	envVar := os.Getenv(name)

	intValue, err := strconv.Atoi(envVar)
	if err != nil {
		return -1
	}

	return intValue
}

func getStringEnvVariable(name string) string {
	return os.Getenv(name)
}

func GetConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
	}

	config := Config{

		// Bluesky
		BskyEnabled:     strings.ToLower(getStringEnvVariable("BSKY_ENABLED")) == "true",
		BskyHandle:      getStringEnvVariable("BSKY_HANDLE"),
		BskyAppPassword: getStringEnvVariable("BSKY_APP_PASSWORD"),
		BskyPostHour:    getIntEnvVariable("BSKY_POST_HOUR"),

		// Telegram Bot
		TelegramBotToken: getStringEnvVariable("TELEGRAM_BOT_TOKEN"),

		// Telegram Channel
		TelegramEnabled:   strings.ToLower(getStringEnvVariable("TELEGRAM_CHANNEL_ENABLED")) == "true",
		TelegramChannelId: getIntEnvVariable("TELEGRAM_CHANNEL_ID"),
		TelegramOwner:     getIntEnvVariable("TELEGRAM_OWNER"),
		TelegramPostHour:  getIntEnvVariable("TELEGRAM_POST_HOUR"),

		// Instagram
		InstagramEnabled:          strings.ToLower(getStringEnvVariable("INSTAGRAM_ENABLED")) == "true",
		InstagramUserId:           getStringEnvVariable("INSTAGRAM_USER_ID"),
		InstagramClientId:         getIntEnvVariable("INSTAGRAM_CLIENT_ID"),
		InstagramClientSecret:     getStringEnvVariable("INSTAGRAM_CLIENT_SECRET"),
		InstagramPostHour:         getIntEnvVariable("INSTAGRAM_POST_HOUR"),
		InstagramLoginRedirectUrl: getStringEnvVariable("INSTAGRAM_LOGIN_REDIRECT_URL"),

		// WebServer
		WebServerUrl:  getStringEnvVariable("WEB_SERVER_URL"),
		WebServerPort: getIntEnvVariable("WEB_SERVER_PORT"),
	}

	if config.TelegramBotToken == "" {
		logger.Fatal("Missing Telegram Bot Token")
	}

	if config.BskyEnabled && (config.BskyAppPassword == "" || config.BskyHandle == "" || config.BskyPostHour == -1) {
		logger.Fatal("Missing Bsky env vars")
	}

	if config.TelegramEnabled && (config.TelegramChannelId == -1 || config.TelegramOwner == -1 || config.TelegramPostHour == -1) {
		logger.Fatal("Missing Telegram env vars")
	}

	if config.InstagramEnabled &&
		(config.InstagramClientId == -1 || config.InstagramClientSecret == "" || config.InstagramLoginRedirectUrl == "" || config.InstagramPostHour == -1 || config.InstagramUserId == "") {
		logger.Fatal("Missing Instagram env vars")
	}

	return &config

}
