package model

import (
	"time"

	"github.com/lib/pq"
)

type Post struct {
	Id int `gorm:"primaryKey";autoIncrement:true"`

	// Data
	Text   string
	Images pq.StringArray `gorm:"type:text[]"`

	// Meta
	Date time.Time `gorm:"type:DATE;"`
}

var PostToPublish *Post
