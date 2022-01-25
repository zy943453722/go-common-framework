package routers

import (
	"fmt"

	"github.com/beego/beego/v2/server/web/context"
	"go-common-framework/errors"
	"go-common-framework/services/log"
	"go-common-framework/types/response"
	"go-common-framework/util"
)

func LoginTestFilter(ctx *context.Context) {
	requestId := ctx.Input.Header("X-Request-Id")
	if requestId == "" {
		requestId = util.GenUuid()
	}
	defer authFailedView(requestId, ctx)
	ctx.Input.SetData("requestId", requestId)
	//@TODO
}

func LoginFilter(ctx *context.Context) {
	requestId := ctx.Input.Header("X-Request-Id")
	if requestId == "" {
		requestId = util.GenUuid()
	}
	defer authFailedView(requestId, ctx)
	ctx.Input.SetData("requestId", requestId)
	//@TODO
}

func authFailedView(requestId string, ctx *context.Context) {
	if e := recover(); e != nil {
		var (
			resp     *response.BaseResponse
			httpCode int
			result   string
		)
		switch t := e.(type) {
		case int:
			resp = baseApiView(t, errors.GetCodeMsg(t), requestId)
		case errors.Error:
			resp = baseApiView(t.GetErrorCode(), t.Error(), requestId)
		default:
			resp = baseApiView(errors.API_UNKOWN_ERROR, fmt.Sprintf("未知错误:%v", t), requestId)
		}
		httpCode = errors.GetHttpCode(resp.GetCode())
		result = resp.GetResponse(resp)
		ctx.Output.SetStatus(httpCode)
		ctx.Output.JSON(resp, true, false)
		log.Info(requestId, "response:%s", result)
	}
}

func baseApiView(code int, msg, requestId string) *response.BaseResponse {
	return &response.BaseResponse{
		Code:      code,
		Message:   msg,
		RequestId: requestId,
	}
}
