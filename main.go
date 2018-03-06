package main

import (
	"fay-blog/config"
	"fay-blog/models"
	"fay-blog/router"

	"github.com/henrylee2cn/faygo"
)

func main() {
	frame := faygo.New("fay-blog", "1.0")
	if err := config.Load(); err != nil {
		frame.Log().Error(err)
	}
	if err := models.AutoMigrate(); err != nil {
		frame.Log().Error(err)
	}

	router.Route(frame)
	faygo.Run()
}
