package handlers

import (
	"regexp"

	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"

	"github.com/mymmrac/telego"
	"github.com/tekofx/crossposter/internal/config"
)

func onKeyboardPress(bh *th.BotHandler) {
	delRegex, _ := regexp.Compile(`delete:\d+`)
	bh.HandleCallbackQuery(func(ctx *th.Context, query telego.CallbackQuery) error {
		ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(int64(config.Conf.TelegramOwner)), "Pulsaste Editar"))
		ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
		return nil
	}, th.CallbackDataEqual("edit"))
	bh.HandleCallbackQuery(func(ctx *th.Context, query telego.CallbackQuery) error {
		ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(int64(config.Conf.TelegramOwner)), query.Data))
		ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
		return nil
	}, th.CallbackDataMatches(delRegex))

}
