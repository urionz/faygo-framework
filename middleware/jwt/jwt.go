package jwt

import (
	"fay-blog/config"
	"fay-blog/models"
	"strconv"
	"time"

	"github.com/henrylee2cn/faygo"
	"github.com/henrylee2cn/faygo/ext/db/gorm"
	"github.com/henrylee2cn/faygo/ext/middleware/jwt"
)

var user *models.User
var administrator *models.Administrator

var Authorizator = func(username string, c *faygo.Context) bool {
	if config.GlobalConfig.Jwt.Enable {
		if username != "" {
			return true
		}
		return false
	}
	return true
}

var Unauthorized = func(c *faygo.Context, code int, message string) {
	c.JSON(code, map[string]interface{}{
		"code":    code,
		"message": message,
	})
}

var PayloadFunc = func(username string) map[string]interface{} {
	return map[string]interface{}{
		"username": user.Username,
	}
}

var UserJwt = &jwt.FaygoJWTMiddleware{
	Realm:      "jwt for user",
	Key:        []byte(config.GlobalConfig.App.AppKey),
	Timeout:    config.GlobalConfig.Jwt.Timeout,
	MaxRefresh: config.GlobalConfig.Jwt.MaxRefresh,
	Authenticator: func(username string, password string, c *faygo.Context) (string, bool) {
		if config.GlobalConfig.Jwt.Enable {
			db := gorm.MustDB()
			user = &models.User{}
			if err := models.NewUserQuerySet(db).UsernameEq(username).PasswordEq(password).One(user); err != nil {
				return "", false
			}
			return strconv.FormatUint(uint64(user.ID), 10), true
		} else {
			return "", true
		}
	},
	Authorizator: Authorizator,
	Unauthorized: Unauthorized,
	PayloadFunc: func(username string) map[string]interface{} {
		return map[string]interface{}{
			"username": user.Username,
		}
	},
	TokenLookup:   "header:Authorization",
	TokenHeadName: "Bearer",
	TimeFunc:      time.Now,
}

var AdministratorJwt = &jwt.FaygoJWTMiddleware{
	Realm:      "jwt for administrator",
	Key:        []byte(config.GlobalConfig.App.AppKey),
	Timeout:    config.GlobalConfig.Jwt.Timeout,
	MaxRefresh: config.GlobalConfig.Jwt.MaxRefresh,
	Authenticator: func(username string, password string, c *faygo.Context) (string, bool) {
		db := gorm.MustDB()
		administrator = &models.Administrator{}
		if err := models.NewAdministratorQuerySet(db).UsernameEq(username).PasswordEq(password).One(administrator); err != nil {
			return "", false
		}
		return strconv.FormatUint(uint64(administrator.ID), 10), true
	},
	Authorizator: Authorizator,
	Unauthorized: Unauthorized,
	PayloadFunc: func(username string) map[string]interface{} {
		return map[string]interface{}{
			"username": administrator.Username,
		}
	},
	TokenLookup:   "header:Authorization",
	TokenHeadName: "Bearer",
	TimeFunc:      time.Now,
}

var UserLoginHandler = faygo.WrapDoc(faygo.HandlerFunc(UserJwt.LoginHandler), "", "", faygo.ParamInfo{
	Name:     "params",
	Required: true,
	In:       "body",
	Model:    "{\n\t\"username\": \"\", \n\t\"password\": \"\", \n\t\"captcha_id\": \"\", \n\t\"captcha\": \"\"\n}",
})

var UserRefreshTokenHandler = faygo.WrapDoc(faygo.HandlerFunc(UserJwt.RefreshHandler), "", "", faygo.ParamInfo{
	Name:     "Authorization",
	Required: true,
	In:       "header",
	Model:    "Bearer ",
})

var AdminLoginHandler = faygo.WrapDoc(faygo.HandlerFunc(AdministratorJwt.LoginHandler), "", "", faygo.ParamInfo{
	Name:     "params",
	Required: true,
	In:       "body",
	Model:    "{\n\t\"username\": \"\", \n\t\"password\": \"\", \n\t\"captcha_id\": \"\", \n\t\"captcha\": \"\"\n}",
})

var AdminRefreshTokenHandler = faygo.WrapDoc(faygo.HandlerFunc(AdministratorJwt.RefreshHandler), "", "", faygo.ParamInfo{
	Name:     "Authorization",
	Required: true,
	In:       "header",
	Model:    "Bearer ",
})
