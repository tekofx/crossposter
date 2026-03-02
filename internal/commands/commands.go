package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/services"
	"github.com/tekofx/crossposter/internal/tasks"
	"github.com/tekofx/crossposter/internal/utils"
)

var commands = []telego.BotCommand{
	{Command: "help", Description: "Mostrar mensaje de ayuda"},
	{Command: "post", Description: "Publica el post"},
	{Command: "cola", Description: "Mostrar post esperando para ser publicado"},
	{Command: "borrar", Description: "Elimina el post en cola"},
}

func AddCommands(bh *th.BotHandler, bot *telego.Bot) {
	postCommand(bh, bot)
	helpCommand(bh, bot)
	queueCommand(bh, bot)
	deleteNewestPostCommand(bh)
	startCommand(bh)

	var PrivateChatCommands = telego.SetMyCommandsParams{
		Commands: commands,
		Scope:    tu.ScopeChat(tu.ID(int64(config.Conf.TelegramOwner))),
	}
	bot.SetMyCommands(context.Background(), &PrivateChatCommands)

}

func startCommand(bh *th.BotHandler) {
	bh.Handle(func(ctx *th.Context, update telego.Update) error {

		utils.SendMessageToOwner(ctx, fmt.Sprintf("Hola %s! Usa /help para obtener info sobre mis comandos", update.Message.From.Username))
		return nil
	}, th.CommandEqual("start"))
}

func deleteNewestPostCommand(bh *th.BotHandler) {
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		post := services.GetNewestPost()

		if post == nil {
			utils.SendMessageToOwner(ctx, post.Text)
			return nil
		}

		err := services.RemovePost(post)
		if err == nil {
			utils.SendMessageToOwner(ctx, post.Text)
			return nil
		}
		utils.SendMessageToOwner(ctx, "Post eliminado")
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
		utils.SendMessageToOwner(ctx, "Obteniendo post")

		if post.HasImages {
			err := utils.SendMediaGroupByFileIDs(bot, int64(config.Conf.TelegramOwner), post)
			if err != nil {
				logger.Error(err)
				return nil
			}
		} else {
			utils.SendMessageToOwner(ctx, post.Text)
		}
		return nil
	}, th.CommandEqual("cola"))
}

func helpCommand(bh *th.BotHandler, bot *telego.Bot) {
	bh.Handle(func(ctx *th.Context, update telego.Update) error {

		var msg strings.Builder
		msg.WriteString(`Proceso de publicación:
		1. Envia el texto o las imágenes (con texto opcional)
		2. Usa el comando /post para publicar
		`)
		msg.WriteString("\nComandos\n")

		for _, command := range commands {
			fmt.Fprintf(&msg, "- /%s: %s\n", command.Command, command.Description)
		}

		utils.SendMessageToOwner(ctx, msg.String())

		return nil
	}, th.CommandEqual("help"))
}

func postCommand(bh *th.BotHandler, bot *telego.Bot) {
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		post := services.GetNewestPost()
		if post == nil {
			utils.SendMessageToOwner(ctx, "No se ha enviado contenido para publicar.")
			return nil
		}

		if post.Scheduled {
			utils.SendMessageToOwner(ctx, "El post ya está programado")
			return nil
		}

		utils.SendMessageToOwner(ctx, "Programando post...")

		go tasks.ScheduleToBsky(bot, post)
		go tasks.ScheduleToTelegram(bot, post)
		go tasks.ScheduleToTwitter(bot, post)
		return nil
	}, th.CommandEqual("post"))
}
