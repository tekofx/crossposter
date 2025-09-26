package services

import (
	"fmt"

	"github.com/tekofx/crossposter/internal/database"
	"github.com/tekofx/crossposter/internal/model"
	"gorm.io/gorm/clause"
)

func InsertPost(post *model.Post) {
	database.Database.Create(post)
}

func InsertOrUpdatePost(post *model.Post) {
	// Replace "id" with your unique key, if different.
	err := database.Database.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}}, // or another unique field
			UpdateAll: true,                          // update all fields on conflict
		}).
		Create(post).Error

	if err != nil {
		// Handle error, e.g. log or return
		fmt.Println("Error inserting or updating post:", err)
	}
}
func PostExistsInDatabase(bskyId string) bool {
	var post model.Post
	err := database.Database.
		Where("bsky_id = ?", bskyId).
		First(&post).
		Error
	return err == nil
}

func GetNewestPost() (*model.Post, error) {
	var post model.Post
	err := database.Database.
		Order("created_at desc").
		First(&post).Error
	if err != nil {
		fmt.Println("asdfadswf")
		return nil, err
	}
	return &post, nil
}

func RemovePostByID(id uint) error {
	return database.Database.Delete(&model.Post{}, id).Error
}
