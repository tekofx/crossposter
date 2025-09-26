package utils

import (
	"github.com/mymmrac/telego"
	"github.com/tekofx/crossposter/internal/config"

	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func SendMessage(ctx *th.Context, chatId int64, text string) (*telego.Message, error) {
	msg, err := ctx.Bot().SendMessage(ctx, tu.Message(
		tu.ID(chatId),
		text,
	))

	if err != nil {
		return nil, err
	}
	return msg, nil
}

func SendMessageToOwner(ctx *th.Context, text string) (*telego.Message, error) {
	msg, err := ctx.Bot().SendMessage(ctx, tu.Message(
		tu.ID(int64(config.Conf.TelegramOwner)),
		text,
	))

	if err != nil {
		return nil, err
	}
	return msg, nil
}
