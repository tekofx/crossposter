package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type PostStatus int

const (
	Created PostStatus = iota
	Scheduled
	Posted
)

func (d PostStatus) String() string {
	return [...]string{"Creado", "Programado", "Publicado"}[d]
}

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
	Status    PostStatus
}

func (post *Post) String() string {
	format := func(ok bool, service string, url string) string {
		if ok {
			return fmt.Sprintf("✅ [%s](%s)", service, url)
		}
		return fmt.Sprintf("❌ %s", service)
	}

	return fmt.Sprintf("ID: %d\nEstado: %s\n%s\n%s\n%s\n%s",
		post.ID,
		post.Status.String(),
		format(post.PublishedOnBsky, "Bluesky", post.BskyLink),
		format(post.PublishedOnInstagram, "Instagram", post.BskyLink),
		format(post.PublishedOnTelegram, "Telegram", post.TelegramLink),
		format(post.PublishedOnTwitter, "Twitter", post.TwitterLink),
	)
}
