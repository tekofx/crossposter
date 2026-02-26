package model

import (
	"fmt"
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

	BskyLink     string
	TelegramLink string
	TwitterLink  string

	// Meta
	CreatedAt time.Time `gorm:"type:DATE;"`
	HasText   bool
	HasImages bool
	Scheduled bool
}

func (post *Post) Message() string {
	format := func(ok bool, service string, url string) string {
		if ok {
			return fmt.Sprintf("✅ [%s](%s)", service, url)
		}
		return fmt.Sprintf("❌ %s", service)
	}

	return fmt.Sprintf("Resultado\n%s\n%s\n%s",
		format(post.PublishedOnBsky, "Bluesky", post.BskyLink),
		format(post.PublishedOnTelegram, "Telegram", post.TelegramLink),
		format(post.PublishedOnTwitter, "Twitter", post.TwitterLink),
	)
}
