package apiTransfer

import (
	"encoding/json"
	"go-common-framework/dao/api"
	serviceRequest "go-common-framework/dao/api/serviceApi/request"
	serviceResponse "go-common-framework/dao/api/serviceApi/response"
	"go-common-framework/util"
	"strconv"
	"time"
)

func GetDeviceList(requestId string, req *serviceRequest.GetDeviceRequest) (*serviceResponse.GetDeviceListResponse, error) {
	api.Locker.Lock()
	a, err := api.GetApi(api.SERVICE_API)
	if err != nil {
		return nil, err
	}
	base := a.GetBaseApi()
	t := time.Now().Unix()
	auth := util.GetSha512(base.Config["token"] + strconv.Itoa(int(t)))
	base.Options = &api.Options{
		Headers: map[string]string{
			"x-timestamp":   strconv.Itoa(int(t)),
			"Authorization": auth,
		},
		Data: map[string]api.IBaseApiRequest{
			api.PARAMS_TYPE_QUERY: req,
		},
	}
	resp, err := api.Send(requestId, a, "getDeviceList")
	if err != nil {
		return nil, err
	}
	getDeviceListResp := new(serviceResponse.GetDeviceListResponse)
	if err = json.Unmarshal(resp, getDeviceListResp); err != nil {
		return nil, err
	}
	return getDeviceListResp, nil
}
