package config

import (
	"os"
	"strconv"
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

	// FileServer
	FileServerUrl  string
	FileServerPort int
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

func getIntEnvVariable(name string, required bool) int {
	envVar := os.Getenv(name)
	if envVar == "" && required {
		logger.Fatal("Env variable %s required", name)
	}

	intValue, err := strconv.Atoi(envVar)
	if err != nil {
		logger.Fatal("Env variable %s must be integer", name)
	}

	return intValue
}

func getStringEnvVariable(name string, required bool) string {
	envVar := os.Getenv(name)
	if envVar == "" && required {
		logger.Fatal("Env variable %s required", name)
	}

	return envVar
}

func GetConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
	}

	config := Config{
		// Bluesky
		BskyHandle:      getStringEnvVariable("BSKY_HANDLE", true),
		BskyAppPassword: getStringEnvVariable("BSKY_APP_PASSWORD", true),
		BskyPostHour:    getIntEnvVariable("BSKY_POST_HOUR", true),

		// Telegram Bot
		TelegramBotToken: getStringEnvVariable("TELEGRAM_BOT_TOKEN", false),

		// Telegram Channel
		TelegramChannelId: getIntEnvVariable("TELEGRAM_CHANNEL_ID", false),
		TelegramOwner:     getIntEnvVariable("TELEGRAM_OWNER", false),
		TelegramPostHour:  getIntEnvVariable("TELEGRAM_POST_HOUR", false),

		// Instagram
		InstagramUserId:           getStringEnvVariable("INSTAGRAM_USER_ID", false),
		InstagramClientId:         getIntEnvVariable("INSTAGRAM_CLIENT_ID", false),
		InstagramClientSecret:     getStringEnvVariable("INSTAGRAM_CLIENT_SECRET", false),
		InstagramPostHour:         getIntEnvVariable("INSTAGRAM_POST_HOUR", false),
		InstagramLoginRedirectUrl: getStringEnvVariable("INSTAGRAM_LOGIN_REDIRECT_URL", false),

		// FileServer
		FileServerUrl:  getStringEnvVariable("FILE_SERVER_URL", false),
		FileServerPort: getIntEnvVariable("FILE_SERVER_PORT", false),
	}

	if config.BskyHandle != "" && config.BskyAppPassword != "" {
		config.BskyEnabled = true
	}

	if config.TelegramBotToken != "" && config.TelegramChannelId != 0 && config.TelegramOwner != 0 {
		config.TelegramEnabled = true
	}

	if config.InstagramClientSecret != "" && config.InstagramUserId != "" && config.InstagramPostHour != 0 {
		config.InstagramEnabled = true
	}

	return &config

}
