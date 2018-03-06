package pprof

import (
	"faygo-framework/config"
	"faygo-framework/middleware/ip"
	"net/http/pprof"

	"github.com/henrylee2cn/faygo"
)

var Instance = PProf{}

type PProf struct {
}

func (p *PProf) RouteWrap(frame *faygo.Framework) *faygo.MuxAPI {
	var router *faygo.MuxAPI
	config := config.GlobalConfig.PProf
	if config.Enable {
		if config.NoLimit {
			router = frame.NewGroup(config.Prefix,
				frame.NewGET("index", p.Index()),
				frame.NewGET("heap", p.Heap()),
				frame.NewGET("goroutine", p.Goroutine()),
				frame.NewGET("block", p.Block()),
				frame.NewGET("threadcreate", p.ThreadCreate()),
				frame.NewGET("cmdline", p.Cmdline()),
				frame.NewGET("profile", p.Profile()),
				frame.NewPOST("symbol", p.Symbol()),
				frame.NewGET("trace", p.Trace()),
				frame.NewGET("mutex", p.Mutex()),
			)
		} else {
			router = frame.Filter(ip.Filter(config.Whitelist, config.RealIp)).NewGroup(config.Prefix,
				frame.NewGET("index", p.Index()),
				frame.NewGET("heap", p.Heap()),
				frame.NewGET("goroutine", p.Goroutine()),
				frame.NewGET("block", p.Block()),
				frame.NewGET("threadcreate", p.ThreadCreate()),
				frame.NewGET("cmdline", p.Cmdline()),
				frame.NewGET("profile", p.Profile()),
				frame.NewPOST("symbol", p.Symbol()),
				frame.NewGET("trace", p.Trace()),
				frame.NewGET("mutex", p.Mutex()),
			)
		}
	} else {
		router = frame.NewGroup("")
	}
	return router
}

func (*PProf) Index() faygo.HandlerFunc {
	return func(ctx *faygo.Context) error {
		pprof.Index(ctx.W, ctx.R)
		return nil
	}
}

func (*PProf) Heap() faygo.HandlerFunc {
	return func(ctx *faygo.Context) error {
		pprof.Handler("heap").ServeHTTP(ctx.W, ctx.R)
		return nil
	}
}

func (*PProf) Goroutine() faygo.HandlerFunc {
	return func(ctx *faygo.Context) error {
		pprof.Handler("goroutine").ServeHTTP(ctx.W, ctx.R)
		return nil
	}
}

func (*PProf) Block() faygo.HandlerFunc {
	return func(ctx *faygo.Context) error {
		pprof.Handler("block").ServeHTTP(ctx.W, ctx.R)
		return nil
	}
}

func (*PProf) ThreadCreate() faygo.HandlerFunc {
	return func(ctx *faygo.Context) error {
		pprof.Handler("threadcreate").ServeHTTP(ctx.W, ctx.R)
		return nil
	}
}

func (*PProf) Cmdline() faygo.HandlerFunc {
	return func(ctx *faygo.Context) error {
		pprof.Cmdline(ctx.W, ctx.R)
		return nil
	}
}

func (*PProf) Profile() faygo.HandlerFunc {
	return func(ctx *faygo.Context) error {
		pprof.Profile(ctx.W, ctx.R)
		return nil
	}
}

func (*PProf) Symbol() faygo.HandlerFunc {
	return func(ctx *faygo.Context) error {
		pprof.Symbol(ctx.W, ctx.R)
		return nil
	}
}

func (*PProf) Trace() faygo.HandlerFunc {
	return func(ctx *faygo.Context) error {
		pprof.Trace(ctx.W, ctx.R)
		return nil
	}
}

func (*PProf) Mutex() faygo.HandlerFunc {
	return func(ctx *faygo.Context) error {
		pprof.Handler("mutex").ServeHTTP(ctx.W, ctx.R)
		return nil
	}
}
