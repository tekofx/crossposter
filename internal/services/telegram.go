package services

import (
	"context"
	"os"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/model"
)

func SendToChannel(bot *telego.Bot, post *model.Post) error {

	if post.HasImages {

		var photos []telego.InputMedia

		for i, image := range post.Images {
			file, err := os.Open(image)
			if err != nil {
				return err
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

		_, err := bot.SendMediaGroup(context.Background(), &telego.SendMediaGroupParams{
			ChatID: tu.ID(int64((config.Conf.TelegramChannelId))),
			Media:  photos,
		})
		return err

	} else {
		_, err := bot.SendMessage(context.Background(), tu.Message(tu.ID(int64(config.Conf.TelegramChannelId)), post.Text))
		if err != nil {
			return err
		}

	}

	return nil

}
