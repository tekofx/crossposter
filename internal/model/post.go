package model

import (
	"time"
)

type Post struct {
	Id      int `gorm:"primaryKey";autoIncrement:true"`
	MediaId int // Documents in the same message have the same MediaId

	// Data
	Text   string
	Images []string `gorm:"type:text;serializer:json"`

	PublishedOnBsky     bool
	PublishedOnTelegram bool
	PublishedOnTwitter  bool

	// Meta
	CreatedAt time.Time `gorm:"type:DATE;"`
	HasText   bool
	HasImages bool
}
