package config

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/tekofx/crossposter/internal/logger"
)

type Config struct {
	// Bluesky
	BskyHandle string

	// Telegram
	TelegramBotToken string
	TelegramChatId   int

	// Twitter
	TwitterConsumerKey    string
	TwitterConsumerSecret string
	TwitterAccessToken    string
	TwitterAccessSecret   string

	// Other config
	PollInterval time.Duration
	StateFile    string
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
		BskyHandle: getStringEnvVariable("BSKY_HANDLE"),

		// Telegram
		TelegramBotToken: getStringEnvVariable("TELEGRAM_BOT_TOKEN"),
		TelegramChatId:   getIntEnvVariable("TELEGRAM_CHAT_ID"),

		// Twitter
		TwitterConsumerKey:    getStringEnvVariable("TWITTER_CONSUMER_KEY"),
		TwitterConsumerSecret: getStringEnvVariable("TWITTER_CONSUMER_SECRET"),
		TwitterAccessToken:    getStringEnvVariable("TWITTER_ACCESS_TOKEN"),
		TwitterAccessSecret:   getStringEnvVariable("TWITTER_ACCESS_SECRET"),

		// Other config
		PollInterval: 60 * time.Second,
		StateFile:    "last_bsky_post.txt",
	}

}
