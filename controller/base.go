package controller

import (
	"net/http"

	"github.com/henrylee2cn/faygo"
)

type Base struct {
}

func (b *Base) Success(ctx *faygo.Context, data interface{}) error {
	return ctx.JSONMsg(http.StatusOK, 0, "")
}

func (b *Base) Error(ctx *faygo.Context, err error) error {
	return ctx.JSONMsg(http.StatusInternalServerError, 1, err.Error())
}
