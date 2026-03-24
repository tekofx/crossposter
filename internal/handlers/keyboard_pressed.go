package handlers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"

	"github.com/mymmrac/telego"
	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/database"
	merrors "github.com/tekofx/crossposter/internal/errors"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/utils"
)

func onKeyboardPress(bh *th.BotHandler) {
	delRegex, _ := regexp.Compile(`delete:\d+`)
	editRegex, _ := regexp.Compile(`delete:\d+`)

	bh.HandleCallbackQuery(func(ctx *th.Context, query telego.CallbackQuery) error {
		if err := onEditPress(ctx, query); err != nil {
			logger.Error(err)
			if err.Code == merrors.NotFoundErrorCode {
				utils.SendMessageToOwner(ctx, "Ese post no existe")
			} else {
				utils.SendMessageToOwner(ctx, fmt.Sprintf("Error al editar post: %s", err.Message))
			}
		}
		return nil
	}, th.CallbackDataMatches(editRegex))

	bh.HandleCallbackQuery(func(ctx *th.Context, query telego.CallbackQuery) error {
		if err := onDeletePress(ctx, query); err != nil {
			logger.Error(err)
			if err.Code == merrors.NotFoundErrorCode {
				utils.SendMessageToOwner(ctx, "Ese post no existe")
			} else {
				utils.SendMessageToOwner(ctx, fmt.Sprintf("Error al borrar post: %s", err.Message))
			}
		}
		return nil
	}, th.CallbackDataMatches(delRegex))
}

func onDeletePress(ctx *th.Context, query telego.CallbackQuery) *merrors.MError {
	postId, err := strconv.Atoi(strings.Split(query.Data, ":")[1])
	if err != nil {
		ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
		return merrors.New(merrors.CannotConvertToIntErrorCode, err.Error())
	}
	merr := database.RemovePostById(postId)
	if merr != nil {
		ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
		return merr
	}
	ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(int64(config.Conf.TelegramOwner)), fmt.Sprintf("Eliminado post %d", postId)))
	ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
	return nil
}

func onEditPress(ctx *th.Context, query telego.CallbackQuery) *merrors.MError {
	// TODO: Implement
	postId, err := strconv.Atoi(strings.Split(query.Data, ":")[1])
	if err != nil {
		ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
		return merrors.New(merrors.CannotConvertToIntErrorCode, err.Error())
	}
	ctx.Bot().SendMessage(ctx, tu.Message(tu.ID(int64(config.Conf.TelegramOwner)), fmt.Sprintf("Editado post %d", postId)))
	ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(query.ID))
	return nil
}
