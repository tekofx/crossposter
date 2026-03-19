package telegram

import (
	"context"
	"os"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/database"
	merrors "github.com/tekofx/crossposter/internal/errors"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/model"
)

func PostToTelegramChannel(bot *telego.Bot, post *model.Post) (*string, *merrors.MError) {
	var postLink *string
	var err *merrors.MError
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
	err = database.UpdatePost(post)
	if err != nil {
		logger.Error(err)
	}

	return postLink, err
}

func postTgText(bot *telego.Bot, post *model.Post) (*string, *merrors.MError) {
	message, err := bot.SendMessage(context.Background(), tu.Message(tu.ID(int64(config.Conf.TelegramChannelId)), post.Text))
	if err != nil {
		return nil, merrors.New(merrors.TelegramCannotSendMessageToChannelErrorCode, err.Error())
	}

	postLink := getTelegramPostLink(*message)
	return &postLink, nil
}

func postTgImages(bot *telego.Bot, post *model.Post) (*string, *merrors.MError) {

	var photos []telego.InputMedia
	for i, image := range post.Images {
		file, err := os.Open(image.Filename)
		if err != nil {
			return nil, merrors.New(merrors.CannotReadFileErrorCode, err.Error())
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
		return nil, merrors.New(merrors.TelegramCannotSendMediaGroupErrorCode, err.Error())
	}
	postLink := getTelegramPostLink(messages[0])
	return &postLink, nil
}
