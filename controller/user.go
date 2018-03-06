package controller

import (
	"faygo-framework/models"

	"github.com/henrylee2cn/faygo"
	"github.com/henrylee2cn/faygo/ext/db/gorm"
	"github.com/henrylee2cn/goutil/errors"
)

type UserRegister struct {
	Base
	Username string `param:"<in:formData> <required> <name:username>"`
	Password string `param:"<in:formData> <required> <name:password>"`
}

func (u *UserRegister) Serve(ctx *faygo.Context) error {
	user := models.User{
		Username: ctx.FormParam("username"),
		Password: ctx.FormParam("password"),
	}
	err := user.Create(gorm.MustDB())
	if err != nil {
		return errors.Errorf(err.Error())
	}
	return ctx.JSONMsg(200, 0, "")
}

type UserLogin struct {
	Base
	Username string `param:"<in:formData> <required> <name:username>"`
	Password string `param:"<in:formData> <required> <name:password>"`
}

func (u *UserLogin) Serve(ctx *faygo.Context) error {
	var user *models.User
	user.Username = ctx.FormParam("username")
	user.Password = ctx.FormParam("password")
	if err := models.NewUserQuerySet(gorm.MustDB()).One(user); err != nil {
		return ctx.JSON(200, "")
	} else {
		return err
	}
}

type CreateUser struct {
	Base
}

func (c *CreateUser) Serve(ctx *faygo.Context) error {
	return nil
}

type UpdateUser struct {
	Base
	Token string `param:"<in:header> <required> <name:Authorization>"`
	UID   uint   `param:"<in:path> <required> <name:uid> <desc:用户ID>"`
}

func (u UpdateUser) Serve(ctx *faygo.Context) error {
	return nil
}

type DeleteUser struct {
	Base
	Token string `param:"<in:header> <required> <name:Authorization>"`
	UID   uint   `param:"<in:path> <required> <name:uid> <desc:用户ID>"`
}

func (d DeleteUser) Serve(ctx *faygo.Context) error {
	return nil
}
