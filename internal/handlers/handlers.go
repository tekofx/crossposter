package handlers

import (
	"github.com/mymmrac/telego"
	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/services"
	"github.com/tekofx/crossposter/internal/utils"

	th "github.com/mymmrac/telego/telegohandler"
)

func AddHandlers(bh *th.BotHandler) {
	onNewPrivateMessage(bh)

}

func onNewPrivateMessage(bh *th.BotHandler) {

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		if len(update.Message.Photo) > 0 {
			utils.SendMessage(ctx, int64(config.Conf.TelegramOwner), "Envía la imagen sin comprimir")
			return nil
		}

		if update.Message.Document != nil {
			utils.SendMessage(ctx, int64(config.Conf.TelegramOwner), "Received message with documents")
		} else {
			services.SendBskyTextPost(update.Message.Text)
			utils.SendMessage(ctx, int64(config.Conf.TelegramOwner), "Received text message")
		}

		return nil

	}, utils.FromBotOwner())

}
