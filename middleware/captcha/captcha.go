package captcha

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"faygo-framework/config"
	"faygo-framework/utils"
	"net/http"

	"github.com/dchest/captcha"
	"github.com/henrylee2cn/faygo"
)

var Instance = Captcha{
	Store: &Store{
		Expiration: config.GlobalConfig.Captcha.Expiration,
	},
}

type Captcha struct {
	Store     *Store `json:"-"`
	CaptchaId string `json:"captcha_id"`
	Captcha   string `json:"captcha"`
}

func (c *Captcha) Serve() faygo.HandlerFunc {
	return func(ctx *faygo.Context) error {
		if config.GlobalConfig.Captcha.Enable {
			width, height, length := captcha.StdWidth, captcha.StdHeight, captcha.DefaultLen
			if config.GlobalConfig.Captcha.Width != 0 {
				width = config.GlobalConfig.Captcha.Width
			}
			if config.GlobalConfig.Captcha.Height != 0 {
				height = config.GlobalConfig.Captcha.Height
			}
			if config.GlobalConfig.Captcha.Length != 0 {
				length = config.GlobalConfig.Captcha.Length
			}
			captcha.SetCustomStore(c.Store)
			captchaId := captcha.NewLen(length)
			buffer := bytes.Buffer{}
			if err := captcha.WriteImage(&buffer, captchaId, width, height); err != nil {
				ctx.Log().Error(err)
				return err
			}
			base64Img := base64.StdEncoding.EncodeToString(buffer.Bytes())
			ctx.JSONMsg(http.StatusOK, 0, map[string]string{
				"captcha_id": captchaId,
				"captcha":    `data:image/png;base64,` + base64Img,
			})
		}
		return nil
	}
}

func (c *Captcha) Middleware() faygo.HandlerFunc {
	return func(ctx *faygo.Context) error {
		if config.GlobalConfig.Captcha.Enable {
			var raw []byte
			var err error
			if raw, err = utils.CopyBody(ctx.R); err != nil {
				ctx.Log().Error(err)
				return nil
			}
			if err := json.Unmarshal(raw, &Instance); err != nil {
				ctx.Log().Error(err)
				return nil
			}
			if Instance.CaptchaId == "" || Instance.Captcha == "" {
				return ctx.JSONMsg(http.StatusForbidden, -1, "请填写验证码和验证码ID")
			}
			if ok := captcha.VerifyString(Instance.CaptchaId, Instance.Captcha); ok {
				return nil
			}
			return ctx.JSONMsg(http.StatusForbidden, -1, "验证码错误")
		}
		return nil
	}
}
