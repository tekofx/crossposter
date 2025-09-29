package model

import "gorm.io/gorm"

type Image struct {
	gorm.Model
	Filename string
	PostID   uint
}
