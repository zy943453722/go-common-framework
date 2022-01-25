package response

type GetDeviceListResponse struct {
	BaseServiceResponse
	Result *GetDeviceListResult `json:"result"`
}

type GetDeviceListResult struct {
	Data       []*GetDeviceListInfo `json:"data"`
	PageNumber int                  `json:"pageNumber"`
	PageSize   int                  `json:"pageSize"`
	TotalCount int                  `json:"totalCount"`
}

type GetDeviceListInfo struct {
	Month        string
	ResourceId   string
	ResourceName string
	ConsumerId   int
	ConsumerName string
	IdcId        int
}
