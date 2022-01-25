package response

type GetConsumerListResponse struct {
	BaseResponse
	Data *ConsumerList `json:"result"`
}

type ConsumerList struct {
	Data       []*ConsumerInfo `json:"data"`
	TotalCount int             `json:"totalCount"`
	PageNumber int             `json:"pageNumber"`
	PageSize   int             `json:"pageSize"`
}

type ConsumerInfo struct {
	ConsumerId    string `json:"consumerId"`
	Source        string `json:"source"`
	SourceName    string `json:"sourceName"`
	ConsumerName  string `json:"consumerName"`
	CompanyName   string `json:"companyName"`
	UserGroup     string `json:"userGroup"`
	UserGroupName string `json:"userGroupName"`
	UpdateTime    string `json:"updateTime"`
	UpdateErp     string `json:"updateErp"`
}
