package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RdsLock ...
type RdsLock struct {
	conn Cache
}

// NewLockWithPasswd ...
func NewLockWithPasswd(dial, username, pswd string) (resLock *RdsLock, err error) {
	config := RedisConfig{
		Addr:     dial,
		Username: username,
		Password: pswd,
	}
	resLock = &RdsLock{
		conn: InitCacheClient(&config),
	}
	return
}

// NewLock ...
func NewLock(dial string) (resLock *RdsLock, err error) {
	config := RedisConfig{
		Addr: dial,
	}
	resLock = &RdsLock{
		conn: InitCacheClient(&config),
	}
	return
}

// Lock ...
type Lock struct {
	resource string
	token    string
	c        Cache
	timeout  time.Duration // sec
}

// TryLock ...
func (p *RdsLock) TryLock(ctx context.Context, resource string, token string, timeout int) (lock *Lock, err error) {
	lock = &Lock{resource, token, p.conn, time.Duration(timeout) * time.Millisecond}
	var ok bool
	ok, err = lock.tryLock(ctx)
	if !ok || err != nil {
		lock = nil
	}
	return
}

func (lock *Lock) tryLock(ctx context.Context) (ok bool, err error) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		err = lock.c.SetNX(ctx, lock.key(), lock.token, lock.timeout)
		if err == redis.Nil {
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
			lock.AddTimeout(ctx, time.Second.Milliseconds())
		}
	}
	go f()
	return
}

// Unlock ....
func (lock *Lock) Unlock(ctx context.Context) (err error) {
	return lock.c.Del(ctx, lock.key())
}

func (lock *Lock) key() string {
	return fmt.Sprintf("redislock:%s", lock.resource)
}

// AddTimeout ...
func (lock *Lock) AddTimeout(ctx context.Context, exTime int64) (ok bool, err error) {
	ttlTime, err := lock.c.TTL(ctx, lock.key())
	if err != nil {
		return false, err
	}
	if ttlTime.Milliseconds() > 0 {
		ok, err = lock.c.Expire(ctx, lock.key(), time.Duration(int(ttlTime.Milliseconds()+exTime))*time.Millisecond)
		if err == redis.Nil {
			return false, nil
		}
		return ok, err
	}
	return false, nil
}
