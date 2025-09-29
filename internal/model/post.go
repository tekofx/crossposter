package model

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	// Data
	Text   string
	Images []Image

	PublishedOnBsky     bool
	PublishedOnTelegram bool
	PublishedOnTwitter  bool

	// Meta
	CreatedAt time.Time `gorm:"type:DATE;"`
	HasText   bool
	HasImages bool
}
