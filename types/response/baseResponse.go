package response

import "encoding/json"

type IResponse interface {
	GetResponse(response IResponse) string
	GetRequestId() string
	GetMessage() string
	GetCode() int
}

type BaseResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"msg"`
	Data      interface{} `json:"result"`
	RequestId string      `json:"requestId"`
}

func (r *BaseResponse) GetRequestId() string {
	return r.RequestId
}

func (r *BaseResponse) GetResponse(response IResponse) string {
	content, _ := json.Marshal(response)
	return string(content)
}

func (r *BaseResponse) GetMessage() string {
	return r.Message
}

func (r *BaseResponse) GetCode() int {
	return r.Code
}
