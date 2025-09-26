package utils

import (
	"github.com/mymmrac/telego"

	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func SendMessage(ctx *th.Context, chatId int64, text string) *telego.Message {
	msg, _ := ctx.Bot().SendMessage(ctx, tu.Message(
		tu.ID(chatId),
		text,
	))
	return msg
}
