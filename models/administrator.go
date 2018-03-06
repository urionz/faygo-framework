package models

import "github.com/jinzhu/gorm"

//go:generate goqueryset -in administrator.go

// gen:qs
type Administrator struct {
	gorm.Model
	Username string `gorm:"unique;index:username;not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
}

func init() {
	RegisterModel(&Administrator{})
}

func (Administrator) TableName() string {
	return "administrator"
}
