package models

import "github.com/jinzhu/gorm"

//go:generate goqueryset -in article.go

// gen:qs
type Article struct {
	gorm.Model
	Title string `json:"title"`
}

func init() {
	RegisterModel(&Article{})
}

func (Article) TableName() string {
	return "article"
}
