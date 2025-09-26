package commands

import (
	"context"
	"fmt"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/services"
	"github.com/tekofx/crossposter/internal/utils"
)

func AddCommands(bh *th.BotHandler, bot *telego.Bot) {
	hi(bh)
	postCommand(bh)
	help(bh)

	var PrivateChatCommands = telego.SetMyCommandsParams{
		Commands: []telego.BotCommand{
			{Command: "hi", Description: "Hello"},
			{Command: "post", Description: "Publica el post"},
			{Command: "help", Description: "Mostrar mensaje de ayuda"},
		},
		Scope: tu.ScopeAllPrivateChats(),
	}
	bot.SetMyCommands(context.Background(), &PrivateChatCommands)
}

func help(bh *th.BotHandler) {
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		msg := "Proceso de publicación:\n1. Envia el texto o las imágenes (sin comprimir, y con texto opcional)\n2. Usa el comando /post para publicar."
		_, err := utils.SendMessage(ctx, int64(config.Conf.TelegramOwner), msg)
		if err != nil {
			logger.Fatal(err)
		}
		return nil
	}, th.CommandEqual("help"))
}

func hi(bh *th.BotHandler) {
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

func postCommand(bh *th.BotHandler) {
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		post, err := services.GetNewestPost()
		if post == nil {
			_, err = utils.SendMessage(ctx, int64(config.Conf.TelegramOwner), "No se ha enviado contenido para publicar.")
			return nil
		}
		reply := "Posteado\n"
		reply += "✅ Bsky\n"
		reply += "❌ Twitter\n"
		reply += "✅ Telegram\n"
		fmt.Println(post.Id)

		// bskyErr := services.PostToBsky(post)
		// if bskyErr == nil {
		// 	reply += "✅ Bsky"
		// }
		_, err = utils.SendMessage(ctx, int64(config.Conf.TelegramOwner), reply)
		if err != nil {
			logger.Error(err)
			return err
		}

		// err = services.RemovePostByID(uint(post.Id))
		// if err != nil {
		// 	logger.Error(err)
		// 	return err
		// }
		return nil
	}, th.CommandEqual("post"))
}
