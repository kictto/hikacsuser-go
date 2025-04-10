package acsapi

import (
	"fmt"
)

// AddFaceByBinaryWithInfo 使用FaceInfo结构体通过二进制方式添加人脸
func (c *ACSClient) AddFaceByBinaryWithInfo(faceInfo FaceInfo) error {
	if c.lUserID < 0 {
		return fmt.Errorf("未登录设备")
	}

	// 设置默认值
	if faceInfo.FaceLibType == "" {
		faceInfo.FaceLibType = "blackFD"
	}
	if faceInfo.FDID == "" {
		faceInfo.FDID = "1"
	}

	// 检查是否提供了人脸图片文件路径
	if faceInfo.FaceFile == "" {
		return fmt.Errorf("未提供人脸图片文件路径")
	}

	// 构建人脸信息JSON
	jsonData := fmt.Sprintf(`{
		"faceLibType": "%s",
		"FDID": "%s",
		"FPID": "%s"
	}`, faceInfo.FaceLibType, faceInfo.FDID, faceInfo.EmployeeNo)

	return c.faceManage.AddFaceByBinaryWithJSON(c.lUserID, jsonData, faceInfo.FaceFile)
}

// AddFacesByBinaryWithInfo 批量添加人脸，使用FaceInfo结构体数组
func (c *ACSClient) AddFacesByBinaryWithInfo(faceInfos []FaceInfo) []error {
	if c.lUserID < 0 {
		return []error{fmt.Errorf("未登录设备")}
	}

	errors := make([]error, 0)
	for _, faceInfo := range faceInfos {
		err := c.AddFaceByBinaryWithInfo(faceInfo)
		if err != nil {
			errors = append(errors, fmt.Errorf("添加人脸 %s 失败: %v", faceInfo.EmployeeNo, err))
		}
	}

	return errors
}

// AddFaceByUrlWithInfo 使用FaceInfo结构体通过URL方式添加人脸
func (c *ACSClient) AddFaceByUrlWithInfo(faceInfo FaceInfo) error {
	if c.lUserID < 0 {
		return fmt.Errorf("未登录设备")
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
		return fmt.Errorf("未提供人脸图片URL")
	}

	// 构建人脸信息JSON
	jsonData := fmt.Sprintf(`{
		"faceLibType": "%s",
		"FDID": "%s",
		"FPID": "%s",
		"faceURL": "%s"
	}`, faceInfo.FaceLibType, faceInfo.FDID, faceInfo.EmployeeNo, faceInfo.FaceURL)

	return c.faceManage.AddFaceByUrlWithJSON(c.lUserID, jsonData)
}

// AddFacesByUrlWithInfo 批量添加人脸，使用FaceInfo结构体数组
func (c *ACSClient) AddFacesByUrlWithInfo(faceInfos []FaceInfo) []error {
	if c.lUserID < 0 {
		return []error{fmt.Errorf("未登录设备")}
	}

	errors := make([]error, 0)
	for _, faceInfo := range faceInfos {
		err := c.AddFaceByUrlWithInfo(faceInfo)
		if err != nil {
			errors = append(errors, fmt.Errorf("添加人脸 %s 失败: %v", faceInfo.EmployeeNo, err))
		}
	}

	return errors
}
