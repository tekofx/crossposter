package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/tekofx/crossposter/internal/config"
	merrors "github.com/tekofx/crossposter/internal/errors"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/model"
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
	helpCommand(bh)
	queueCommand(bh, bot)
	deletePostCommand(bh)
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

func deletePostCommand(bh *th.BotHandler) {
	bh.Handle(func(ctx *th.Context, update telego.Update) error {

		num, err := utils.GetIntArgument(update.Message.Text)
		if err != nil {
			logger.Error("Delete Command", err.Message)
			switch err.Code {
			case merrors.CannotConvertToIntErrorCode:
				utils.SendMessageToOwner(ctx, "El argumento no es un número válido")
			case merrors.TelegramArgumentNotProvidedErrorCode:
				utils.SendMessageToOwner(ctx, "Falta el id del post a eliminar. Usa /borrar id")
			}
			return nil
		}

		err = services.RemovePostById(*num)
		if err != nil {
			utils.SendMessageToOwner(ctx, "No se ha podido eliminar el post debido a un error")
			logger.Error("Delete Command", err)
			return nil
		}
		utils.SendMessageToOwner(ctx, "Post eliminado")
		return nil
	}, th.CommandEqual("borrar"))
}

func queueCommand(bh *th.BotHandler, bot *telego.Bot) {
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		posts, err := services.GetPosts()
		if err != nil {
			logger.Error("Queue command", err)
			return nil
		}
		if len(posts) == 0 {
			utils.SendMessageToOwner(ctx, "No hay contenido en cola")
			return nil
		}
		utils.SendMessageToOwner(ctx, "Obteniendo posts")

		for _, post := range posts {
			err := utils.SendPostToOwner(bot, &post)
			if err != nil {
				logger.Error("Queue command", err)
				return nil
			}
		}

		return nil
	}, th.CommandEqual("cola"))
}

func helpCommand(bh *th.BotHandler) {
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
		num, err := utils.GetIntArgument(update.Message.Text)
		if err != nil {
			logger.Error("Post Command", err.Message)
			switch err.Code {
			case merrors.CannotConvertToIntErrorCode:
				utils.SendMessageToOwner(ctx, "El argumento no es un número válido")
			case merrors.TelegramArgumentNotProvidedErrorCode:
				utils.SendMessageToOwner(ctx, "Falta el id del post. Usa /post id")
			}
			return nil
		}

		post, err := services.GetPostById(*num)
		if err != nil {
			logger.Error("Post command", err)
			utils.SendMessageToOwner(ctx, "Error al usar comando post")
			return nil
		}
		if post == nil {
			utils.SendMessageToOwner(ctx, "No se ha enviado contenido para publicar.")
			return nil
		}

		if post.Scheduled {
			utils.SendMessageToOwner(ctx, "El post ya está programado")
			return nil
		}

		utils.SendMessageToOwner(ctx, "Programando post...")

		go tasks.SchedulePost(model.Bluesky, bot, post, 20, 0)
		go tasks.SchedulePost(model.Instagram, bot, post, 20, 0)
		go tasks.SchedulePost(model.Telegram, bot, post, 20, 0)
		go tasks.SchedulePost(model.Twitter, bot, post, 20, 0)

		return nil
	}, th.CommandEqual("post"))
}
