package ratelimit

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis"
	"github.com/spf13/cast"
)

var (
	standAloneCounterRateLimiterOnce sync.Once
	OneArgs                          = 1
	Limit                            = 0
)

// StandAloneCounterRateLimiter ...
type StandAloneCounterRateLimiter struct {
	snippet         time.Duration
	currentRequests int32
	allowRequests   int32
}

func NewStandAloneCounterRateLimiter(snippet time.Duration, allowRequests int32) *StandAloneCounterRateLimiter {
	return &StandAloneCounterRateLimiter{snippet: snippet, allowRequests: allowRequests}
}

func (l *StandAloneCounterRateLimiter) Take() error {
	standAloneCounterRateLimiterOnce.Do(func() {
		go func() {
			for {
				select {
				case <-time.After(l.snippet):
					atomic.StoreInt32(&l.currentRequests, 0)
				}
			}
		}()
	})

	curRequest := atomic.LoadInt32(&l.currentRequests)
	if curRequest >= l.allowRequests {
		return ErrExceededLimit
	}
	if !atomic.CompareAndSwapInt32(&l.currentRequests, curRequest, curRequest+1) {
		return ErrExceededLimit
	}
	return nil
}

// DistributedCounterRateLimiter ...
type DistributedCounterRateLimiter struct {
	snippet       time.Duration
	allowRequests int32
	redisClient   *redis.Client
	hashID        string
	key           string
}

func NewDistributedCounterRateLimiter(
	snippet time.Duration,
	allowRequests int32,
	client *redis.Client,
	hashID string,
	key string) *DistributedCounterRateLimiter {
	c := &DistributedCounterRateLimiter{
		snippet:       snippet,
		allowRequests: allowRequests,
		redisClient:   client,
		hashID:        hashID,
		key:           key,
	}

	return c
}

func (l *DistributedCounterRateLimiter) Take() error {
	result := l.redisClient.Do("Evalsha", l.hashID, OneArgs, l.key, l.snippet, l.allowRequests)
	ok, err := result.Result()
	if err != nil {
		return err
	}
	if cast.ToInt(ok) == Limit {
		return ErrExceededLimit
	}
	return nil
}

/*

local key = KEYS[1]

local capacity = tonumber(ARGV[1])
local expire = tonumber(ARGV[2])

local current_num = tonumber(redis.call("incr", key))

local t = redis.call('ttl',key)
if current_num - 1 == 0 and t == -1 then
  redis.call('expire',key,expire)
end

local ok = 0
if current_num < capacity then
  ok = 1
end

return {ok}

➜  redis-6.0.8 ./src/redis-cli -x script load < lua/count.lua // 加载脚本
"a61ab0c62af637febce19dd0c563a8cb05fc1ac6"
- SCRIPT FLUSH // 删除脚本
-
*/
