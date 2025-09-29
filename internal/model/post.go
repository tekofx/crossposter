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

func (post *Post) Message() string {
	msg := "Resultado\n"
	if post.PublishedOnBsky {
		msg += "✅ Bluesky\n"
	} else {
		msg += "❌ Bluesky\n"
	}

	if post.PublishedOnTelegram {
		msg += "✅ Telegram Channel\n"
	} else {
		msg += "❌ Telegram Channel\n"
	}

	if post.PublishedOnTwitter {
		msg += "✅ Twitter\n"
	} else {
		msg += "❌ Twitter\n"
	}

	return msg
}
