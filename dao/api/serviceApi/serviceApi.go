package serviceApi

import (
	"fmt"
	"os"
	"time"

	"github.com/beego/beego/v2/server/web"
	"go-common-framework/dao/api"
)

var apis = map[string]*api.UrlPath{
	"getDeviceList": {
		Method: "GET",
		Path:   "/v1/device/getDeviceList",
	},
}

type ServiceApi struct {
	api.BaseApi
}

func (a *ServiceApi) ServiceConfig() error {
	var err error
	a.Host, err = web.AppConfig.String("service.host")
	if err != nil {
		return err
	}
	a.Domain, err = web.AppConfig.String("service.domain")
	if err != nil {
		return err
	}
	timeout, err := web.AppConfig.Int("service.timeout")
	if err != nil {
		return err
	}
	a.ConnectTimeOut = time.Second * time.Duration(timeout)
	token, err := web.AppConfig.String("service.token")
	if err != nil {
		return err
	}
	a.Config["token"] = token
	return nil
}

func (a *ServiceApi) GetBaseApi() *api.BaseApi {
	return &a.BaseApi
}

func NewServiceApi() (*ServiceApi, error) {
	a := &ServiceApi{
		api.BaseApi{
			Apis:   apis,
			Config: make(map[string]string),
		},
	}
	if err := a.ServiceConfig(); err != nil {
		return nil, err
	}
	return a, nil
}

func init() {
	a, err := NewServiceApi()
	if err != nil {
		fmt.Printf("init service api error:%s", err.Error())
		os.Exit(1)
	}
	api.ApiMap[api.SERVICE_API] = a
}
