package services

import (
	"errors"
	"log"

	"github.com/tekofx/crossposter/internal/database"
	merrors "github.com/tekofx/crossposter/internal/errors"
	"github.com/tekofx/crossposter/internal/model"
	"gorm.io/gorm"
)

func CreatePost() *model.Post {
	var post model.Post
	database.Database.Create(&post)
	return &post
}

func UpdatePost(post *model.Post) *merrors.MError {
	result := database.Database.Save(post)

	if result.Error != nil {
		return merrors.New(merrors.UpdatePostErrorCode, result.Error.Error())
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

func RemovePost(post *model.Post) *merrors.MError {

	err := database.Database.Delete(post)
	if err.Error != nil {
		return merrors.New(merrors.RemovePostErrorCode, err.Error.Error())
	}

	return nil
}
