package api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-common-framework/services/http"
	"go-common-framework/services/log"
)

const (
	//服务
	SERVICE_API = "service-api"

	//参数类型
	PARAMS_TYPE_QUERY      = "query"
	PARAMS_TYPE_JSON       = "json"
	PARAMS_TYPE_FORM       = "form"
	PARAMS_TYPE_URLENCODED = "urlencoded"
)

var (
	ApiMap = make(map[string]IBaseApi)
	Locker sync.Mutex
)

type IBaseApi interface {
	ServiceConfig() error
	GetBaseApi() *BaseApi
}

type BaseApi struct {
	Apis           map[string]*UrlPath
	Host           string
	Domain         string
	Config         map[string]string
	ConnectTimeOut time.Duration
	Options        *Options
}

type UrlPath struct {
	Method string
	Path   string
}

type Options struct {
	Attributes map[string]interface{}
	Data       map[string]IBaseApiRequest
	Headers    map[string]string
}

func GetApi(name string) (IBaseApi, error) {
	api, ok := ApiMap[name]
	if !ok {
		return nil, fmt.Errorf("unknown api %s", name)
	}
	return api, nil
}

func (a *BaseApi) getURL(apiName string) (*UrlPath, error) {
	if urlPath, ok := a.Apis[apiName]; ok {
		return urlPath, nil
	}
	return nil, fmt.Errorf("%s apiname not found", apiName)
}

func Send(requestId string, a IBaseApi, apiName string) ([]byte, error) {
	base := a.GetBaseApi()
	urlPath, err := base.getURL(apiName)
	if err != nil {
		return nil, err
	}
	//url中的替换参数
	u := urlPath.Path
	if len(base.Options.Attributes) > 0 {
		for k, v := range base.Options.Attributes {
			switch v.(type) {
			case string:
				u = strings.Replace(u, "{"+k+"}", v.(string), 1)
			case int:
				u = strings.Replace(u, "{"+k+"}", strconv.Itoa(v.(int)), 1)
			}
		}
	}

	params := make(map[string]string)
	body := ""

	if base.Options.Data != nil {
		//处理query参数
		if data, ok := base.Options.Data[PARAMS_TYPE_QUERY]; ok {
			query, err := data.Decode()
			if err != nil {
				return nil, err
			}
			//利用反射处理不同的结构体变量类型
			typeData := reflect.TypeOf(data).Elem()
			for k, v := range query {
				switch v.(type) {
				case json.Number:
					for i := 0; i < typeData.NumField(); i++ {
						tag := typeData.Field(i).Tag.Get("json")
						ret := strings.Split(tag, ",")
						if ret[0] == k {
							dataType := typeData.Field(i).Type.Name()
							if dataType == "int64" || dataType == "int32" || dataType == "int" {
								res, _ := v.(json.Number).Int64()
								params[k] = strconv.Itoa(int(res))
							} else if dataType == "float64" || dataType == "float32" {
								res, _ := v.(json.Number).Float64()
								params[k] = fmt.Sprintf("%.2f", res)
							} else if dataType == "string" {
								params[k] = v.(json.Number).String()
							}
							break
						}
					}
				case string:
					params[k] = v.(string)
				}
			}
		}
		//处理application/json参数
		if data, ok := base.Options.Data[PARAMS_TYPE_JSON]; ok {
			bodyData, err := json.Marshal(data)
			if err != nil {
				return nil, err
			}
			body = string(bodyData)
			base.Options.Headers["Content-Type"] = "application/json"
		} else if data, ok = base.Options.Data[PARAMS_TYPE_URLENCODED]; ok { //处理application/x-www-form-urlencoded参数
			urlencoded, err := data.Decode()
			if err != nil {
				return nil, err
			}
			for k, v := range urlencoded {
				body += k + "=" + fmt.Sprintf("%s", v) + "&"
			}
			body = strings.Trim(body, "&")
			base.Options.Headers["Content-Type"] = "application/x-www-form-urlencoded"
		} else if data, ok = base.Options.Data[PARAMS_TYPE_FORM]; ok { //处理multipart/form-data参数
			form, err := data.Decode()
			if err != nil {
				return nil, err
			}
			for k, v := range form {
				body += k + "=" + fmt.Sprintf("%s", v) + ";"
			}
			body = strings.Trim(body, ";")
			base.Options.Headers["Content-Type"] = "multipart/form-data"
		}
	}

	reqURL := base.Domain + u
	base.Options.Headers["X-Request-Id"] = requestId
	Locker.Unlock()
	response, err := http.DoHttpRequest(
		reqURL,
		urlPath.Method,
		body,
		base.Host,
		base.ConnectTimeOut,
		base.Options.Headers,
		params,
	)
	if err != nil {
		log.Error(requestId, "do http request error:%v, params:%v", err, params)
		return nil, err
	}

	log.Info(requestId, "do http response:%s", string(response))
	return response, nil
}
