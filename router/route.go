package router

import (
	"fay-blog/controller"
	"fay-blog/middleware/captcha"
	"fay-blog/middleware/host"
	"fay-blog/middleware/jwt"
	"fay-blog/middleware/pprof"
	"fay-blog/middleware/throttle"

	"github.com/henrylee2cn/faygo"
)

// Route register router
func Route(frame *faygo.Framework) {
	frame.Route(
		frame.NewGroup("",
			// pprof 路由注册
			pprof.Instance.RouteWrap(frame),
			// 公共路由模块
			frame.NewNamedGroup("公共接口", "",
				frame.NewNamedGET("首页", "index", &controller.Index{}),
				frame.NewNamedGET("验证码", "captcha", captcha.Instance.Serve()),
			),
			frame.NewNamedGroup("用户接口", "user",
				frame.NewNamedGroup("需要认证", "",
					frame.NewNamedGET("用户token刷新", "refresh_token", jwt.UserRefreshTokenHandler),
					frame.NewNamedPOST("用户创建", "create", &controller.CreateUser{}),
					frame.NewNamedPUT("用户更新", "update/:uid", &controller.UpdateUser{}),
					frame.NewNamedDELETE("删除用户", "delete/:uid", &controller.DeleteUser{}),
				).Use(jwt.UserJwt.MiddlewareFunc()),
				frame.NewNamedGroup("不需要认证", "",
					frame.NewNamedPOST("注册Api", "register", &controller.UserRegister{}),
					frame.NewNamedPOST("登录Api", "login", jwt.UserLoginHandler).Use(captcha.Instance.Middleware()),
				),
			),
			frame.NewNamedGroup("后台接口", "admin",
				frame.NewNamedGroup("需要认证", "",
					frame.NewNamedGET("管理员token刷新", "refresh_token", jwt.AdminRefreshTokenHandler),
				).Use(jwt.AdministratorJwt.MiddlewareFunc()),
				frame.NewNamedGroup("不需要认证", "",
					frame.NewNamedPOST("登录API", "login", jwt.AdminLoginHandler).Use(captcha.Instance.Middleware()),
				),
			),
		).Use(host.MapTplHostToCtx("{subdomain}.urionz.com"), throttle.Instance.Middleware()),
	)
}
