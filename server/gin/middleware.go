package gin

import (

	"net/http"
	"time"

	"github.com/juju/ratelimit"
	rgin "github.com/gin-gonic/gin"
)
func Logger() rgin.HandlerFunc{
	return rgin.Logger()
}

func Recovery() rgin.HandlerFunc{
	return rgin.Recovery()
}

func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *rgin.Context) {
	bucket := ratelimit.NewBucket(fillInterval, cap)
	return func(c *rgin.Context) {
		// 如果取不到令牌就中断本次请求返回 rate limit...
		if bucket.TakeAvailable(1) < 1 {
			c.String(http.StatusOK, "rate limit...")
			c.Abort()
			return
		}
		c.Next()
	}
}