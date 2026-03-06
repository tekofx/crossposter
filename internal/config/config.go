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
	BskyHandle      string
	BskyAppPassword string
	BskyPostHour    int

	// Telegram
	TelegramBotToken  string
	TelegramChannelId int
	TelegramOwner     int
	TelegramPostHour  int

	// Twitter
	TwitterUsername       string
	TwitterConsumerKey    string
	TwitterConsumerSecret string
	TwitterAccessToken    string
	TwitterAccessSecret   string
	TwitterPostHour       int

	// Instagram
	InstagramUserId      string
	InstagramAccessToken string
	InstagramPostHour    int
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
	if envVar == "" {
		logger.Fatal("Env variable %s required", name)
	}

	intValue, err := strconv.Atoi(envVar)
	if err != nil {
		logger.Fatal("Env variable %s must be integer", name)
	}

	return intValue
}

func getStringEnvVariable(name string) string {
	envVar := os.Getenv(name)
	if envVar == "" {
		logger.Fatal("Env variable %s required", name)
	}

	return envVar
}

func GetConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
	}

	return &Config{
		// Bluesky
		BskyHandle:      getStringEnvVariable("BSKY_HANDLE"),
		BskyAppPassword: getStringEnvVariable("BSKY_APP_PASSWORD"),
		BskyPostHour:    getIntEnvVariable("BSKY_POST_HOUR"),

		// Telegram
		TelegramBotToken:  getStringEnvVariable("TELEGRAM_BOT_TOKEN"),
		TelegramChannelId: getIntEnvVariable("TELEGRAM_CHANNEL_ID"),
		TelegramOwner:     getIntEnvVariable("TELEGRAM_OWNER"),
		TelegramPostHour:  getIntEnvVariable("TELEGRAM_POST_HOUR"),

		// Twitter
		TwitterUsername:       getStringEnvVariable("TWITTER_USERNAME"),
		TwitterConsumerKey:    getStringEnvVariable("TWITTER_CONSUMER_KEY"),
		TwitterConsumerSecret: getStringEnvVariable("TWITTER_CONSUMER_SECRET"),
		TwitterAccessToken:    getStringEnvVariable("TWITTER_ACCESS_TOKEN"),
		TwitterAccessSecret:   getStringEnvVariable("TWITTER_ACCESS_SECRET"),
		TwitterPostHour:       getIntEnvVariable("TWITTER_POST_HOUR"),

		// Instagram
		InstagramUserId:      getStringEnvVariable("INSTAGRAM_USER_ID"),
		InstagramAccessToken: getStringEnvVariable("INSTAGRAM_ACCESS_TOKEN"),
		InstagramPostHour:    getIntEnvVariable("INSTAGRAM_POST_HOUR"),
	}

}
