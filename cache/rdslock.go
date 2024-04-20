package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RdsLock ...
type RdsLock struct {
	conn Cache
}

// NewLockWithPasswd ...
func NewLockWithPasswd(dial, username, pswd string) (*RdsLock, error) {
	if dial == "" {
		return nil, errors.New("dial can not be empty")
	}
	if username == "" {
		return nil, errors.New("username can not be empty")
	}
	if pswd == "" {
		return nil, errors.New("pswd can not be empty")
	}
	config := RedisConfig{
		Addr:     dial,
		Username: username,
		Password: pswd,
	}
	conn := InitCacheClient(&config)
	if conn == nil {
		return nil, fmt.Errorf("failed to init cache client")
	}
	return &RdsLock{
		conn: conn,
	}, nil
}

// NewLock ...
func NewLock(dial string) (*RdsLock, error) {
	if dial == "" {
		return nil, errors.New("dial can not be empty")
	}
	config := RedisConfig{
		Addr: dial,
	}
	conn := InitCacheClient(&config)
	if conn == nil {
		return nil, fmt.Errorf("failed to init cache client")
	}
	return &RdsLock{
		conn: conn,
	}, nil
}

// Lock ...
type Lock struct {
	resource string
	token    string
	c        Cache
	timeout  time.Duration // sec
	done     chan struct{}
}

// TryLock ...
func (p *RdsLock) TryLock(ctx context.Context, resource string, token string, timeout int) (lock *Lock, err error) {
	if p.conn == nil {
		return nil, fmt.Errorf("RdsLock is not initialized")
	}
	if resource == "" {
		return nil, fmt.Errorf("resource can not be empty")
	}
	if token == "" {
		return nil, fmt.Errorf("token can not be empty")
	}
	if timeout <= 0 {
		return nil, fmt.Errorf("timeout must be greater than 0")
	}
	lock = &Lock{
		resource: resource,
		token:    token,
		c:        p.conn,
		timeout:  time.Duration(timeout) * time.Millisecond,
		done:     make(chan struct{}, 1),
	}
	ok, err := lock.tryLock(ctx)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("failed to acquire lock")
	}
	return lock, nil
}

func (lock *Lock) tryLock(ctx context.Context) (ok bool, err error) {
	for {
		err = lock.c.SetNX(ctx, lock.key(), lock.token, lock.timeout)
		if errors.Is(err, redis.Nil) {
			return false, err
		}
		if err != nil {
			return false, err
		}
		ok = true
		break
	}

	go func() {
		ticker := time.NewTicker(lock.timeout / 2 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-lock.done:
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				_, err := lock.AddTimeout(ctx, time.Second.Milliseconds())
				if err != nil {
					return
				}
			}
		}
	}()

	return
}

// Unlock ...
func (lock *Lock) Unlock(ctx context.Context) (err error) {
	if lock == nil || lock.done == nil || lock.c == nil {
		return fmt.Errorf("lock, lock.done or lock.c is nil")
	}
	close(lock.done)
	return lock.c.Del(ctx, lock.key())
}

func (lock *Lock) key() string {
	if lock == nil || lock.resource == "" {
		panic("lock or lock.resource is nil")
	}
	return fmt.Sprintf("redislock:%s", lock.resource)
}

// AddTimeout ...
func (lock *Lock) AddTimeout(ctx context.Context, exTime int64) (ok bool, err error) {
	if lock == nil {
		return false, fmt.Errorf("lock is nil")
	}
	if lock.c == nil {
		return false, fmt.Errorf("lock.c is nil")
	}

	ttlTime, err := lock.c.TTL(ctx, lock.key())
	if err != nil {
		return false, err
	}

	totalTime := ttlTime.Milliseconds() + exTime
	if totalTime < 0 {
		return false, fmt.Errorf("totalTime is less than 0")
	}

	expireDuration := time.Duration(totalTime) * time.Millisecond
	ok, err = lock.c.Expire(ctx, lock.key(), expireDuration)
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	return ok, err
}
