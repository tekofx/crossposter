package handlers

import (
	"github.com/mymmrac/telego"
	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/utils"

	th "github.com/mymmrac/telego/telegohandler"
)

func AddHandlers(bh *th.BotHandler) {
	onNewPrivateMessage(bh)

}

func onNewPrivateMessage(bh *th.BotHandler) {

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		utils.SendMessage(ctx, int64(config.Conf.TelegramOwner), "Received")
		return nil

	}, utils.FromBotOwner())

}
