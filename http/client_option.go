package http

import (
	"net/http"
	"time"
)

type options struct {
	maxConnectionNum    int
	timeout             time.Duration
	dialTimeout         time.Duration
	idleConnTimeout     time.Duration
	keepAlive           time.Duration
	tlsHandshakeTimeout time.Duration
	name                string
}

type Option func(*options)

func WithMaxConnectionNum(maxConnectionNum int) Option {
	return func(opt *options) {
		opt.maxConnectionNum = maxConnectionNum
	}
}

func WithTimeout(duration time.Duration) Option {
	return func(opt *options) {
		opt.timeout = duration
	}
}

func WithDialTimeout(duration time.Duration) Option {
	return func(opt *options) {
		opt.dialTimeout = duration
	}
}

func WithIdleConnTimeout(duration time.Duration) Option {
	return func(opt *options) {
		opt.idleConnTimeout = duration
	}
}

func WithKeepAlive(duration time.Duration) Option {
	return func(opt *options) {
		opt.keepAlive = duration
	}
}

func WithTLSHandshakeTimeout(duration time.Duration) Option {
	return func(opt *options) {
		opt.tlsHandshakeTimeout = duration
	}
}

func WithServiceName(name string) Option {
	return func(opt *options) {
		opt.name = name
	}
}

type RequestOptions struct {
	ContentType string
	Header      http.Header
	RespHandler func(*http.Response) error
}

type RequestOption func(*RequestOptions)

func WithContentType(contentType string) RequestOption {
	return func(opt *RequestOptions) {
		opt.ContentType = contentType
	}
}

func WithHeader(m http.Header) RequestOption {
	return func(opt *RequestOptions) {
		opt.Header = m
	}
}

func WithHeaderSet(kv ...string) RequestOption {
	return func(opt *RequestOptions) {
		for i := 0; i < len(kv)-1; i += 2 {
			opt.Header.Set(kv[i], kv[i+1])
		}
	}
}

func WithResponse(fn func(*http.Response) error) RequestOption {
	return func(opt *RequestOptions) {
		opt.RespHandler = fn
	}
}
