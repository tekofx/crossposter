package handlers

import (
	"regexp"
	"sync"
	"time"

	"github.com/mymmrac/telego"
	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/database"
	merrors "github.com/tekofx/crossposter/internal/errors"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/model"
	"github.com/tekofx/crossposter/internal/utils"

	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

var (
	groupMu sync.Mutex
	groups  = make(map[string][]telego.Message)
	timers  = make(map[string]*time.Timer)
)

func AddHandlers(bh *th.BotHandler, bot *telego.Bot) {
	onNewPrivateMessage(bh, bot)
	onKeyboardPress(bh)

}
func flushGroup(bot *telego.Bot, groupId string) {
	groupMu.Lock()
	messages := groups[groupId]
	delete(groups, groupId)
	delete(timers, groupId)
	groupMu.Unlock()

	// Process entire media group once
	err := processMediaGroup(bot, messages)
	if err != nil {
		logger.Error(err)
	}
}

func handleSingleMessage(bot *telego.Bot, msg telego.Message) {
	post := database.CreatePost()
	post.Text = msg.Text
	post.HasText = true
	utils.SendPostToOwner(bot, post)
	err := database.UpdatePost(post)
	if err != nil {
		logger.Error(err)
	}
}

func processMediaGroup(bot *telego.Bot, msgs []telego.Message) *merrors.MError {
	post := database.CreatePost()
	for _, m := range msgs {
		if m.Photo != nil {
			photoLen := len(m.Photo)
			file, err := utils.DownloadImage(bot, m.Photo[photoLen-1].FileID)
			if err != nil {
				return merrors.New(merrors.UnexpectedErrorCode, err.Error())
			}
			post.Images = append(post.Images,
				model.Image{
					Filename: *file,
					MimeType: "image/jpeg",
					FileSize: m.Photo[photoLen-1].FileSize,
				},
			)
			post.HasImages = true
		}
	}
	err := database.UpdatePost(post)
	if err != nil {
		return err
	}
	err = utils.SendPostToOwner(bot, post)
	if err != nil {
		return err
	}
	return nil
}
func onNewPrivateMessage(bh *th.BotHandler, bot *telego.Bot) {
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		if update.Message.Document != nil {
			utils.SendMessageToOwner(ctx, "Envía el archivo como imagen")
			return nil
		}
		if update.Message.MediaGroupID == "" {
			handleSingleMessage(bot, *update.Message)
			return nil
		}
		groupId := update.Message.MediaGroupID
		groupMu.Lock()

		// Append message to group
		groups[groupId] = append(groups[groupId], *update.Message)

		// Reset timer on each new message in group
		if timer, exists := timers[groupId]; exists {
			timer.Stop()
		}
		timers[groupId] = time.AfterFunc(500*time.Millisecond, func() {
			flushGroup(bot, groupId)
		})

		groupMu.Unlock()
		return nil

	}, utils.FromBotOwner())

}

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
