package model

import (
	"time"

	"github.com/lib/pq"
)

type Post struct {
	Id int `gorm:"primaryKey";autoIncrement:true"`

	// Bsky
	BskyId      string
	IsQuote     bool
	IsRepost    bool
	IsReply     bool
	IsSelfQuote bool

	// Telegram
	TelegramId          int
	PublishedOnTelegram bool

	// Twitter
	TwitterUrl         string
	PublishedOnTwitter bool

	// Data
	Text   string
	Images pq.StringArray `gorm:"type:text[]"`

	// Meta
	Date time.Time `gorm:"type:DATE;"`
}
