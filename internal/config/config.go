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

	// Telegram
	TelegramBotToken  string
	TelegramChannelId int
	TelegramOwner     int

	// Twitter
	TwitterUsername       string
	TwitterConsumerKey    string
	TwitterConsumerSecret string
	TwitterAccessToken    string
	TwitterAccessSecret   string

	// Instagram
	InstagramUserId      string
	InstagramAccessToken string
}

var lock = &sync.Mutex{}

var Conf *Config

func InitializeConfig() {
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

		// Telegram
		TelegramBotToken:  getStringEnvVariable("TELEGRAM_BOT_TOKEN"),
		TelegramChannelId: getIntEnvVariable("TELEGRAM_CHANNEL_ID"),
		TelegramOwner:     getIntEnvVariable("TELEGRAM_OWNER"),

		// Twitter
		TwitterUsername:       getStringEnvVariable("TWITTER_USERNAME"),
		TwitterConsumerKey:    getStringEnvVariable("TWITTER_CONSUMER_KEY"),
		TwitterConsumerSecret: getStringEnvVariable("TWITTER_CONSUMER_SECRET"),
		TwitterAccessToken:    getStringEnvVariable("TWITTER_ACCESS_TOKEN"),
		TwitterAccessSecret:   getStringEnvVariable("TWITTER_ACCESS_SECRET"),

		// Instagram
		InstagramUserId:      getStringEnvVariable("INSTAGRAM_USER_ID"),
		InstagramAccessToken: getStringEnvVariable("INSTAGRAM_ACCESS_TOKEN"),
	}

}
