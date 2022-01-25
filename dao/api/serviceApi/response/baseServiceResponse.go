package response

type BaseServiceResponse struct {
	Result    interface{}         `json:"result"`
	RequestId string              `json:"requestId"`
	Err       *ServiceResponseErr `json:"error"`
}

type ServiceResponseErr struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
