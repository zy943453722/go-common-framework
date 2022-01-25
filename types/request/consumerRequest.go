package request

type EditConsumerRequest struct {
	BaseRequest
	Cid       string `json:"cid" valid:"Required"`
	Company   string `json:"company"`
	UserGroup string `json:"userGroup"`
}

type GetConsumerInfoRequest struct {
	BaseRequest
	Cid string `json:"cid" valid:"Required"`
}

type GetConsumerListRequest struct {
	BaseRequest
	Name       string `json:"name"`
	Company    string `json:"company"`
	Group      string `json:"group"`
	PageNumber int    `json:"pageNumber" valid:"Min(0)"`
	PageSize   int    `json:"pageSize" valid:"Range(0,100)"`
	All        int    `json:"all" valid:"Range(0,1)"`
}
