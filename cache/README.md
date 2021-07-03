# cache

```go
package main

import (
	"context"
	"time"

	"github.com/colinrs/pkgx/cache"
	"github.com/colinrs/pkgx/logger"
)



type Student struct {
	Name string `json:"name"`
	Age int `json:"age"`

}

func main(){
	// 配置
	conf := &cache.RedisConfig{
		Addr:           "127.0.0.1:6379",
		DB:             0,
		PoolSize:       100,
		IdleTimeout:    3600,
		Prefix:         "test",
		LocalCacheSize: 1,
	}
	// 缓存实例
	redisCache := cache.InitCacheClient(conf)
	ctx := context.Background()
	var err error
	// 设置缓存+过期时间
	err = redisCache.Set(ctx, "key1", "val1", 100*time.Second)
	if err!=nil{
		logger.Error("set cache err:%s", err.Error())
	}
	var result []byte
	// 获取缓存，从原始数据获取设为nil
	result,err = redisCache.Get(ctx, "key1", nil)
	if err!=nil{
		logger.Error("get cache err:%s", err.Error())
	}
	if result!=nil{
		logger.Info("getResult:%s", result)
	}
	s := Student{
		Name: "s1",
		Age: 22,
	}
	err = redisCache.Set(ctx, "s1", s, 100*time.Second)
	if err!=nil{
		logger.Error("set cache err:%s", err.Error())
	}
	result,err = redisCache.Get(ctx, "s1", nil)
	if err!=nil{
		logger.Error("get cache err:%s", err.Error())
	}
	if result!=nil{
		logger.Info("getResult s1:%s", result)
	}
	result,err = redisCache.Get(ctx, "s1", nil)
	if err!=nil{
		logger.Error("get cache err:%s", err.Error())
	}
	if result!=nil{
		logger.Info("getResult s1:%s", result)
	}
	time.Sleep(10000*time.Second)
}
```
