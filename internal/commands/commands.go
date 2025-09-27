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
	postCommand(bh)
	helpCommand(bh, bot)
	queueCommand(bh, bot)
	deleteNewestPostCommand(bh)

	var PrivateChatCommands = telego.SetMyCommandsParams{
		Commands: []telego.BotCommand{
			{Command: "help", Description: "Mostrar mensaje de ayuda"},
			{Command: "post", Description: "Publica el post"},
			{Command: "cola", Description: "Mostrar post esperando para ser publicado"},
			{Command: "borrar", Description: "Elimina el post en cola"},
		},
		Scope: tu.ScopeAllPrivateChats(),
	}
	bot.SetMyCommands(context.Background(), &PrivateChatCommands)

}

func deleteNewestPostCommand(bh *th.BotHandler) {
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		post, err := services.GetNewestPost()

		if post == nil {
			_, err = utils.SendMessageToOwner(ctx, post.Text)
			logger.Error(err)
			return err
		}

		err = services.RemovePostByID(uint(post.Id))
		if post == nil {
			_, err = utils.SendMessageToOwner(ctx, post.Text)
			logger.Error(err)
			return err
		}

		_, err = utils.SendMessageToOwner(ctx, "Post eliminado")
		if err != nil {
			logger.Error(err)
			return err
		}
		return nil
	}, th.CommandEqual("borrar"))
}

func queueCommand(bh *th.BotHandler, bot *telego.Bot) {
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		_, err := utils.SendMessageToOwner(ctx, "Obteniendo post...")
		if err != nil {
			logger.Error(err)
			return err
		}
		post, err := services.GetNewestPost()
		if post == nil {
			utils.SendMessageToOwner(ctx, "No hay contenido en cola")
			return nil
		}

		if post.HasImages {
			err := utils.SendMediaGroupByFileIDs(bot, int64(config.Conf.TelegramOwner), post.Images, post.Text)
			if err != nil {
				logger.Error(err)
				return err
			}
		} else {
			_, err = utils.SendMessageToOwner(ctx, post.Text)
			if err != nil {
				logger.Error(err)
				return err
			}
		}
		return nil
	}, th.CommandEqual("cola"))
}

func helpCommand(bh *th.BotHandler, bot *telego.Bot) {
	bh.Handle(func(ctx *th.Context, update telego.Update) error {

		commands, _ := bot.GetMyCommands(context.Background(), &telego.GetMyCommandsParams{
			Scope: tu.ScopeAllPrivateChats(),
		})
		msg := `Proceso de publicación:
		1. Envia el texto o las imágenes (sin comprimir, y con texto opcional)
		2. Usa el comando /post para publicar.`

		msg += "\nComandos\n"

		for _, command := range commands {
			msg += fmt.Sprintf("- /%s: %s\n", command.Command, command.Description)
		}
		_, err := utils.SendMessage(ctx, int64(config.Conf.TelegramOwner), msg)
		if err != nil {
			logger.Fatal(err)
		}
		return nil
	}, th.CommandEqual("help"))
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
