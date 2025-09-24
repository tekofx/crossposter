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
	BskyHandle       string
	TelegramBotToken string
	TelegramChatId   int
	PollInterval     time.Duration
	StateFile        string
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
		BskyHandle:       getStringEnvVariable("BSKY_HANDLE"),
		TelegramBotToken: getStringEnvVariable("TELEGRAM_BOT_TOKEN"),
		TelegramChatId:   getIntEnvVariable("TELEGRAM_CHAT_ID"),
		PollInterval:     60 * time.Second,
		StateFile:        "last_bsky_post.txt",
	}

}
