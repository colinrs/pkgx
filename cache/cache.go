package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/colinrs/pkgx/logger"
	"github.com/go-redis/redis"
)

const (
	cmdGet    = "get"
	cmdSet    = "set"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error)
	Get(ctx context.Context, key string, result interface{}) (exist bool, err error)
	AddPlugin(p Plugin)
}

type RedisConfig struct {
	Addr        string `yaml:"addr" json:"addr"`
	DB          int    `yaml:"db" json:"db"`
	PoolSize    int    `yaml:"pool_size" json:"pool_size"`
	IdleTimeout int    `yaml:"idle_timeout" json:"idle_timeout"`
	Prefix      string `yaml:"prefix" json:"prefix"`
}

type RedisClient struct {
	client  *redis.Client
	prefix  string
	plugins []Plugin
}

var _ Cache = (*RedisClient)(nil)

var DefaultRedisClient *RedisClient

func getFullKey(prefix, key string) string {
	return prefix + "_" + key
}

func InitCacheClient(conf *RedisConfig) Cache {
	DefaultRedisClient = &RedisClient{}
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

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error) {
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


func (r *RedisClient) Get(ctx context.Context, key string, result interface{}) (exist bool, err error) {
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
		return false, nil
	}
	// something err get key from redis
	if err != nil {
		logger.Error("get redis key: %v, error: %v", fullKey, err)
		return false, err
	}
	return true, json.Unmarshal([]byte(byteValue), result)
}

func (r *RedisClient) AddPlugin(p Plugin) {
	r.plugins = append(r.plugins, p)
}
