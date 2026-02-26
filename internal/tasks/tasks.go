package tasks

import (
	"fmt"
	"time"

	"github.com/mymmrac/telego"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/model"
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

func ScheduleToTelegram(bot *telego.Bot, post *model.Post) {
	post.Scheduled = true
	services.UpdatePost(post)
	//waitUntilHour(16, 00)
	logger.Log("Telegram Post Scheduled")

	waitUntilHour(12, 58)

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
	checkToRemovePost(bot)

}

func ScheduleToBsky(bot *telego.Bot, post *model.Post) {
	logger.Log("Bsky Post Scheduled")

	post.Scheduled = true
	services.UpdatePost(post)
	//waitUntilHour(20, 00)
	waitUntilHour(13, 01)

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
	checkToRemovePost(bot)
}

func ScheduleToTwitter(bot *telego.Bot, post *model.Post) {
	logger.Log("Twitter Post Scheduled")

	post.Scheduled = true
	services.UpdatePost(post)
	//waitUntilHour(20, 00)
	waitUntilHour(13, 01)

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
	checkToRemovePost(bot)

}

// Checks if the post on database have been posted. If not, schedule it
func CheckUnpostedPost(bot *telego.Bot) {
	post := services.GetNewestPost()

	if post == nil {
		return
	}

	if !post.PublishedOnBsky {
		go ScheduleToBsky(bot, post)
	}

	if !post.PublishedOnTelegram {
		go ScheduleToTelegram(bot, post)
	}

	if !post.PublishedOnTwitter {
		go ScheduleToTwitter(bot, post)
	}

}

// If post have been posted to all socials, remove it from database
func checkToRemovePost(bot *telego.Bot) {
	post := services.GetNewestPost()

	if post.PublishedOnBsky && post.PublishedOnTelegram && post.PublishedOnTwitter {
		_, err := utils.SendMessageToOwnerUsingBot(bot, "Se ha publicado el post en todas las redes sociales. Eliminado de la cola.")
		if err != nil {
			logger.Error("checkToRemovePost", "Could not send message to owner", err)
		}
		err = services.RemovePost(post)
		if err != nil {
			logger.Error("checkToRemovePost", "Could not remove post from database", err)
		}
	}

}
