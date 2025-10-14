package asynqueue

import (
	"errors"
	"fmt"
)

var (
	ErrAsynqAddrEmpty = errors.New("asynq addr is empty")
)

type AsynqConfig struct {
	Addr        string
	Concurrency int
	DB          int `json:",default=0,optional"` // default 0
	Queue       string
	Monitoring  *Monitoring `json:",optional"`
}

type Monitoring struct {
	Enable   bool   `json:",default=false,optional"`
	Path     string `json:",default=/monitoring,optional"` // default /monitoring
	Port     string `json:",default=8089,optional"`        // default 8089
	UserName string `json:",default=userName,optional"`
	Password string `json:",default=password,optional"`
}

func (c *AsynqConfig) Validate() error {
	if c.Addr == "" {
		return ErrAsynqAddrEmpty
	}
	if c.Concurrency <= 0 {
		c.Concurrency = defaultConcurrency
	}
	if c.Queue == "" {
		fmt.Printf("asynq queue is empty, use default queue: %s\n", defaultQueue)
		c.Queue = defaultQueue
	}
	if c.Monitoring == nil {
		return nil
	}
	if !c.Monitoring.Enable {
		return nil
	}
	if c.Monitoring.Path == "" {
		c.Monitoring.Path = defaultPath
	}
	if c.Monitoring.Port == "" {
		c.Monitoring.Port = defaultPort
	}
	if c.Monitoring.UserName == "" {
		c.Monitoring.UserName = defaultUserName
	}
	if c.Monitoring.Password == "" {
		c.Monitoring.Password = defaultPassword
	}
	return nil
}
