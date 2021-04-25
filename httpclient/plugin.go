package httpclient

import (
	"context"
	"net/http"
	"time"
)

type Plugin interface {
	OnRequestStart(context.Context, *http.Request)
	OnRequestEnd(context.Context, *http.Request, *http.Response, time.Time, error)
}

var _ Plugin = (*defaultHTTPPlugin)(nil)

type defaultHTTPPlugin struct {
}

var DefaultHTTPPlugin = &defaultHTTPPlugin{}

func (d *defaultHTTPPlugin) OnRequestStart(ctx context.Context, req *http.Request) {
}

func (d *defaultHTTPPlugin) OnRequestEnd(ctx context.Context, req *http.Request, resp *http.Response, startTime time.Time, err error) {
}
