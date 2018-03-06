package config

import (
	"fay-blog/utils"
	"time"

	"github.com/henrylee2cn/faygo"
	"github.com/urionz/store"
)

var Cache *store.Container

const (
	CONFIG_FILE = faygo.CONFIG_DIR + "config.ini"
)

type (
	Config struct {
		App     AppConfig     `ini:"application" comment:"application section"`
		Cache   CacheConfig   `ini:"cache" comment:"Cache section"`
		PProf   PProfConfig   `ini:"pprof" comment:"pprof section"`
		Jwt     JwtConfig     `ini:"jwt" comment:"jwt section"`
		Captcha CaptchaConfig `ini:"captcha" comment:"captcha section"`
	}
	AppConfig struct {
		AppKey string `ini:"app_key" comment:"the application global app key"`
		AppEnv string `ini:"app_env" comment:"the application runtime env"`
	}
	CacheConfig struct {
		Enable          bool          `ini:"enable" comment:"Enable the config section"`
		Driver          string        `ini:"driver" comment:"redis(go-redis) | memory(gocache)"`
		Prefix          string        `ini:"prefix" comment:"cache key prefix"`
		Expiration      time.Duration `ini:"expiration" comment:"default expiration time(minute)"`
		CleanupInterval time.Duration `ini:"cleanup_interval" comment:"default cleanup interval time(minute)"`
		Addr            string        `ini:"addr" comment:"redis tcp connect addr"`
		Password        string        `ini:"password" comment:"redis password"`
		DB              int           `ini:"default_db" comment:"redis default db"`
	}
	PProfConfig struct {
		Enable    bool     `ini:"enable" comment:"Enable the config section"`
		Prefix    string   `ini:"prefix" comment:"pprof route prefix"`
		NoLimit   bool     `ini:"no_limit" comment:"If true, access is not restricted"`
		RealIp    bool     `ini:"real_ip" comment:"if true, means verifying the real IP of the visitor"`
		Whitelist []string `ini:"whitelist" delim:"|" comment:"'whitelist=192.*|202.122.246.170' means: only IP addresses that are prefixed with '192.' or equal to '202.122.246.170' are allowed"`
	}
	JwtConfig struct {
		Enable     bool          `ini:"enable" comment:"Enable the config section"`
		Timeout    time.Duration `ini:"timeout" comment:"jwt timeout"`
		MaxRefresh time.Duration `ini:"max_refresh" comment:"jwt token max refresh time"`
	}
	CaptchaConfig struct {
		Enable            bool          `ini:"enable" comment:"Enable the config section"`
		Expiration        time.Duration `ini:"expiration" comment:"captcha expiration"`
		Length            int           `ini:"length" comment:"captcha length"`
		Width             int           `ini:"width" comment:"captcha width"`
		Height            int           `ini:"height" comment:"captcha height"`
		TriggerRecycleNum int           `ini:"trigger_recycle_num" comment:"The number of captchas created that triggers garbage collection used"`
	}
)

var (
	GlobalConfig = Config{
		App: AppConfig{
			AppKey: utils.Base64Encode(faygo.RandomString(20)),
			AppEnv: "development",
		},
		Cache: CacheConfig{
			Enable:          true,
			Driver:          "redis",
			Prefix:          "",
			Expiration:      60,
			CleanupInterval: 60,
			Addr:            "localhost:6379",
			Password:        "",
			DB:              0,
		},
		PProf: PProfConfig{
			Enable:    false,
			Prefix:    "debug",
			NoLimit:   false,
			RealIp:    false,
			Whitelist: []string{"127.*", "192.168.*"},
		},
		Jwt: JwtConfig{
			Enable:     true,
			Timeout:    time.Hour,
			MaxRefresh: time.Hour,
		},
		Captcha: CaptchaConfig{
			Enable:     true,
			Expiration: time.Minute * 2,
			Length:     4,
			Width:      240,
			Height:     80,
		},
	}
)

func Load() error {
	if err := faygo.SyncINI(&GlobalConfig, func(onecUpdateFunc func() error) error {
		return onecUpdateFunc()
	}, CONFIG_FILE); err != nil {
		return err
	}
	newGlobalCache(&GlobalConfig.Cache)
	return nil
}

func newGlobalCache(config *CacheConfig) {
	Cache = store.New(config.Driver, store.Container{
		Addr:            config.Addr,
		Password:        config.Password,
		DB:              config.DB,
		Prefix:          config.Prefix,
		Expiration:      config.Expiration,
		CleanupInterval: config.CleanupInterval,
	})
}
