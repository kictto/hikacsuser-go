package acsapi

import (
	"fmt"
)

// AddFaceByBinaryWithInfo 使用FaceInfo结构体通过二进制方式添加人脸
func (c *ACSClient) AddFaceByBinaryWithInfo(faceInfo FaceInfo) (ResponseData, error) {
	var response ResponseData
	if c.lUserID < 0 {
		return response, fmt.Errorf("未登录设备")
	}

	// 设置默认值
	if faceInfo.FaceLibType == "" {
		faceInfo.FaceLibType = "blackFD"
	}
	if faceInfo.FDID == "" {
		faceInfo.FDID = "1"
	}

	// 检查是否提供了人脸图片文件路径或二进制数据
	if faceInfo.FaceFile == "" && len(faceInfo.FaceData) == 0 {
		return response, fmt.Errorf("未提供人脸图片文件路径或二进制数据")
	}

	// 构建人脸信息JSON
	jsonData := fmt.Sprintf(`{
		"faceLibType": "%s",
		"FDID": "%s",
		"FPID": "%s"
	}`, faceInfo.FaceLibType, faceInfo.FDID, faceInfo.EmployeeNo)

	// 如果提供了二进制数据，则直接使用
	var respData []byte
	var err error
	if len(faceInfo.FaceData) > 0 {
		respData, err = c.faceManage.AddFaceByBinaryWithJSONAndData(c.lUserID, jsonData, faceInfo.FaceData)
		if err != nil {
			return response, err
		}
	} else {
		// 否则使用文件路径
		respData, err = c.faceManage.AddFaceByBinaryWithJSON(c.lUserID, jsonData, faceInfo.FaceFile)
		if err != nil {
			return response, err
		}
	}

	// 解析响应数据
	response, err = ParseResponseData(respData)
	return response, err
}

// AddFacesByBinaryWithInfo 批量添加人脸，使用FaceInfo结构体数组
func (c *ACSClient) AddFacesByBinaryWithInfo(faceInfos []FaceInfo) ([]ResponseData, []error) {
	if c.lUserID < 0 {
		return nil, []error{fmt.Errorf("未登录设备")}
	}

	responses := make([]ResponseData, 0)
	errors := make([]error, 0)
	for _, faceInfo := range faceInfos {
		resp, err := c.AddFaceByBinaryWithInfo(faceInfo)
		responses = append(responses, resp)
		if err != nil {
			errors = append(errors, fmt.Errorf("添加人脸 %s 失败: %v", faceInfo.EmployeeNo, err))
		}
	}

	return responses, errors
}

// AddFaceByUrlWithInfo 使用FaceInfo结构体通过URL方式添加人脸
func (c *ACSClient) AddFaceByUrlWithInfo(faceInfo FaceInfo) (ResponseData, error) {
	var response ResponseData
	if c.lUserID < 0 {
		return response, fmt.Errorf("未登录设备")
	}

	// 设置默认值
	if faceInfo.FaceLibType == "" {
		faceInfo.FaceLibType = "blackFD"
	}
	if faceInfo.FDID == "" {
		faceInfo.FDID = "1"
	}

	// 检查是否提供了人脸图片URL
	if faceInfo.FaceURL == "" {
		return response, fmt.Errorf("未提供人脸图片URL")
	}

	// 构建人脸信息JSON
	jsonData := fmt.Sprintf(`{
		"faceLibType": "%s",
		"FDID": "%s",
		"FPID": "%s",
		"faceURL": "%s"
	}`, faceInfo.FaceLibType, faceInfo.FDID, faceInfo.EmployeeNo, faceInfo.FaceURL)

	respData, err := c.faceManage.AddFaceByUrlWithJSON(c.lUserID, jsonData)
	if err != nil {
		return response, err
	}

	// 解析响应数据
	response, err = ParseResponseData(respData)
	return response, err
}

// AddFacesByUrlWithInfo 批量添加人脸，使用FaceInfo结构体数组
func (c *ACSClient) AddFacesByUrlWithInfo(faceInfos []FaceInfo) ([]ResponseData, []error) {
	if c.lUserID < 0 {
		return nil, []error{fmt.Errorf("未登录设备")}
	}

	responses := make([]ResponseData, 0)
	errors := make([]error, 0)
	for _, faceInfo := range faceInfos {
		resp, err := c.AddFaceByUrlWithInfo(faceInfo)
		responses = append(responses, resp)
		if err != nil {
			errors = append(errors, fmt.Errorf("添加人脸 %s 失败: %v", faceInfo.EmployeeNo, err))
		}
	}

	return responses, errors
}
