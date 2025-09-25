package services

import (
	"github.com/tekofx/crossposter/internal/database"
	"github.com/tekofx/crossposter/internal/model"
)

func InsertPost(post *model.Post) {
	database.Database.Create(post)
}
func PostExistsInDatabase(bskyId string) bool {
	var post model.Post
	err := database.Database.
		Where("bsky_id = ?", bskyId).
		First(&post).
		Error
	return err == nil
}
