package main

import (
	"context"
	"fmt"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/tekofx/crossposter/internal/commands"
	config "github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/database"
	"github.com/tekofx/crossposter/internal/handlers"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/services/twitter"
	"github.com/tekofx/crossposter/internal/tasks"
)

func main() {

	targetDate, duration := tasks.GetScheduledTime(12, 20)
	remainingHours, remamainingMinutes := tasks.GetDuration(duration)
	fmt.Println(remainingHours)
	fmt.Println(remamainingMinutes)

	fmt.Printf("Publicación en Telegram: %s (%d horas y %d minutos)", targetDate.Format("02-01-2006 15:04"), int(duration.Hours()), int(duration.Minutes()))

	return

	config.InitializeConfig()
	//services.InitializeTelegram()
	database.InitializeDb()

	// err := services.InitializeBluesky()
	// if err != nil {
	// 	logger.Fatal("Bluesky", err)
	// }
	err := twitter.InitializeTwitter()
	if err != nil {
		logger.Fatal(err)
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
	if err != nil {
		logger.Fatal(err)
	}

	// Add commands
	commands.AddCommands(bh, bot)
	handlers.AddHandlers(bh, bot)

	// Stop handling updates
	defer func() { _ = bh.Stop() }()
	logger.Log("Bot started as", bot.Username())
	tasks.CheckUnpostedPost(bot)
	botErr = bh.Start()
	if err != nil {
		logger.Fatal(err)
	}
}
