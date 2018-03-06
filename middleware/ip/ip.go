package ip

import (
	"net/http"
	"strings"

	"github.com/henrylee2cn/faygo"
)

func Filter(whitelist []string, realIP bool) faygo.HandlerFunc {
	var noAccess bool
	var match []string
	var prefix []string

	if len(whitelist) == 0 {
		noAccess = true
	} else {
		for _, s := range whitelist {
			if strings.HasSuffix(s, "*") {
				prefix = append(prefix, s[:len(s)-1])
			} else {
				match = append(match, s)
			}
		}
	}

	return func(ctx *faygo.Context) error {
		if noAccess {
			ctx.Error(http.StatusForbidden, "no access")
			return nil
		}

		var ip string
		if realIP {
			ip = ctx.RealIP()
		} else {
			ip = ctx.IP()
		}
		for _, ipMatch := range match {
			if ipMatch == ip {
				ctx.Next()
				return nil
			}
		}
		for _, ipPrefix := range prefix {
			if strings.HasPrefix(ip, ipPrefix) {
				ctx.Next()
				return nil
			}
		}
		ctx.Error(http.StatusForbidden, "not allow to access: "+ip)
		return nil
	}
}
