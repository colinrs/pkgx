package lock

import (
	"fmt"
	"time"

	"github.com/colinrs/pkgx/cache"
	"github.com/redis/go-redis/v9"
)

// RdsLock ...
type RdsLock struct {
	conn cache.Cache
}

// NewLockWithPasswd ...
func NewLockWithPasswd(dial, username, pswd string) (resLock *RdsLock, err error) {
	config := cache.RedisConfig{
		Addr:     dial,
		Username: username,
		Password: pswd,
	}
	resLock = &RdsLock{
		conn: cache.InitCacheClient(&config),
	}
	return
}

// NewLock ...
func NewLock(dial string) (resLock *RdsLock, err error) {
	config := cache.RedisConfig{
		Addr: dial,
	}
	resLock = &RdsLock{
		conn: cache.InitCacheClient(&config),
	}
	return
}

// Lock ...
type Lock struct {
	resource string
	token    string
	c        cache.Cache
	timeout  int
}

// TryLock ...
func (p *RdsLock) TryLock(resource string, token string, timeout int) (lock *Lock, err error) {
	lock = &Lock{resource, token, p.conn, timeout}
	var ok bool
	ok, err = lock.tryLock()
	if !ok || err != nil {
		lock = nil
	}
	return
}

func (lock *Lock) tryLock() (ok bool, err error) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		_, err = redis.String(lock.conn.Do("SET", lock.key(), lock.token, "EX", int(lock.timeout), "NX"))
		if err == redis.ErrNil {
			// The lock was not successful, it already exists.
			<-ticker.C
			continue
		}
		if err != nil {
			<-ticker.C
			continue
		}
		ok = true
		fmt.Println("lock sucess")
		break
	}
	f := func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			<-ticker.C
			lock.AddTimeout(1)
		}

	}
	go f()
	return
}

// Unlock ....
func (lock *Lock) Unlock() (err error) {
	_, err = lock.conn.Do("del", lock.key())
	return
}

func (lock *Lock) key() string {
	return fmt.Sprintf("redislock:%s", lock.resource)
}

// AddTimeout ...
func (lock *Lock) AddTimeout(exTime int64) (ok bool, err error) {
	var ttlTime int64
	ttlTime, err = redis.Int64(lock.conn.Do("TTL", lock.key()))
	if err != nil {
		return
	}
	fmt.Println(ttlTime)
	if ttlTime > 0 {
		_, err = redis.String(lock.conn.Do("SET", lock.key(), lock.token, "EX", int(ttlTime+exTime)))
		if err == redis.ErrNil {
			return
		}
		if err != nil {
			return
		}
	}
	return
}
