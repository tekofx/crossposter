package main

import (
	"context"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/tekofx/crossposter/internal/commands"
	config "github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/database"
	"github.com/tekofx/crossposter/internal/handlers"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/services"
)

func main() {
	config.InitializeConfig()
	//services.InitializeTelegram()
	database.InitializeDb()

	services.InitializeBluesky()
	services.InitializeTwitter()
	bot, botErr := telego.NewBot(config.Conf.TelegramBotToken)

	if botErr != nil {
		logger.Fatal(botErr)
	}

	// Get updates channel
	updates, err := bot.UpdatesViaLongPolling(context.Background(), nil)
	if err != nil {
		logger.Fatal(err)
	}

	// Create bot handler and specify from where to get updates
	bh, err := th.NewBotHandler(bot, updates)
	if err != nil {
		logger.Fatal(err)
	}

	// Add commands
	commands.AddCommands(bh, bot)
	handlers.AddHandlers(bh, bot)

	// Stop handling updates
	defer func() { _ = bh.Stop() }()
	logger.Log("Bot started as", bot.Username())
	err = bh.Start()
	if err != nil {
		logger.Fatal(err)
	}
}
