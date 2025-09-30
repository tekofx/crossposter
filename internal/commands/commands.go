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
	postCommand(bh, bot)
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
		Scope: tu.ScopeChat(tu.ID(int64(config.Conf.TelegramOwner))),
	}
	bot.SetMyCommands(context.Background(), &PrivateChatCommands)

}

func deleteNewestPostCommand(bh *th.BotHandler) {
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		post := services.GetNewestPost()

		if post == nil {
			_, err := utils.SendMessageToOwner(ctx, post.Text)
			logger.Error(err)
			return err
		}

		err := services.RemovePost(post)
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
		post := services.GetNewestPost()
		if post == nil {
			utils.SendMessageToOwner(ctx, "No hay contenido en cola")
			return nil
		}
		_, err := utils.SendMessageToOwner(ctx, "Obteniendo post...")
		if err != nil {
			logger.Error(err)
			return err
		}
		if post.HasImages {
			err := utils.SendMediaGroupByFileIDs(bot, int64(config.Conf.TelegramOwner), post)
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

		commands, _ := bot.GetMyCommands(context.Background(), &telego.GetMyCommandsParams{})
		msg := `Proceso de publicación:
		1. Envia el texto o las imágenes (sin comprimir, y con texto opcional)
		2. Usa el comando /post para publicar`

		msg += "\nComandos\n"

		for _, command := range commands {
			msg += fmt.Sprintf("- /%s: %s\n", command.Command, command.Description)
		}
		_, err := utils.SendMessageToOwner(ctx, msg)
		if err != nil {
			logger.Fatal(err)
		}
		return nil
	}, th.CommandEqual("help"))
}

func postCommand(bh *th.BotHandler, bot *telego.Bot) {
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		post := services.GetNewestPost()
		if post == nil {
			_, err := utils.SendMessage(ctx, int64(config.Conf.TelegramOwner), "No se ha enviado contenido para publicar.")
			return err
		}
		utils.SendMessageToOwner(ctx, "Publicando post...")

		// bskyErr := services.PostToBsky(post)
		// if bskyErr != nil {
		// 	logger.Error(bskyErr)
		// }

		// tgErr := services.SendToChannel(bot, post)
		// if tgErr != nil {
		// 	logger.Error(tgErr)
		// }

		twitterErr := services.PostToTwitter(post)
		if twitterErr != nil {
			logger.Error(twitterErr)
		}

		_, err := utils.SendMessage(ctx, int64(config.Conf.TelegramOwner), post.Message())
		if err != nil {
			logger.Error(err)
			return err
		}

		// err = services.RemovePost(post)
		// if err != nil {
		// 	logger.Error(err)
		// 	return err
		// }
		return nil
	}, th.CommandEqual("post"))
}
