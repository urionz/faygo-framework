package throttle

import (
	"crypto/sha1"
	"faygo-framework/config"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/henrylee2cn/faygo"
)

var Instance = Throttle{
	MaxAttempts:   60,
	DecayDuration: time.Second * 10,
}

type Throttle struct {
	MaxAttempts   int
	DecayDuration time.Duration
	Default       faygo.HandlerFunc
}

const (
	// DefaultMaxAttempts 默认最大请求次数为60次
	DefaultMaxAttempts = 60
	// DefaultDecayDuration 默认频率控制周期为1分钟
	DefaultDecayDuration = time.Minute
)

func (t *Throttle) Prepare() {
	if t.MaxAttempts == 0 {
		t.MaxAttempts = DefaultMaxAttempts
	}

	if t.DecayDuration == 0 {
		t.DecayDuration = DefaultDecayDuration
	}
}

func (t *Throttle) Middleware() faygo.HandlerFunc {
	t.Prepare()
	return func(ctx *faygo.Context) error {
		// 获取唯一签名
		key := resolveRequestSignature(ctx)
		// 超过尝试次数 返回 429
		if tooManyAttempts(key, t.MaxAttempts) {
			var timer time.Time
			config.Cache.GetScan(key+":timer", &timer)
			retryAfter := timer.Sub(time.Now())
			headers := getHeaders(t.MaxAttempts, calculateRemainingAttempts(key, t.MaxAttempts, retryAfter), retryAfter)
			for k, v := range headers {
				ctx.SetHeader(k, v)
			}
			return ctx.JSON(http.StatusTooManyRequests, "")
		}
		// 命中存在的访问者
		hit(key, t.DecayDuration)
		// 设置header
		for k, v := range getHeaders(t.MaxAttempts, calculateRemainingAttempts(key, t.MaxAttempts, 0), 0) {
			ctx.SetHeader(k, v)
		}
		return nil
	}
}

// 命中唯一访问者
func hit(key string, decayMinutes time.Duration) int {
	hits := 0
	config.Cache.Add(key+":timer", availableAt(decayMinutes), decayMinutes)
	added := config.Cache.Add(key, 0, decayMinutes)
	config.Cache.Increment(key, 1)
	if err := config.Cache.GetScan(key, &hits); err != nil && !added && hits == 1 {
		config.Cache.Put(key, 1, decayMinutes)
	}
	return hits
}

// 计算下一次可用时间
func availableAt(delay time.Duration) time.Time {
	return time.Now().Add(delay)
}

func getHeaders(maxAttempts, remainingAttempts int, retryAfter time.Duration) map[string]string {
	headers := map[string]string{
		"X-RateLimit-Limit":     strconv.Itoa(maxAttempts),
		"X-RateLimit-Remaining": strconv.Itoa(remainingAttempts),
	}
	if retryAfter.Nanoseconds() != 0 {
		headers["Retry-After"] = retryAfter.String()
		headers["X-RateLimit-Reset"] = availableAt(retryAfter).Format("2006-01-02 15:04:05")
	}
	return headers
}

// 判断是否超出尝试次数
func tooManyAttempts(key string, maxAttempts int) bool {
	attempts := 0
	if err := config.Cache.GetScan(key, &attempts); err != nil {
		return false
	}
	if attempts >= maxAttempts {
		if config.Cache.Has(key + ":timer") {
			return true
		}
		config.Cache.Forget(key)
	}
	return false
}

// 计算唯一签名
func resolveRequestSignature(context *faygo.Context) string {
	domain := context.Domain()
	ip := context.IP()
	hash := sha1.New()
	io.WriteString(hash, domain+"|"+ip)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// 计算剩余尝试次数
func calculateRemainingAttempts(key string, maxAttempts int, retryAfter time.Duration) int {
	attempts := 0
	err := config.Cache.GetScan(key, &attempts)
	if err != nil {
		return 0
	}
	if retryAfter.Nanoseconds() == 0 {
		return maxAttempts - attempts
	}
	return 0
}
