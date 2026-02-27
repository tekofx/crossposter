package telegram

import (
	"fmt"
	"strings"

	"github.com/mymmrac/telego"
)

func getTelegramPostLink(message telego.Message) string {
	if message.Chat.Username == "" {
		return fmt.Sprintf("https://t.me/c/%s/%d", strings.Split(message.Chat.ChatID().String(), "100")[1], message.MessageID)
	} else {
		return fmt.Sprintf("https://t.me/c/%s/%d", message.Chat.Username, message.MessageID)
	}
}
