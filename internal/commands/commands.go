package commands

import (
	"context"
	"fmt"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/tekofx/crossposter/internal/logger"
)

func AddCommands(bh *th.BotHandler, bot *telego.Bot) {
	hi(bh)

	var PrivateChatCommands = telego.SetMyCommandsParams{
		Commands: []telego.BotCommand{
			{Command: "hi", Description: "Hello"},
		},
		Scope: tu.ScopeAllPrivateChats(),
	}
	bot.SetMyCommands(context.Background(), &PrivateChatCommands)
}

func hi(bh *th.BotHandler) {
	fmt.Println("asdf")
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		_, err := ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(update.Message.Chat.ID),
			fmt.Sprintf("Hello %s!", update.Message.From.FirstName),
		))
		if err != nil {
			logger.Fatal(err)
		}
		return nil
	}, th.CommandEqual("hi"))
}
