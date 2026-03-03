package tasks

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/mymmrac/telego"
	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/model"
	"github.com/tekofx/crossposter/internal/services"
	"github.com/tekofx/crossposter/internal/services/telegram"
	"github.com/tekofx/crossposter/internal/utils"
)

var tasksManager *TasksManager

func Initialize() {
	tasksManager = newTasksManager()
}

func StopTasksOfPost(postId int) {
	tasksManager.StopTask(fmt.Sprintf("%dBluesky", postId))
	tasksManager.StopTask(fmt.Sprintf("%dInstagram", postId))
	tasksManager.StopTask(fmt.Sprintf("%dTelegram", postId))
	tasksManager.StopTask(fmt.Sprintf("%dTwitter", postId))

}

func formatSchedule(text string, targetTime time.Time, duration time.Duration) string {
	return fmt.Sprintf("%s: %s (%d hours and %d minutes)", text, targetTime.Format("02-01-2006 15:04"), int(duration.Hours()), int(duration.Minutes())%60)
}

// Returns the time and the remaining time until a date
func getScheduledTime(hour int, minute int) (time.Time, time.Duration) {

	now := time.Now()
	target := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())

	if now.After(target) {
		target = target.Add(24 * time.Hour)
	}

	return target, target.Sub(now)
}

func SchedulePost(social model.SocialNetWork, bot *telego.Bot, post *model.Post, hour int, minute int) {

	taskId := strconv.Itoa(int(post.ID)) + social.String()
	tasksManager.StartTask(taskId, func(ctx context.Context) {
		post.Status = model.Scheduled
		services.UpdatePost(post)

		targetTime, duration := getScheduledTime(hour, minute)
		duration = time.Second * 15
		logger.Log("Task", taskId, social.String(), formatSchedule("Post Schedule", targetTime, duration))
		utils.SendMessageToOwnerUsingBot(bot, formatSchedule(social.String(), targetTime, duration))
		select {
		case <-ctx.Done():
			logger.Log("Task", taskId, "Stopped")
			return // Exit early if cancelled
		case <-time.After(duration):
			// Proceed after delay
		}

		postLink, tgErr := telegram.PostToTelegramChannel(bot, post)
		if tgErr != nil {
			logger.Error(social.String(), "Schedule", tgErr)
			return
		}

		_, err := utils.SendMessageToOwnerUsingBot(bot, fmt.Sprintf("Publicado en [%s](%s)", social.String(), *postLink))

		if err != nil {
			logger.Error(social.String(), "Scheduled Post", "Could not send post confirmation", err)
			return
		}
		checkToRemovePost(bot, post)
		return

	})

}

func GetAllTasks() string {
	return tasksManager.GetAllTasks()
}

// Checks if the post on database have been posted. If ncoot, schedule it
func CheckUnpostedPosts(bot *telego.Bot) {
	posts, err := services.GetPosts()
	if err != nil {
		logger.Error("CheckUnpostedPost", err)
		return
	}

	if len(posts) == 0 {
		return
	}

	for _, post := range posts {
		if post.Status != model.Scheduled {
			return
		}

		if !post.PublishedOnBsky {
			SchedulePost(model.Bluesky, bot, &post, config.Conf.InstagramPostHour, 0)
		}

		if !post.PublishedOnInstagram {
			SchedulePost(model.Instagram, bot, &post, config.Conf.InstagramPostHour, 0)
		}

		if !post.PublishedOnTelegram {
			SchedulePost(model.Telegram, bot, &post, config.Conf.TelegramPostHour, 0)
		}

		if !post.PublishedOnTwitter {
			SchedulePost(model.Twitter, bot, &post, config.Conf.TwitterPostHour, 0)
		}
	}

}

// If post have been posted to all socials, remove it from database
func checkToRemovePost(bot *telego.Bot, post *model.Post) {

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
