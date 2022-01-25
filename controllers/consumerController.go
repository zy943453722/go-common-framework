package controllers

import (
	"encoding/json"

	"go-common-framework/errors"
	"go-common-framework/models"
	"go-common-framework/services/log"
	"go-common-framework/types/request"
	"go-common-framework/types/response"
)

type ConsumerController struct {
	BaseController
}

//GetConsumerList 获取客户列表 GET /api/consumer/list
func (c *ConsumerController) GetConsumerList() {
	requestId := c.GetRequestId()
	defer c.RecoverException(requestId)
	req := &request.GetConsumerListRequest{
		Name:    c.GetString("name"),
		Company: c.GetString("company"),
		Group:   c.GetString("group"),
	}
	pageNumber, err := c.GetInt("pageNumber", models.DEFAULT_PAGENUMBER)
	if err != nil {
		panic(errors.NewCommonError(errors.API_ARGUMENTS_ERROR, "参数解析失败", err))
	}
	pageSize, err := c.GetInt("pageSize", models.DEFAULT_PAGESIZE)
	if err != nil {
		panic(errors.NewCommonError(errors.API_ARGUMENTS_ERROR, "参数解析失败", err))
	}
	all, err := c.GetInt("all", 0)
	if err != nil {
		panic(errors.NewCommonError(errors.API_ARGUMENTS_ERROR, "参数解析失败", err))
	}
	req.PageNumber = pageNumber
	req.PageSize = pageSize
	req.All = all

	if err = req.Validate(req); err != nil {
		log.Info(requestId, "validate error: %s", err.Error())
		panic(errors.NewCommonError(errors.API_ARGUMENTS_ERROR, "参数验证失败", err))
	}

	c.resp = &response.GetConsumerListResponse{
		BaseResponse: *c.BaseApiView(errors.API_OK, errors.GetCodeMsg(errors.API_OK), requestId),
		Data:         models.GetConsumerList(requestId, req),
	}
}

//ExportConsumerList 导出客户列表 GET /api/consumer/export
func (c *ConsumerController) ExportConsumerList() {
	requestId := c.GetRequestId()
	defer c.RecoverException(requestId)
	req := &request.GetConsumerListRequest{
		Name:    c.GetString("name"),
		Company: c.GetString("company"),
		Group:   c.GetString("group"),
	}
	pageNumber, err := c.GetInt("pageNumber", models.DEFAULT_PAGENUMBER)
	if err != nil {
		panic(errors.NewCommonError(errors.API_ARGUMENTS_ERROR, "参数解析失败", err))
	}
	pageSize, err := c.GetInt("pageSize", models.DEFAULT_PAGESIZE)
	if err != nil {
		panic(errors.NewCommonError(errors.API_ARGUMENTS_ERROR, "参数解析失败", err))
	}
	all, err := c.GetInt("all", 0)
	if err != nil {
		panic(errors.NewCommonError(errors.API_ARGUMENTS_ERROR, "参数解析失败", err))
	}
	req.PageNumber = pageNumber
	req.PageSize = pageSize
	req.All = all

	if err = req.Validate(req); err != nil {
		log.Info(requestId, "validate error: %s", err.Error())
		panic(errors.NewCommonError(errors.API_ARGUMENTS_ERROR, "参数验证失败", err))
	}

	fileName, downloadFileName := models.ExportConsumerList(requestId, req)
	c.Ctx.Input.SetData(DOWNLOAD, "yes")
	c.Ctx.Output.Download(fileName, downloadFileName)
}

//GetConsumerInfo 获取客户详情 GET /api/consumer/info
func (c *ConsumerController) GetConsumerInfo() {
	requestId := c.GetRequestId()
	defer c.RecoverException(requestId)

	req := &request.GetConsumerInfoRequest{
		Cid: c.GetString("cid"),
	}

	if err := req.Validate(req); err != nil {
		log.Info(requestId, "validate error: %s", err.Error())
		panic(errors.NewCommonError(errors.API_ARGUMENTS_ERROR, "参数验证失败", err))
	}

	c.resp = &response.GetConsumerInfoResponse{
		BaseResponse: *c.BaseApiView(errors.API_OK, errors.GetCodeMsg(errors.API_OK), requestId),
		Data:         models.GetConsumerInfo(requestId, req),
	}
}

//EditConsumer 编辑客户 PUT /api/consumer/edit
func (c *ConsumerController) EditConsumer() {
	requestId := c.GetRequestId()
	defer c.RecoverException(requestId)

	req := &request.EditConsumerRequest{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, req); err != nil {
		log.Info(requestId, "json unmarshal error: %s", err.Error())
		panic(errors.NewCommonError(errors.API_ARGUMENTS_ERROR, "参数解析失败", err))
	}

	if err := req.Validate(req); err != nil {
		log.Info(requestId, "validate error: %s", err.Error())
		panic(errors.NewCommonError(errors.API_ARGUMENTS_ERROR, "参数验证失败", err))
	}

	c.resp = &response.CommonSuccessResponse{
		BaseResponse: *c.BaseApiView(errors.API_OK, errors.GetCodeMsg(errors.API_OK), requestId),
		Data:         models.EditConsumer(requestId, c.GetErp(), req),
	}
}
