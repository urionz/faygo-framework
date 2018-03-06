package host

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/henrylee2cn/faygo"
)

// MapTplHostToCtx 映射模板域名到上下文变量中
// 用法：MapTplHostToCtx("{subDomain}.{domain}.com")
// 获取：ctx.Data("subDomain")  ctx.Data("domain")
func MapTplHostToCtx(tpl string) faygo.HandlerFunc {
	return func(ctx *faygo.Context) error {
		hs, ts, tmp := strings.Split(ctx.Domain(), "."), strings.Split(tpl, "."), ""
		if len(hs) < 2 {
			ctx.Log().Warn("host length less 2")
			return nil
		}
		reg := regexp.MustCompile(`{(\w+)}`)
		tmp = hs[1]
		if len(hs) != 3 {
			hs[0], hs[1] = "", hs[0]
			hs = append(hs, tmp)
		}
		tmp = ts[1]
		if len(ts) != 3 {
			ts[0], ts[1] = "", ts[0]
			ts = append(ts, tmp)
		}
		for index, name := range ts {
			if strings.Contains(name, "{") {
				match := reg.FindSubmatchIndex([]byte(name))
				ctxDataName := fmt.Sprintf("%s", reg.Expand([]byte(""), []byte("$1"), []byte(name), match))
				ctx.SetData(ctxDataName, hs[index])
			} else {
				continue
			}
		}
		return nil
	}
}
