package model

import (
	"time"

	"github.com/lib/pq"
)

type Post struct {
	Id      int `gorm:"primaryKey";autoIncrement:true"`
	MediaId int // Documents in the same message have the same MediaId

	// Data
	Text   string
	Images pq.StringArray `gorm:"type:text[]"`

	PublishedOnBsky     bool
	PublishedOnTelegram bool
	PublishedOnTwitter  bool

	// Meta
	CreatedAt time.Time `gorm:"type:DATE;"`
	HasText   bool
	HasImages bool
}
