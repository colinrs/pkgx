package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	defaultHTTPTimeout = 10
	httpMethodGet      = "GET"
	httpMethodPost     = "POST"
)

type Client interface {
	Get(ctx context.Context, url string, param map[string]string, reply interface{}, opts ...RequestOption) error
	Post(ctx context.Context, url string, body interface{}, reply interface{}, opts ...RequestOption) error
	AddPlugin(p Plugin)
}

type ClientImpl struct {
	client  *http.Client
	options options
	plugins []Plugin
}

// Determine whether ClientImpl implements Client
var _ Client = (*ClientImpl)(nil)

var C Client

func (c *ClientImpl) Post(ctx context.Context, url string, body interface{}, reply interface{}, opts ...RequestOption) (err error) {
	reader, err := bodyReader(body)
	if err != nil {
		return err
	}
	return c.handle(ctx, httpMethodPost, url, reader, reply, opts...)
}

func bodyReader(body interface{}) (io.Reader, error) {
	if _, ok := body.(io.Reader); ok {
		return body.(io.Reader), nil
	}
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(buf), nil
}

func (c *ClientImpl) Get(ctx context.Context, url string, param map[string]string, reply interface{}, opts ...RequestOption) (err error) {
	var paramList []string
	for k, v := range param {
		paramList = append(paramList, k+"="+v)
	}
	if len(paramList) > 0 {
		url += "?" + strings.Join(paramList, "&")
	}
	return c.handle(ctx, httpMethodGet, url, nil, reply, opts...)
}

func (c *ClientImpl) do(ctx context.Context, method string, url string, reader io.Reader, header http.Header) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, c.getHTTPTimeout(ctx, getPathKey(req.URL.Path)))
	defer cancel()
	req = req.WithContext(ctx)
	c.reportRequestStart(ctx, req)
	startTime := time.Now()
	req.Header = header
	resp, err = c.client.Do(req)
	c.reportRequestEnd(ctx, req, resp, err, startTime)
	return resp, err
}

func (c *ClientImpl) handle(ctx context.Context, method string, url string, reader io.Reader, reply interface{}, opts ...RequestOption) (err error) {
	opt := RequestOptions{
		ContentType: "application/json",
		Header:      http.Header{},
	}
	for _, o := range opts {
		o(&opt)
	}

	opt.Header.Set("Content-Type", opt.ContentType)
	resp, err := c.do(ctx, method, url, reader, opt.Header)
	if err != nil {
		return err
	}
	if opt.RespHandler != nil {
		return opt.RespHandler(resp)
	}

	if resp.StatusCode == http.StatusNoContent {
		errNoContent := fmt.Errorf("http code:%v url:%v", resp.StatusCode, url)
		return errNoContent
	}

	if resp.StatusCode != http.StatusOK {
		errStatusNotOK := fmt.Errorf("http code:%v url:%v", resp.StatusCode, url)
		return errStatusNotOK
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, reply)
}

func (c *ClientImpl) getHTTPTimeout(ctx context.Context, path string) time.Duration {
	return time.Duration(defaultHTTPTimeout) * time.Second
}

func (c *ClientImpl) reportRequestStart(ctx context.Context, request *http.Request) {
	for _, plugin := range c.plugins {
		plugin.OnRequestStart(ctx, request)
	}
}

func (c *ClientImpl) reportRequestEnd(ctx context.Context, request *http.Request, response *http.Response, err error, startTime time.Time) {
	for _, plugin := range c.plugins {
		plugin.OnRequestEnd(ctx, request, response, startTime, err)
	}
}

func (c *ClientImpl) AddPlugin(p Plugin) {
	c.plugins = append(c.plugins, p)
}

func InitClient(opts ...Option) Client {
	opt := options{
		maxConnectionNum:    100,
		timeout:             60 * time.Second,
		dialTimeout:         30 * time.Second,
		idleConnTimeout:     90 * time.Second,
		keepAlive:           30 * time.Second,
		tlsHandshakeTimeout: 10 * time.Second,
	}
	for _, o := range opts {
		o(&opt)
	}
	C = &ClientImpl{
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   opt.dialTimeout,
					KeepAlive: opt.keepAlive,
				}).DialContext,
				MaxIdleConns:          opt.maxConnectionNum,
				MaxIdleConnsPerHost:   opt.maxConnectionNum,
				IdleConnTimeout:       opt.idleConnTimeout,
				TLSHandshakeTimeout:   opt.tlsHandshakeTimeout,
				ExpectContinueTimeout: 1 * time.Second,
			},
			Timeout: opt.timeout,
		},
		options: opt,
	}
	return C
}

func getPathKey(path string) string {
	s := strings.Split(path, "/")
	if len(s) > 0 {
		return s[len(s)-1]
	}
	return "emptyPathKey"
}
