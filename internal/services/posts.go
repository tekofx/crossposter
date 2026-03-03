package services

import (
	"errors"

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

func GetNewestPost() (*model.Post, *merrors.MError) {
	var post model.Post

	err := database.Database.
		Order("created_at DESC").
		Preload("Images"). // Load associated images (optional)
		First(&post).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No posts in the database
		}
		return nil, merrors.New(merrors.DatabaseErrorCode, err.Error())
	}

	return &post, nil
}

func GetPosts() ([]model.Post, *merrors.MError) {
	var posts []model.Post

	err := database.Database.
		Order("created_at DESC").
		Preload("Images"). // Load associated images (optional)
		Find(&posts).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return posts, nil // No posts in the database
		}
		return posts, merrors.New(merrors.DatabaseErrorCode, err.Error())
	}

	return posts, nil
}

func RemovePost(post *model.Post) *merrors.MError {
	err := database.Database.Delete(post)
	if err.Error != nil {
		return merrors.New(merrors.RemovePostErrorCode, err.Error.Error())
	}

	return nil
}

func RemovePostById(postId int) *merrors.MError {
	err := database.Database.Delete(&model.Post{}, postId)
	if err.Error != nil {
		return merrors.New(merrors.RemovePostErrorCode, err.Error.Error())
	}
	return nil
}

func GetPostById(postId string) (*model.Post, *merrors.MError) {
	var post model.Post
	err := database.Database.Find(&post, postId).Error
	if err != nil {
		return nil, merrors.New(merrors.DatabaseErrorCode, err.Error())
	}

	return &post, nil
}
