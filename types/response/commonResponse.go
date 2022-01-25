package response

type SuccessResponse struct {
	Success bool `json:"success"`
}

type CommonSuccessResponse struct {
	BaseResponse
	Data *SuccessResponse `json:"result"`
}
