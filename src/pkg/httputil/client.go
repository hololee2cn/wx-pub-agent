package httputil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hololee2cn/pkg/ginx"

	log "github.com/sirupsen/logrus"
)

const (
	XCode = "x-code"
	XMsg  = "x-msg"
)

type CustomErrMsg struct {
	XCode int    `json:"x-code"`
	XMsg  string `json:"x-msg"`
}

type RequestProperty struct {
	Method  string
	URI     string
	Payload []byte
	Header  map[string]string
}

var (
	DefaultHTTPClient *http.Client
)

func init() {
	DefaultHTTPClient = &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSHandshakeTimeout: time.Second * 10,
			DisableKeepAlives:   false,
			DisableCompression:  false,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
			IdleConnTimeout:     10 * time.Second,
		},
	}
}

func request(ctx context.Context, client *http.Client, req RequestProperty) (int, []byte, http.Header, error) {
	method, uri, payload, header := req.Method, req.URI, req.Payload, req.Header

	// 判断url是否有效
	_, err := url.ParseRequestURI(uri)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("invalid request url:%s, error: %+v", uri, err)
	}

	request, err := http.NewRequest(method, uri, bytes.NewReader(payload))
	if ctx != nil {
		request, err = http.NewRequestWithContext(ctx, method, uri, bytes.NewReader(payload))
	}

	if err != nil {
		return 0, nil, nil, fmt.Errorf("new request failed, url:%s, error: %+v", uri, err)
	}

	for key, value := range header {
		request.Header.Set(key, value)
	}

	var resp *http.Response
	for i := 0; i < 3; i++ {
		resp, err = client.Do(request)
		if err != nil {
			// reset Request.Body
			request.Body = ioutil.NopCloser(bytes.NewReader(payload))
			time.Sleep(time.Millisecond * 10)
			continue
		}
		break
	}
	if err != nil {
		return 0, nil, nil, fmt.Errorf("failed to send request, method:%s, url:%s, err: %+v", method, uri, err)
	}
	defer func() {
		resp.Body.Close()
	}()

	// 南凌内部服务错误获取方式
	errCode := resp.Header.Get(XCode)
	errMsg := resp.Header.Get(XMsg)
	if len(errCode) > 0 || len(errMsg) > 0 {
		customErr, err := GetNovaServiceErrorResponseFromHeader(resp.Header)
		if err != nil {
			return 0, nil, resp.Header, fmt.Errorf("request errCode:%s strconv atoi failed, method:%s, url:%s, err: %+v", errCode, method, uri, err)
		}
		if customErr != nil {
			body, err := json.Marshal(customErr)
			if err != nil {
				return 0, nil, nil, fmt.Errorf("request custom error marshal failed, method:%s, url:%s, err: %+v", method, uri, err)
			}
			return resp.StatusCode, body, resp.Header, nil
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, resp.Header, fmt.Errorf("failed to read response payload, method:%s, url:%s, err: %+v", method, uri, err)
	}

	return resp.StatusCode, body, resp.Header, nil
}

func IsStatusCodeOK(statusCode int) bool {
	if statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices {
		return true
	}
	return false
}

func RequestWithRepeat(traceID string, method, uri string, payload []byte, header map[string]string) (statusCode int, body []byte, respHeader http.Header, err error) {
	log.Debugf("http request, traceID:%s, requests url: %s, method: %s", traceID, uri, method)
	var start = time.Now()

	defer func() {
		log.Debugf("get http request response, traceID:%s, statusCode:%d, use time: %d, err: %+v", traceID, statusCode, time.Since(start).Milliseconds(), err)
	}()
	if len(header[ginx.HTTPTraceIDHeader]) == 0 {
		header[ginx.HTTPTraceIDHeader] = traceID
	}
	statusCode, body, respHeader, err = request(context.TODO(), DefaultHTTPClient, GetRequestProperty(method, uri, payload, header))
	return
}

func RequestWithContextAndRepeat(ctx context.Context, req RequestProperty, traceID string) (statusCode int, body []byte, respHeader http.Header, err error) {
	log.Debugf("http context request, traceID:%s, requests url: %s, method: %s", traceID, req.URI, req.Method)
	var start = time.Now()

	defer func() {
		log.Debugf("get http context request response, traceID:%s, statusCode:%d, use time: %d, err: %+v", traceID, statusCode, time.Since(start).Milliseconds(), err)
	}()
	if len(req.Header[ginx.HTTPTraceIDHeader]) == 0 {
		req.Header[ginx.HTTPTraceIDHeader] = traceID
	}
	statusCode, body, respHeader, err = request(ctx, DefaultHTTPClient, req)
	return
}

func GetRequestProperty(method, uri string, payload []byte, header map[string]string) RequestProperty {
	return RequestProperty{
		Method:  method,
		URI:     uri,
		Payload: payload,
		Header:  header,
	}
}

func GetNovaServiceErrorResponseFromHeader(header http.Header) (*CustomErrMsg, error) {
	// 南凌内部服务错误获取方式
	errCode := header.Get(XCode)
	errMsg := header.Get(XMsg)
	if len(errCode) == 0 {
		return nil, fmt.Errorf("response x-code header is not exist")
	}
	code, err := strconv.Atoi(errCode)
	if err != nil {
		return nil, err
	}
	if code == 0 {
		return nil, nil
	}
	customErr := CustomErrMsg{
		XCode: code,
		XMsg:  errMsg,
	}
	return &customErr, nil
}
