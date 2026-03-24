package handlers

import (
	"github.com/mymmrac/telego"
	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/database"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/model"
	"github.com/tekofx/crossposter/internal/utils"

	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func AddHandlers(bh *th.BotHandler, bot *telego.Bot) {
	onNewPrivateMessage(bh, bot)
	onKeyboardPress(bh)

}

func onNewPrivateMessage(bh *th.BotHandler, bot *telego.Bot) {

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		if update.Message.Document != nil {
			utils.SendMessageToOwner(ctx, "Envía el archivo como imagen")
			return nil
		}
		post := database.CreatePost()

		if len(update.Message.Photo) > 0 {
			photoLen := len(update.Message.Photo)
			file, err := utils.DownloadImage(bot, update.Message.Photo[photoLen-1].FileID)
			if err != nil {
				logger.Error(err)
				return err
			}

			post.Images = append(post.Images,
				model.Image{
					Filename: *file,
					MimeType: "image/jpeg",
					FileSize: update.Message.Photo[photoLen-1].FileSize,
				},
			)
			post.HasImages = true
		} else {
			post.Text = update.Message.Text
			post.HasText = true
		}

		utils.SendPostToOwner(ctx, post)

		err := database.UpdatePost(post)
		if err != nil {
			logger.Error(err)
		}
		return nil

	}, utils.FromBotOwner())

}

func onKeyboardPress(bh *th.BotHandler) {
	bh.HandleCallbackQuery(func(ctx *th.Context, query telego.CallbackQuery) error {
		ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(int64(config.Conf.TelegramOwner)), "Pulsaste Editar"))
		ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
		return nil
	}, th.CallbackDataEqual("edit"))
	bh.HandleCallbackQuery(func(ctx *th.Context, query telego.CallbackQuery) error {
		ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(int64(config.Conf.TelegramOwner)), "Pulsaste Borrar"))
		ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
		return nil
	}, th.CallbackDataEqual("delete"))

}
