package model

import (
	"time"

	"gorm.io/gorm"
)

type InstagramLogin struct {
	gorm.Model
	AccessToken string
	TokenType   string
	ExpireDate  time.Time `gorm:"type:DATETIME"`
}
