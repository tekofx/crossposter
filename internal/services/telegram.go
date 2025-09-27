package services

import (
	"github.com/mymmrac/telego"
	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/logger"
)

var bot *telego.Bot

func InitializeTelegram() {
	var botErr error
	bot, botErr = telego.NewBot(config.Conf.TelegramBotToken)
	logger.Log("Logged in Telegram as", bot.Username())

	if botErr != nil {
		logger.Fatal(botErr)
	}
}
