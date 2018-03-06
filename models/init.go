package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/henrylee2cn/faygo/ext/db/gorm"
	"github.com/henrylee2cn/goutil/errors"
)

var models []interface{}

func RegisterModel(model interface{}) {
	models = append(models, model)
}

func AutoMigrate() error {
	var err error
	if db, ok := gorm.DB("default"); ok {
		for _, model := range models {
			if db.HasTable(model) {
				err = db.AutoMigrate(model).Error
			} else {
				err = db.CreateTable(model).Error
			}
		}
	} else {
		err = errors.Errorf("db connection fail...")
	}
	return err
}
