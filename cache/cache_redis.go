package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/colinrs/pkgx/logger"
	"github.com/coocood/freecache"
	"github.com/golang/groupcache/singleflight"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/mathx"
)

var (
	// make the unstable expiry to be [0.95, 1.05] * seconds
	expiryDeviation = 0.05
	defaultExpire   = 5 * time.Minute
	NoneValue       = []byte("NoneValue")
)

type RedisConfig struct {
	Addr           string `yaml:"addr" json:"addr"`
	DB             int    `yaml:"db" json:"db"`
	PoolSize       int    `yaml:"pool_size" json:"pool_size"`
	IdleTimeout    int    `yaml:"idle_timeout" json:"idle_timeout"`
	Prefix         string `yaml:"prefix" json:"prefix"`
	LocalCacheSize int    `yaml:"local_cache_size" json:"local_cache_size"` // M
	Username       string `yaml:"username" json:"username"`
	Password       string `yaml:"password" json:"password"`
}

type flightGroup interface {
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
	localCache     *freecache.Cache
}

var _ Cache = (*RedisCacheClient)(nil)

var DefaultRedisClient *RedisCacheClient

func getFullKey(prefix, key string) string {
	return prefix + "_" + key
}

func InitCacheClient(conf *RedisConfig) Cache {
	DefaultRedisClient = &RedisCacheClient{
		status:         newCacheStat(),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		loadGroup:      &singleflight.Group{},
		DefaultExpire:  defaultExpire,
		// conf.LocalCacheSize: M
		localCache: freecache.NewCache(conf.LocalCacheSize * 1024 * 1024),
	}
	DefaultRedisClient.client = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		DB:       conf.DB,
		PoolSize: conf.PoolSize,
		Username: conf.Username,
		Password: conf.Password,
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
	err = r.client.Set(ctx, fullKey, byteValue, expiration).Err()
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

func (r *RedisCacheClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error) {
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
	err = r.client.SetNX(ctx, fullKey, byteValue, expiration).Err()
	elapsed := time.Since(startTime).Milliseconds()
	for _, p := range r.plugins {
		p.OnSetRequestEnd(ctx, cmdSetNX, elapsed, fullKey, err)
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
	if val, err := r.localCache.Get(fullKeyByte); err == nil {
		r.status.IncrementLocalCacheHit()
		return val, nil
	}
	r.status.IncrementLocalCacheMiss()
	startTime := time.Now()
	byteValue, err = r.client.Get(ctx, fullKey).Bytes()
	elapsed := time.Since(startTime).Milliseconds()
	for _, p := range r.plugins {
		p.OnGetRequestEnd(ctx, cmdGet, elapsed, fullKey, err)
	}
	// not found key
	if err == redis.Nil {
		r.status.IncrementMiss()
		if fetch != nil {

			var b []byte
			_, err = r.loadGroup.Do(fullKey, func() (interface{}, error) {
				var fetchResult interface{}
				if val, err := r.localCache.Get(fullKeyByte); err == nil {
					return val, nil
				}
				v, e := fetch()
				if e != nil {
					logger.Error("get redis key: %v, from fetch error: %v", fullKey, e)
					// set none value
					expiration := r.unstableExpiry.AroundDuration(r.DefaultExpire)
					_ = r.localCache.Set(fullKeyByte, NoneValue, int(expiration.Seconds()))
					return NoneValue, nil
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
	_, err = r.client.Del(ctx, key).Result()
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

func (r *RedisCacheClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	fullKey := getFullKey(r.prefix, key)
	startTime := time.Now()
	ttl := r.client.TTL(ctx, key).Val()
	elapsed := time.Since(startTime).Milliseconds()
	for _, p := range r.plugins {
		p.OnGetRequestEnd(ctx, cmdTTL, elapsed, fullKey, nil)
	}
	return ttl, nil
}

func (r *RedisCacheClient) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	fullKey := getFullKey(r.prefix, key)
	startTime := time.Now()
	ok, err := r.client.Expire(ctx, key, expiration).Result()
	elapsed := time.Since(startTime).Milliseconds()
	for _, p := range r.plugins {
		p.OnGetRequestEnd(ctx, cmdTTL, elapsed, fullKey, err)
	}
	return ok, err
}

func (r *RedisCacheClient) AddPlugin(p Plugin) {
	r.plugins = append(r.plugins, p)
}
