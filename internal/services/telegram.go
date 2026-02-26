package services

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/model"
)

func PostToTelegramChannel(bot *telego.Bot, post *model.Post) (*string, error) {
	var postLink *string
	var err error
	if post.HasImages {
		postLink, err = postTgImages(bot, post)
	} else {
		postLink, err = postTgText(bot, post)
	}

	if err != nil {
		return nil, err
	}

	post.TelegramLink = *postLink
	post.PublishedOnTelegram = true
	err = UpdatePost(post)
	if err != nil {
		logger.Error(err)
	}

	return postLink, err
}

func postTgText(bot *telego.Bot, post *model.Post) (*string, error) {
	message, err := bot.SendMessage(context.Background(), tu.Message(tu.ID(int64(config.Conf.TelegramChannelId)), post.Text))
	if err != nil {
		return nil, err
	}

	postLink := getTelegramPostLink(*message)
	return &postLink, err
}

func postTgImages(bot *telego.Bot, post *model.Post) (*string, error) {

	var photos []telego.InputMedia
	for i, image := range post.Images {
		file, err := os.Open(image.Filename)
		if err != nil {
			return nil, err
		}
		inputFile := telego.InputFile{
			File: file,
		}
		media := telego.InputMediaPhoto{
			Type:  "photo",
			Media: inputFile,
		}
		if post.Text != "" && i == 0 {
			media.Caption = post.Text
		}
		photos = append(photos, &media)
	}

	messages, err := bot.SendMediaGroup(context.Background(), &telego.SendMediaGroupParams{
		ChatID: tu.ID(int64((config.Conf.TelegramChannelId))),
		Media:  photos,
	})
	if err != nil {
		return nil, err
	}
	postLink := getTelegramPostLink(messages[0])
	return &postLink, nil
}

func getTelegramPostLink(message telego.Message) string {
	if message.Chat.Username == "" {
		return fmt.Sprintf("https://t.me/c/%s/%d", strings.Split(message.Chat.ChatID().String(), "100")[1], message.MessageID)
	} else {
		return fmt.Sprintf("https://t.me/c/%s/%d", message.Chat.Username, message.MessageID)
	}
}
