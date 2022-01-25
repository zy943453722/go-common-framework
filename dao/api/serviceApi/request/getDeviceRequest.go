package request

import (
	"bytes"
	"encoding/json"
)

type GetDeviceRequest struct {
	PageNumber int    `json:"pageNumber,omitempty"`
	PageSize   int    `json:"pageSize,omitempty"`
	Month      string `json:"month"`
	All        int    `json:"all"`
}

func (r *GetDeviceRequest) Decode() (map[string]interface{}, error) {
	var err error
	data, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	mapData := make(map[string]interface{})
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err = decoder.Decode(&mapData); err != nil {
		return nil, err
	}

	return mapData, nil
}
