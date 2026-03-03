package model

import (
	"fmt"
	"time"

	"github.com/tekofx/crossposter/internal/types"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model

	// Data
	Text   string
	Images []Image

	PublishedOnBsky      bool
	PublishedOnInstagram bool
	PublishedOnTelegram  bool
	PublishedOnTwitter   bool

	BskyLink      string
	InstagramLink string
	TelegramLink  string
	TwitterLink   string

	// Meta
	CreatedAt time.Time `gorm:"type:DATE;"`
	HasText   bool
	HasImages bool
	Status    types.PostStatus
}

func (post *Post) String() string {
	format := func(ok bool, service string, url string) string {
		if ok {
			return fmt.Sprintf("✅ [%s](%s)", service, url)
		}
		return fmt.Sprintf("❌ %s", service)
	}

	var msg string

	if post.HasText {
		msg += fmt.Sprintf("%s\n", post.Text)
	}

	msg += fmt.Sprintf("ID: %d\nEstado: %s\n%s\n%s\n%s\n%s",
		post.ID,
		post.Status.String(),
		format(post.PublishedOnBsky, "Bluesky", post.BskyLink),
		format(post.PublishedOnInstagram, "Instagram", post.BskyLink),
		format(post.PublishedOnTelegram, "Telegram", post.TelegramLink),
		format(post.PublishedOnTwitter, "Twitter", post.TwitterLink),
	)

	return msg
}
