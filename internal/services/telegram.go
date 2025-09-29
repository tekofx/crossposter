package services

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/model"
)

func SendToChannel(bot *telego.Bot, post *model.Post) (*string, error) {

	var err error
	var message *telego.Message
	var messages []telego.Message
	if post.HasImages {
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

		messages, err = bot.SendMediaGroup(context.Background(), &telego.SendMediaGroupParams{
			ChatID: tu.ID(int64((config.Conf.TelegramChannelId))),
			Media:  photos,
		})

		message = &messages[0]

	} else {
		message, err = bot.SendMessage(context.Background(), tu.Message(tu.ID(int64(config.Conf.TelegramChannelId)), post.Text))
	}

	post.PublishedOnTelegram = err == nil

	var url string
	if message.Chat.Username == "" {
		url = fmt.Sprintf("https://t.me/c/%s/%d", strings.Split(message.Chat.ChatID().String(), "100")[1], message.MessageID)
	} else {
		url = fmt.Sprintf("https://t.me/c/%s/%d", message.Chat.Username, message.MessageID)
	}
	return &url, err
}
