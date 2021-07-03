package cache

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/colinrs/pkgx/logger"
)

const (
	cmdGet    = "get"
	cmdSet    = "set"
	cmdDel    = "del"

	statInterval     = time.Minute

)

type fetchFunc func() (interface{}, error)

// Cache ...
type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error)
	Get(ctx context.Context, key string, fetch fetchFunc) (result []byte,err error)
	Del(ctx context.Context, key string) (err error)
	AddPlugin(p Plugin)
}

type cacheStat struct {
	hit          uint64
	miss         uint64
	localCacheHit uint64
	localCacheMiss uint64
}


func newCacheStat() *cacheStat {
	st := &cacheStat{
	}
	go st.statLoop()
	return st
}

func (cs *cacheStat) IncrementHit() {
	atomic.AddUint64(&cs.hit, 1)
}

func (cs *cacheStat) IncrementMiss() {
	atomic.AddUint64(&cs.miss, 1)
}

func (cs *cacheStat) IncrementLocalCacheHit() {
	atomic.AddUint64(&cs.localCacheHit, 1)
}

func (cs *cacheStat) IncrementLocalCacheMiss() {
	atomic.AddUint64(&cs.localCacheMiss, 1)
}

func (cs *cacheStat) statLoop() {
	ticker := time.NewTicker(statInterval)
	defer ticker.Stop()

	for range ticker.C {
		hit := atomic.SwapUint64(&cs.hit, 0)
		miss := atomic.SwapUint64(&cs.miss, 0)
		localCacheHit := atomic.SwapUint64(&cs.localCacheHit, 0)
		localCacheMiss := atomic.SwapUint64(&cs.localCacheMiss, 0)
		total := hit + miss + localCacheHit + localCacheMiss
		if total == 0 {
			continue
		}
		percent := 100 * float32(hit+localCacheHit) / float32(total)
		logger.Info("hit_ratio: %0.2f, elements get total: %d, hit: %d, miss: %d, local cache hit:%d, local cache miss:%d" ,
			percent, total, hit, miss, localCacheHit, localCacheMiss)
	}
}

