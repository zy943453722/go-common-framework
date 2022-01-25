package response

type GetConsumerInfoResponse struct {
	BaseResponse
	Data *ConsumerInfo `json:"result"`
}
