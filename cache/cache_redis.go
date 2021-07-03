package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/colinrs/pkgx/logger"
	"github.com/coocood/freecache"
	"github.com/go-redis/redis"
	"github.com/golang/groupcache/singleflight"
	"github.com/tal-tech/go-zero/core/mathx"
)

var (
	// make the unstable expiry to be [0.95, 1.05] * seconds
	expiryDeviation = 0.05
	defaultExpire = 5*time.Minute
)

type RedisConfig struct {
	Addr        string `yaml:"addr" json:"addr"`
	DB          int    `yaml:"db" json:"db"`
	PoolSize    int    `yaml:"pool_size" json:"pool_size"`
	IdleTimeout int    `yaml:"idle_timeout" json:"idle_timeout"`
	Prefix      string `yaml:"prefix" json:"prefix"`
	LocalCacheSize int `yaml:"local_cache_size" json:"local_cache_size"` // M
}

type flightGroup interface {
	// Done is called when Do is done.
	Do(key string, fn func() (interface{}, error)) (interface{}, error)
}


type RedisCacheClient struct {
	client         *redis.Client
	prefix         string
	plugins        []Plugin
	status         *cacheStat
	unstableExpiry mathx.Unstable
	loadGroup      flightGroup
	DefaultExpire  time.Duration
	localCache  *freecache.Cache
}

var _ Cache = (*RedisCacheClient)(nil)

var DefaultRedisClient *RedisCacheClient

func getFullKey(prefix, key string) string {
	return prefix + "_" + key
}

func InitCacheClient(conf *RedisConfig) Cache {
	DefaultRedisClient = &RedisCacheClient{
		status: newCacheStat(),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		loadGroup:  &singleflight.Group{},
		DefaultExpire: defaultExpire,
		// conf.LocalCacheSize: M
		localCache: freecache.NewCache(conf.LocalCacheSize * 1024 * 1024),
	}
	DefaultRedisClient.client = redis.NewClient(&redis.Options{
		Addr:        conf.Addr,
		DB:          conf.DB,
		PoolSize:    conf.PoolSize,
		IdleTimeout: time.Duration(conf.IdleTimeout) * time.Second,
	})
	if conf.Prefix != "" {
		DefaultRedisClient.prefix = conf.Prefix
	}
	return DefaultRedisClient
}

func (r *RedisCacheClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error) {
	var byteValue []byte
	startTime := time.Now()
	fullKey := getFullKey(r.prefix, key)
	fullKeyByte, _ := json.Marshal(fullKey)
	if byteValue, err = json.Marshal(value); err != nil {
		logger.Error("json.Marshal redis value: %v, error: %v", value, err)
		return err
	}
	_ = r.localCache.Set(fullKeyByte, byteValue, int(expiration.Seconds()))
	expiration = r.unstableExpiry.AroundDuration(expiration)
	err = r.client.Set(fullKey, byteValue, expiration).Err()
	elapsed := time.Since(startTime).Milliseconds()
	for _, p := range r.plugins {
		p.OnSetRequestEnd(ctx, cmdSet, elapsed, fullKey, err)
	}
	if err != nil {
		logger.Error("set redis key: %v, error: %v", fullKey, err)
		return err
	}
	return nil
}


func (r *RedisCacheClient) Get(ctx context.Context, key string, fetch fetchFunc) (result []byte, err error) {
	var byteValue []byte
	fullKey := getFullKey(r.prefix, key)
	fullKeyByte, _ := json.Marshal(fullKey)
	if val, err := r.localCache.Get(fullKeyByte); err==nil {
		r.status.IncrementLocalCacheHit()
		return val, nil
	}
	r.status.IncrementLocalCacheMiss()
	startTime := time.Now()
	byteValue, err = r.client.Get(fullKey).Bytes()
	elapsed := time.Since(startTime).Milliseconds()
	for _, p := range r.plugins {
		p.OnGetRequestEnd(ctx, cmdGet, elapsed, fullKey, err)
	}
	// not found key
	if err == redis.Nil {
		r.status.IncrementMiss()
		if fetch!=nil {

			var b []byte
			_, err = r.loadGroup.Do(fullKey, func() (interface{}, error) {
				var fetchResult interface{}
				if val, err := r.localCache.Get(fullKeyByte); err == nil {
					return val, nil
				}
				v, e := fetch()
				if e != nil {
					logger.Error("get redis key: %v, from fetch error: %v", fullKey, e)
					return nil, e
				}
				expiration := r.unstableExpiry.AroundDuration(r.DefaultExpire)
				b, _ = json.Marshal(fetchResult)

				_ = r.localCache.Set(fullKeyByte, b, int(expiration.Seconds()))
				return v, nil
			})
			if err != nil {
				return nil, err
			}
			return b, nil
		}
	}
	// something err get key from redis
	if err != nil {
		logger.Error("get redis key: %v, error: %v", fullKey, err)
		return nil, err
	}
	r.status.IncrementHit()
	return byteValue, nil
}

func (r *RedisCacheClient) Del(ctx context.Context, key string) (err error) {
	fullKey := getFullKey(r.prefix, key)
	startTime := time.Now()
	_, err = r.client.Del(key).Result()
	elapsed := time.Since(startTime).Milliseconds()
	for _, p := range r.plugins {
		p.OnGetRequestEnd(ctx, cmdDel, elapsed, fullKey, err)
	}
	// not found key
	if err == redis.Nil {
		return nil
	}
	// something err get key from redis
	if err != nil {
		return err
	}
	return nil
}


func (r *RedisCacheClient) AddPlugin(p Plugin) {
	r.plugins = append(r.plugins, p)
}
