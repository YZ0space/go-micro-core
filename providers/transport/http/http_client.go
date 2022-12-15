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
	"strconv"
	"strings"
	"time"
)

type HttpClient struct {
	client  *http.Client
	options options
}

func (h *HttpClient) Post(ctx context.Context, url string, body interface{}, v interface{}, opts ...RequestOption) (err error) {
	return h.handle(ctx, "POST", url, body, v, opts...)
}

func (h *HttpClient) handle(ctx context.Context, method string, url string, reqBody interface{}, v interface{}, opts ...RequestOption) (err error) {
	opt := RequestOptions{
		ContentType: "application/json",
		Header:      http.Header{},
	}
	for _, o := range opts {
		o(&opt)
	}

	opt.Header.Set("Content-Type", opt.ContentType)

	retryTimesStr, retryIntervalStr := opt.Header.Get("RETRY-TIMES"), opt.Header.Get("RETRY-INTERVAL")
	retryTimes, retryInterval := 1, 0
	if retryTimesStr != "" {
		retryTimes, _ = strconv.Atoi(retryTimesStr)
		opt.Header.Del("RETRY-TIMES")
	}
	if retryIntervalStr != "" {
		retryInterval, _ = strconv.Atoi(retryIntervalStr)
		opt.Header.Del("RETRY-INTERVAL")
	}

	resp, err := h.doWithRetry(ctx, method, url, reqBody, opt.Header, retryTimes, retryInterval)
	if err != nil {
		return
	}

	if opt.RespHandler != nil {
		return opt.RespHandler(resp)
	}

	if resp.StatusCode == http.StatusNoContent {
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, tempErr := ioutil.ReadAll(resp.Body)
		if tempErr != nil {
			return fmt.Errorf("http code:%v url:%v", resp.StatusCode, url)
		}
		_ = json.Unmarshal(body, v)
		return fmt.Errorf("http code:%v url:%v", resp.StatusCode, url)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, v)
	return
}

func (h *HttpClient) Get(ctx context.Context, url string, v interface{}, opts ...RequestOption) (err error) {
	return h.handle(ctx, "GET", url, nil, v, opts...)
}

func (h *HttpClient) do(method string, url string, reader io.Reader, header http.Header) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}

	req.Header = header
	return h.client.Do(req)
}

func (h *HttpClient) doWithRetry(ctx context.Context, method string, url string, body interface{}, header http.Header, retryTimes, retryInterval int) (resp *http.Response, err error) {
	defer func() {
		if err != nil && strings.Contains(err.Error(), "reset by peer") && retryTimes > 0 {
			if retryInterval != 0 {
				time.Sleep(time.Duration(retryInterval * int(time.Second)))
			}
			resp, err = h.doWithRetry(ctx, method, url, body, header, retryTimes-1, retryInterval)
		}
		return
	}()

	//io.reader不能被多次读取，因此使用body新建reader；GET请求reader为nil
	var reader io.Reader
	if method == "POST" {
		reader, err = h.bodyReader(body)
		if err != nil {
			return
		}
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}

	req.Header = header

	resp, err = h.client.Do(req)
	return
}

func (h *HttpClient) bodyReader(body interface{}) (io.Reader, error) {
	if _, ok := body.(io.Reader); ok {
		return body.(io.Reader), nil
	}

	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(buf), nil
}

func NewHttpClient(opts ...Option) *HttpClient {
	// anyway  先默认吧 不走配置了
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

	return &HttpClient{
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
}
