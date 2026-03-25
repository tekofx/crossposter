package handlers

import (
	"sync"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/tekofx/crossposter/internal/database"
	merrors "github.com/tekofx/crossposter/internal/errors"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/model"
	"github.com/tekofx/crossposter/internal/utils"
)

var (
	groupMu sync.Mutex
	groups  = make(map[string][]telego.Message)
	timers  = make(map[string]*time.Timer)
)

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

func downloadPhotoAndUpdatePost(bot *telego.Bot, msg telego.Message, post *model.Post) *merrors.MError {
	photoLen := len(msg.Photo)
	file, err := utils.DownloadImage(bot, msg.Photo[photoLen-1].FileID)
	if err != nil {
		return err
	}
	post.Images = append(post.Images,
		model.Image{
			Filename: *file,
			MimeType: "image/jpeg",
			FileSize: msg.Photo[photoLen-1].FileSize,
		},
	)
	post.HasImages = true
	return nil
}

func processPhoto(bot *telego.Bot, msg telego.Message) *merrors.MError {
	post := database.CreatePost()
	err := downloadPhotoAndUpdatePost(bot, msg, post)
	if err != nil {
		return err
	}
	err = database.UpdatePost(post)
	if err != nil {
		return err
	}
	err = utils.SendPostToOwner(bot, post)
	if err != nil {
		return err
	}

	return nil
}

func processMediaGroup(bot *telego.Bot, msgs []telego.Message) *merrors.MError {
	post := database.CreatePost()
	for _, m := range msgs {
		err := downloadPhotoAndUpdatePost(bot, m, post)
		if err != nil {
			return err
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

		// If there's no media group
		if update.Message.MediaGroupID == "" {
			if len(update.Message.Photo) == 0 {
				handleSingleMessage(bot, *update.Message)
			} else {
				processPhoto(bot, *update.Message)
			}
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
