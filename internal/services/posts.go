package services

import (
	"errors"
	"fmt"
	"log"

	"github.com/tekofx/crossposter/internal/database"
	"github.com/tekofx/crossposter/internal/model"
	"gorm.io/gorm"
)

func CreatePost() *model.Post {
	var post model.Post
	database.Database.Create(&post)
	return &post
}

func UpdatePost(post *model.Post) error {
	result := database.Database.Save(post)

	if result.Error != nil {
		fmt.Println("Error updating post:", result)
		return result.Error
	}

	return nil
}
func PostExistsInDatabase(bskyId string) bool {
	var post model.Post
	err := database.Database.
		Where("bsky_id = ?", bskyId).
		First(&post).
		Error
	return err == nil
}

func GetNewestPost() *model.Post {
	var post model.Post

	err := database.Database.
		Order("created_at DESC").
		Preload("Images"). // Load associated images (optional)
		First(&post).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // No posts in the database
		}
		// Log unexpected errors
		log.Printf("Database error fetching newest post: %v", err)
		return nil
	}

	return &post
}

func RemovePost(post *model.Post) error {
	return database.Database.Delete(post).Error
}
