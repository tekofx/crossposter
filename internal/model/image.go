package model

import "gorm.io/gorm"

type Image struct {
	gorm.Model
	Filename string
	MimeType string
	FileSize int64
	PostID   uint
}
