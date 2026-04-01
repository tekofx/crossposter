package model

import (
	"gorm.io/gorm"
)

type InstagramLogin struct {
	gorm.Model
	AccessToken string
	TokenType   string
	ExpiresIn   int
}
