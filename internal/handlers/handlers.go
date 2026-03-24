package handlers

import (
	"github.com/mymmrac/telego"

	th "github.com/mymmrac/telego/telegohandler"
)

func AddHandlers(bh *th.BotHandler, bot *telego.Bot) {
	onNewPrivateMessage(bh, bot)
	onKeyboardPress(bh)
}
