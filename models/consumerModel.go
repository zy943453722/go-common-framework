package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
	"go-common-framework/dao/table"
	"go-common-framework/errors"
	"go-common-framework/services/excel"
	"go-common-framework/types/request"
	"go-common-framework/types/response"
)

func GetConsumerList(requestId string, req *request.GetConsumerListRequest) *response.ConsumerList {
	params := map[string]interface{}{
		"name":      req.Name,
		"company":   req.Company,
		"userGroup": req.Group,
	}
	offset := (req.PageNumber - 1) * req.PageSize
	limit := req.PageSize
	consumerList, count, err := table.GetConsumerList(requestId, params, limit, offset, req.All)
	if err != nil {
		panic(errors.NewCommonError(errors.API_DB_ERROR, "查询客户列表失败", err))
	}
	resp := &response.ConsumerList{
		TotalCount: count,
		PageNumber: req.PageNumber,
		PageSize:   req.PageSize,
	}
	if count > 0 {
		consumerResp := make([]*response.ConsumerInfo, 0)
		for _, consumerInfo := range consumerList {
			tmp := &response.ConsumerInfo{
				ConsumerId:    consumerInfo.Cid,
				ConsumerName:  consumerInfo.Name,
				CompanyName:   consumerInfo.Company,
				UserGroup:     consumerInfo.UserGroup,
				UserGroupName: table.UserGroupMap[consumerInfo.UserGroup],
				UpdateTime:    time.Unix(consumerInfo.UpdatedAt, 0).Format(table.TIME),
				UpdateErp:     consumerInfo.UpdatedErp,
			}
			consumerResp = append(consumerResp, tmp)
		}
		resp.Data = consumerResp
	}
	return resp
}

func ExportConsumerList(requestId string, req *request.GetConsumerListRequest) (string, string) {
	resp := GetConsumerList(requestId, req)
	e := excel.NewExcelInstance()
	index := e.SetNewSheet("Sheet1")
	sheetData := make([][]string, 0)
	for _, data := range resp.Data {
		sheetData = append(sheetData, []string{
			data.ConsumerId,
			data.ConsumerName,
			data.CompanyName,
			data.UserGroupName,
			data.UpdateTime,
			data.UpdateErp,
		})
	}
	if err := e.SetSheetData("Sheet1", ConsumerHeader, sheetData); err != nil {
		panic(errors.NewCommonError(errors.API_INTERNAL_ERROR, "新建excel失败", err))
	}
	e.File.SetActiveSheet(index)
	downloadFileName := CONSUMER_FILE_PREFIX + time.Now().Format(table.DATE) + ".xlsx"
	fileName, err := e.SaveFile(downloadFileName)
	if err != nil {
		panic(errors.NewCommonError(errors.API_INTERNAL_ERROR, "保存excel到本地失败", err))
	}
	return fileName, downloadFileName
}

func GetConsumerInfo(requestId string, req *request.GetConsumerInfoRequest) *response.ConsumerInfo {
	consumerInfo, err := table.GetConsumerInfo(requestId, req.Cid)
	if err != nil {
		panic(errors.NewCommonError(errors.API_DB_ERROR, "查询客户详情失败", err))
	}
	if consumerInfo != nil {
		return &response.ConsumerInfo{
			ConsumerId:    consumerInfo.Cid,
			ConsumerName:  consumerInfo.Name,
			CompanyName:   consumerInfo.Company,
			UserGroup:     consumerInfo.UserGroup,
			UserGroupName: table.UserGroupMap[consumerInfo.UserGroup],
			UpdateTime:    time.Unix(consumerInfo.UpdatedAt, 0).Format(table.TIME),
			UpdateErp:     consumerInfo.UpdatedErp,
		}
	} else {
		return nil
	}
}

func EditConsumer(requestId, erp string, req *request.EditConsumerRequest) *response.SuccessResponse {
	consumerInfo, err := table.GetConsumerInfo(requestId, req.Cid)
	if err != nil {
		panic(errors.NewCommonError(errors.API_DB_ERROR, "查询客户详情失败", err))
	}
	params := make(orm.Params, 0)
	if req.Company != consumerInfo.Company && req.Company != "" {
		params["company"] = req.Company
	}
	if req.UserGroup != consumerInfo.UserGroup && req.UserGroup != "" {
		params["user_group"] = req.UserGroup
	}
	if len(params) > 0 {
		params["updated_erp"] = erp
		params["updated_at"] = time.Now().Unix()
		if err = table.UpdateConsumer(requestId, req.Cid, params); err != nil {
			panic(errors.NewCommonError(errors.API_DB_ERROR, "编辑客户失败", err))
		}
	}
	return &response.SuccessResponse{
		Success: true,
	}
}
