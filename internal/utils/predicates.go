package utils

import (
	"context"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/tekofx/crossposter/internal/config"
)

func FromBotOwner() telegohandler.Predicate {
	return func(ctx context.Context, update telego.Update) bool {
		if update.Message == nil {
			return false
		}
		return update.Message.Chat.Type == "private" && update.Message.From.ID == int64(config.Conf.TelegramOwner)
	}
}
