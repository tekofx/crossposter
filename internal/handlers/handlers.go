package handlers

import (
	"fmt"

	"github.com/mymmrac/telego"
	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/model"
	"github.com/tekofx/crossposter/internal/services"
	"github.com/tekofx/crossposter/internal/utils"

	th "github.com/mymmrac/telego/telegohandler"
)

func AddHandlers(bh *th.BotHandler, bot *telego.Bot) {
	onNewPrivateMessage(bh, bot)

}

func onNewPrivateMessage(bh *th.BotHandler, bot *telego.Bot) {

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		if update.Message.Document != nil {
			utils.SendMessage(ctx, int64(config.Conf.TelegramOwner), "Envía el archivo como imagen")
			return nil
		}
		post := services.GetNewestPost()
		if post == nil {
			post = services.CreatePost()
		}

		if len(update.Message.Photo) > 0 {
			photoLen := len(update.Message.Photo)
			utils.SendMessageToOwner(ctx, "Recibida imagen")
			file, err := utils.DownloadImage(bot, update.Message.Photo[photoLen-1].FileID)
			if err != nil {
				logger.Error(err)
				return err
			}

			post.Images = append(post.Images,
				model.Image{
					Filename: *file,
					MimeType: "image/jpeg",
					FileSize: update.Message.Photo[photoLen-1].FileSize,
				},
			)
			post.HasImages = true
		} else {
			utils.SendMessageToOwner(ctx, fmt.Sprintf("Recibido texto %s", update.Message.Text))
			post.Text = update.Message.Text
			post.HasText = true
		}

		err := services.UpdatePost(post)
		if err != nil {
			return err
		}
		return nil

	}, utils.FromBotOwner())

}
