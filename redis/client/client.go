package client

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisConfig struct {
	Addr           string `yaml:"addr" json:"addr"`
	DB             int    `yaml:"db" json:"db"`
	PoolSize       int    `yaml:"pool_size" json:"pool_size"`
	ReadTimeout    int    `yaml:"read_timeout" json:"read_timeout"`   // ms
	WriteTimeout   int    `yaml:"write_timeout" json:"write_timeout"` // s
	IdleTimeout    int    `yaml:"idle_timeout" json:"idle_timeout"`
	Prefix         string `yaml:"prefix" json:"prefix"`
	LocalCacheSize int    `yaml:"local_cache_size" json:"local_cache_size"` // M
}

func New(config *RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		ReadTimeout:  time.Duration(config.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(config.ReadTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.IdleTimeout) * time.Second,
	})
}
