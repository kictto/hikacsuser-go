package acsapi

import (
	"encoding/json"
)

// ResponseData 设备响应数据结构体
type ResponseData struct {
	StatusCode    int    `json:"statusCode"`    // 状态码
	StatusString  string `json:"statusString"`  // 状态描述
	SubStatusCode string `json:"subStatusCode"` // 子状态码
	ErrorCode     int    `json:"errorCode"`     // 错误码
	ErrorMsg      string `json:"errorMsg"`      // 错误信息
	RawData       []byte `json:"-"`             // 原始响应数据
}

// ParseResponseData 解析设备响应数据
func ParseResponseData(data []byte) (ResponseData, error) {
	var response ResponseData
	response.RawData = data

	// 如果响应数据为空，返回空响应
	if len(data) == 0 {
		return response, nil
	}

	// 尝试解析JSON数据
	err := json.Unmarshal(data, &response)
	if err != nil {
		// 如果解析失败，保留原始数据但不返回错误
		// 因为有些API可能返回非JSON格式的数据
	}

	return response, nil
}
