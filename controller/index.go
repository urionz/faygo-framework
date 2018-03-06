package controller

import (
	"fmt"

	"github.com/henrylee2cn/faygo"
	"github.com/henrylee2cn/goutil/errors"
)

type Index struct {
	Base
	Name string `param:"<in:query> <required> <desc: name>" json:"name"`
}

func (index *Index) Serve(ctx *faygo.Context) error {
	name := ctx.QueryParam("name")
	fmt.Println(name)
	return index.Error(ctx, errors.New("asds"))
}
