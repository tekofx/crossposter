package database

import (
	"errors"

	merrors "github.com/tekofx/crossposter/internal/errors"
	"github.com/tekofx/crossposter/internal/model"
	"gorm.io/gorm"
)

func CreatePost() *model.Post {
	var post model.Post
	Database.Create(&post)
	return &post
}

func UpdatePost(post *model.Post) *merrors.MError {
	result := Database.Save(post)

	if result.Error != nil {
		return merrors.New(merrors.UpdatePostErrorCode, result.Error.Error())
	}

	return nil
}
func PostExistsInDatabase(bskyId string) bool {
	var post model.Post
	err := Database.
		Where("bsky_id = ?", bskyId).
		First(&post).
		Error
	return err == nil
}

func GetPosts() ([]model.Post, *merrors.MError) {
	var posts []model.Post

	err := Database.
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
	err := Database.Delete(post)
	if err.Error != nil {
		return merrors.New(merrors.RemovePostErrorCode, err.Error.Error())
	}

	return nil
}

func RemovePostById(postId int) *merrors.MError {
	err := Database.Delete(&model.Post{}, postId)
	if err.Error != nil {
		return merrors.New(merrors.RemovePostErrorCode, err.Error.Error())
	}
	return nil
}

func GetPostById(postId int) (*model.Post, *merrors.MError) {
	var post model.Post
	result := Database.Preload("Images").First(&post, postId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, merrors.New(merrors.NotFoundErrorCode, "post not found")
		}
		return nil, merrors.New(merrors.DatabaseErrorCode, result.Error.Error())
	}
	return &post, nil
}
