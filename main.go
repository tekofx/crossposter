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
	"github.com/tekofx/crossposter/internal/model"
	"github.com/tekofx/crossposter/internal/services/bsky"
	"github.com/tekofx/crossposter/internal/services/twitter"
	"github.com/tekofx/crossposter/internal/tasks"
)

func main() {
	config.Initialize()

	tasks.Initialize()
	database.InitializeDb()

	twitter.PostTweet("hola")
	return

	err := twitter.Initialize()
	if err != nil {
		logger.Fatal("Twitter", err)
	}

	_, err = twitter.PostToTwitter(&model.Post{Text: "Test"})
	if err != nil {
		logger.Error(err)
	}
	return

	err = bsky.Initialize()
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
