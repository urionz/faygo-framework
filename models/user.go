package models

import (
	"github.com/jinzhu/gorm"
)

//go:generate goqueryset -in user.go

// gen:qs
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null;index:username" json:"username"`
	Password string `gorm:"not null" json:"password"`
}

func init() {
	RegisterModel(&User{})
}

func (User) TableName() string {
	return "user"
}
