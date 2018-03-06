package models

import "github.com/jinzhu/gorm"

//go:generate goqueryset -in tag.go

// gen:qs
type Tag struct {
	gorm.Model
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

func init() {
	RegisterModel(&Tag{})
}

func (Tag) TableName() string {
	return "tag"
}
