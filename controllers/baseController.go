package controllers

import (
	"fmt"
	"go-common-framework/util"
	"net/url"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/server/web"
	"go-common-framework/errors"
	"go-common-framework/services/log"
	"go-common-framework/types/response"
)

const (
	DOWNLOAD = "download"
)

type BaseController struct {
	resp response.IResponse
	web.Controller
}

func (c *BaseController) Prepare() {
	web.ReadFromRequest(&c.Controller)
	requestId, ok := c.GetHeader()["requestId"]
	if !ok && requestId == "" {
		if c.GetRequestId() == "" {
			requestId = util.GenUuid()
		} else {
			requestId = c.GetRequestId()
		}
	}
	c.Ctx.Input.SetData("requestId", requestId)
	log.Info(requestId, genRequestLog(c))
}

func (c *BaseController) Finish() {
	res := c.Ctx.Input.GetData(DOWNLOAD)
	if res == nil {
		var (
			httpCode  int
			requestId string
			result    string
		)
		if c.resp != nil {
			httpCode = errors.GetHttpCode(c.resp.GetCode())
			requestId = c.resp.GetRequestId()
			result = c.resp.GetResponse(c.resp)
		} else {
			httpCode = 200
		}
		c.Ctx.Output.SetStatus(httpCode)
		c.Data["json"] = c.resp
		c.ServeJSON()
		log.Info(requestId, "response:%s", result)
	}
}

func (c *BaseController) BaseApiView(code int, msg, requestId string) *response.BaseResponse {
	return &response.BaseResponse{
		Code:      code,
		Message:   msg,
		RequestId: requestId,
	}
}

func (c *BaseController) RecoverException(requestId string) {
	if e := recover(); e != nil {
		switch t := e.(type) {
		case int:
			c.resp = c.BaseApiView(t, errors.GetCodeMsg(t), requestId)
		case errors.Error:
			c.resp = c.BaseApiView(t.GetErrorCode(), t.Error(), requestId)
		default:
			c.resp = c.BaseApiView(errors.API_UNKOWN_ERROR, fmt.Sprintf("未知错误:%v", t), requestId)
		}
	}
}

func (c *BaseController) GetRequestId() string {
	return c.Ctx.Input.GetData("requestId").(string)
}

func (c *BaseController) GetErp() string {
	return c.Ctx.Input.GetData("loginErp").(string)
}

func (c *BaseController) GetHeader() map[string]string {
	return map[string]string{
		"requestId": c.Ctx.Input.Header("X-Request-Id"),
	}
}

func genRequestLog(c *BaseController) string {
	cmd := make([]string, 0)
	cmd = append(cmd, "curl", "-X", c.Ctx.Request.Method)
	for key, values := range c.Ctx.Request.Header {
		for _, value := range values {
			cmd = append(cmd, "-H", fmt.Sprintf("'%s:%s'", key, value))
		}
	}
	if c.Ctx.Request.Host != "" {
		cmd = append(cmd, "-H", fmt.Sprintf("'%s:%s'", "Host", c.Ctx.Request.Host))
	}

	if contentType := c.Ctx.Input.Header("Content-Type"); contentType != "" {
		body := strings.Replace(string(c.Ctx.Input.RequestBody), "\n", "", -1)
		if len(contentType) > 0 && body != "" {
			cmd = append(cmd, "-d", "'"+body+"'")
		}
	}

	uri, _ := url.QueryUnescape(c.Ctx.Input.URI())
	reqUrl := c.Ctx.Input.Site() + ":" + strconv.Itoa(c.Ctx.Input.Port()) + uri
	cmd = append(cmd, reqUrl)
	curl := strings.Join(cmd, " ")
	return curl
}
