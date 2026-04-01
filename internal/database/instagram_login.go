package database

import (
	"errors"
	"time"

	"github.com/tekofx/crossposter/internal/model"
	"gorm.io/gorm"
)

func CreateInstagramLogin(accessToken string, tokenType string, expiresIn int) *model.InstagramLogin {

	if ExistsInstagramLogin() {
		return nil
	}

	var instagramLogin model.InstagramLogin
	Database.Create(model.InstagramLogin{
		AccessToken: accessToken,
		TokenType:   tokenType,
		ExpireDate:  time.Now().Add(time.Duration(expiresIn) * time.Second),
	})
	return &instagramLogin
}

func GetInstagramLogin() *model.InstagramLogin {
	var post model.InstagramLogin
	result := Database.First(&post)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// No records found
		return nil
	}

	return (*model.InstagramLogin)(result.Statement.ReflectValue.UnsafePointer())
}

func ExistsInstagramLogin() bool {
	var count int64
	Database.Model(&model.InstagramLogin{}).Count(&count)
	return count > 0
}
