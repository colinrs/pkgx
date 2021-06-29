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

// Cache ...
type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error)
	Get(ctx context.Context, key string, result interface{}) (exist bool, err error)
	Del(ctx context.Context, key string) (err error)
	AddPlugin(p Plugin)
}

type cacheStat struct {
	hit          uint64
	miss         uint64
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

func (cs *cacheStat) statLoop() {
	ticker := time.NewTicker(statInterval)
	defer ticker.Stop()

	for range ticker.C {
		hit := atomic.SwapUint64(&cs.hit, 0)
		miss := atomic.SwapUint64(&cs.miss, 0)
		total := hit + miss
		if total == 0 {
			continue
		}
		percent := 100 * float32(hit) / float32(total)
		logger.Debug("hit_ratio: %.1f%%, elements: %d, hit: %d, miss: %d", percent, total, hit, miss)
	}
}

