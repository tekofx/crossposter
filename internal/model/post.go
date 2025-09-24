package model

import "time"

type Post struct {
	Id                  int    `gorm:"primaryKey";autoIncrement:false"`
	BskyId              string `gorm:"primaryKey";autoIncrement:false"`
	PublishedOnTelegram bool
	PublishedOnTwitter  bool
	Date                time.Time `gorm:"type:DATE;"`
}
