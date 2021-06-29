package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/colinrs/pkgx/logger"
	"github.com/go-redis/redis"
)

type RedisConfig struct {
	Addr        string `yaml:"addr" json:"addr"`
	DB          int    `yaml:"db" json:"db"`
	PoolSize    int    `yaml:"pool_size" json:"pool_size"`
	IdleTimeout int    `yaml:"idle_timeout" json:"idle_timeout"`
	Prefix      string `yaml:"prefix" json:"prefix"`
}

type RedisCacheClient struct {
	client  *redis.Client
	prefix  string
	plugins []Plugin
	status *cacheStat
}

var _ Cache = (*RedisCacheClient)(nil)

var DefaultRedisClient *RedisCacheClient

func getFullKey(prefix, key string) string {
	return prefix + "_" + key
}

func InitCacheClient(conf *RedisConfig) Cache {
	DefaultRedisClient = &RedisCacheClient{
		status: newCacheStat(),
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
	if byteValue, err = json.Marshal(value); err != nil {
		logger.Error("set redis key: %v, error: %v", fullKey, err)
		return err
	}
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


func (r *RedisCacheClient) Get(ctx context.Context, key string, result interface{}) (exist bool, err error) {
	var byteValue []byte
	fullKey := getFullKey(r.prefix, key)
	startTime := time.Now()
	byteValue, err = r.client.Get(fullKey).Bytes()
	elapsed := time.Since(startTime).Milliseconds()
	for _, p := range r.plugins {
		p.OnGetRequestEnd(ctx, cmdGet, elapsed, fullKey, err)
	}
	// not found key
	if err == redis.Nil {
		r.status.IncrementMiss()
		return false, nil
	}
	// something err get key from redis
	if err != nil {
		logger.Error("get redis key: %v, error: %v", fullKey, err)
		return false, err
	}
	r.status.IncrementHit()
	return true, json.Unmarshal([]byte(byteValue), result)
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
