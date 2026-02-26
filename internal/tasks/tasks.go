package tasks

import (
	"fmt"
	"time"

	"github.com/mymmrac/telego"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/services"
	"github.com/tekofx/crossposter/internal/utils"
)

func waitUntilHour(hour int, minute int) {
	now := time.Now()
	target := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())

	if now.After(target) {
		target = target.Add(24 * time.Hour)
	}

	time.Sleep(time.Until(target))
}

func ScheduleToTelegram(bot *telego.Bot) {
	//waitUntilHour(16, 00)
	waitUntilHour(11, 56)

	post := services.GetNewestPost()
	postLink, tgErr := services.PostToTelegramChannel(bot, post)
	if tgErr != nil {
		logger.Error("Telegram", tgErr)
		return
	}

	_, err := utils.SendMessageToOwnerUsingBot(bot, fmt.Sprintf("Publicado en [Telegram](%s)", *postLink))

	if err != nil {
		logger.Error("Telegram Scheduled Post", "Could not send post confirmation", err)
		return
	}
}

func ScheduleToBsky(bot *telego.Bot) {
	//waitUntilHour(20, 00)
	waitUntilHour(11, 58)

	post := services.GetNewestPost()
	postLink, err := services.PostToBsky(post)
	if err != nil {
		logger.Error("Bluesky Scheduled Post", err)
		return
	}

	_, err = utils.SendMessageToOwnerUsingBot(bot, fmt.Sprintf("Publicado en [Bluesky](%s)", *postLink))

	if err != nil {
		logger.Error("Bluesky", "Could not send post confirmation", err)
		return
	}
}

func ScheduleToTwitter(bot *telego.Bot) {
	//waitUntilHour(20, 00)
	waitUntilHour(11, 58)

	post := services.GetNewestPost()
	postLink, err := services.PostToTwitter(post)
	if err != nil {
		logger.Error("Twitter Scheduled Post", err)
		return
	}

	_, err = utils.SendMessageToOwnerUsingBot(bot, fmt.Sprintf("Publicado en [Twitter](%s)", *postLink))

	if err != nil {
		logger.Error("Twitter", "Could not send post confirmation", err)
		return
	}
}
