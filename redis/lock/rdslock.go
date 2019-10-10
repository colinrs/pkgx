package lock

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

// RdsLock ...
type RdsLock struct {
	conn redis.Conn
}

// NewLockWithPasswd ...
func NewLockWithPasswd(dial, pswd string) (resLock *RdsLock, err error) {
	rpswd := redis.DialPassword(pswd)
	var rdsc redis.Conn
	rdsc, err = redis.Dial("tcp", dial, rpswd)
	if err != nil {
		return
	}
	resLock = &RdsLock{
		conn: rdsc,
	}
	return
}

// NewLock ...
func NewLock(dial string) (resLock *RdsLock, err error) {
	var rdsc redis.Conn
	rdsc, err = redis.Dial("tcp", dial)
	if err != nil {
		return
	}
	resLock = &RdsLock{
		conn: rdsc,
	}
	return
}

// Lock ...
type Lock struct {
	resource string
	token    string
	conn     redis.Conn
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
