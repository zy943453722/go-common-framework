package http

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/httplib"
	"go-common-framework/services/log"
)

const (
	RETRY_TIMES = 3 //失败重试3次
	RETRY_DELAY = 1 //重试延迟1s
)

func DoHttpRequest(reqUrl, method, bodyStr, host string, timeout time.Duration, header, params map[string]string) ([]byte, error) {
	if len(params) > 0 {
		values := url.Values{}
		for k, v := range params {
			values.Set(k, v)
		}
		u, err := url.ParseRequestURI(reqUrl)
		if err != nil {
			return nil, err
		}
		u.RawQuery = values.Encode()
		reqUrl = fmt.Sprintf("%v", u)
	}
	req := httplib.NewBeegoRequest(reqUrl, method)
	req.Retries(RETRY_TIMES) //接口调用重试
	req.RetryDelay(time.Duration(RETRY_DELAY))
	req.SetFilters(LogFilter)
	req.SetTimeout(timeout, timeout)
	for k, v := range header {
		req.Header(k, v)
	}
	if host != "" {
		req.SetHost(host)
	}
	if bodyStr != "" {
		req.Body(bodyStr)
	}
	return req.Bytes()
}

func LogFilter(next httplib.Filter) httplib.Filter {
	return func(ctx context.Context, req *httplib.BeegoHTTPRequest) (*http.Response, error) {
		r := req.GetRequest()
		cmd := make([]string, 0)
		cmd = append(cmd, "curl", "-X", r.Method)
		for key, values := range r.Header {
			for _, value := range values {
				cmd = append(cmd, "-H", fmt.Sprintf("'%s:%s'", key, value))
			}
		}
		if r.Host != "" {
			cmd = append(cmd, "-H", fmt.Sprintf("'%s:%s'", "Host", r.Host))
		}

		if contentType, ok := r.Header["Content-Type"]; ok {
			if len(contentType) > 0 && r.Body != nil {
				//为了防止r.Body读完后文件指针指向最后，所以需要重新赋值
				body, _ := ioutil.ReadAll(r.Body)
				r.Body = ioutil.NopCloser(bytes.NewReader(body))
				cmd = append(cmd, "-d", "'"+string(body)+"'")
			}
		}

		reqUrl, _ := url.QueryUnescape(r.URL.String())
		cmd = append(cmd, reqUrl)
		curl := strings.Join(cmd, " ")
		log.Info(r.Header.Get("X-Request-Id"), curl)
		return next(ctx, req)
	}
}
