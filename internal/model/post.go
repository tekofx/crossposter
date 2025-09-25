package model

import "time"

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
	Images []string

	// Meta
	Date time.Time `gorm:"type:DATE;"`
}
