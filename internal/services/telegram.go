package services

import (
	"context"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/model"
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

func NotifyOwner(message string) error {
	_, err := bot.SendMessage(context.Background(), &telego.SendMessageParams{
		ChatID: tu.ID(int64(config.Conf.TelegramOwner)),
		Text:   message,
	})
	if err != nil {
		return err
	}

	return nil
}

func PostToTelegram(post model.BskyPost) error {
	if len(post.Post.Embed.Images) == 0 {
		_, err := bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID: tu.ID(int64(config.GetConfig().TelegramChannelId)),
			Text:   post.Post.Record.Text,
		})
		return err
	} else {
		err := postImages(post)
		return err
	}
}

func postImages(post model.BskyPost) error {
	var media []telego.InputMedia
	for i, image := range post.Post.Embed.Images {
		inputFile := telego.InputFile{
			URL: image.Fullsize,
		}

		if i == 0 {
			media = append(media, &telego.InputMediaPhoto{
				Type:    "photo",
				Media:   inputFile,
				Caption: post.Post.Record.Text,
			})
		} else {
			media = append(media, &telego.InputMediaPhoto{
				Type:  "photo",
				Media: inputFile,
			})
		}

	}

	_, err := bot.SendMediaGroup(
		context.Background(),
		&telego.SendMediaGroupParams{
			ChatID: tu.ID(int64(config.Conf.TelegramChannelId)),
			Media:  media,
		},
	)

	return err
}
