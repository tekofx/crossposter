package handlers

import (
	"fmt"

	"github.com/mymmrac/telego"
	"github.com/tekofx/crossposter/internal/config"
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
		if len(update.Message.Photo) > 0 {
			utils.SendMessage(ctx, int64(config.Conf.TelegramOwner), "Envía la imagen sin comprimir")
			return nil
		}
		post, _ := services.GetNewestPost()
		if post == nil {
			post = &model.Post{}
		}

		if update.Message.Document != nil {
			if !utils.IsImageExtension(update.Message.Document.FileName) {
				utils.SendMessageToOwner(ctx, "No se admite este archivo. Envía una imagen.")
				return nil
			}
			utils.SendMessageToOwner(ctx, fmt.Sprintf("Recibido archivo %s", update.Message.Document.FileName))
			post.Images = append(post.Images, update.Message.Document.FileID)
			post.HasImages = true
		} else {
			utils.SendMessageToOwner(ctx, fmt.Sprintf("Recibido texto %s", update.Message.Text))
			post.Text = update.Message.Text
			post.HasText = true
		}

		services.InsertOrUpdatePost(post)
		return nil

	}, utils.FromBotOwner())

}
