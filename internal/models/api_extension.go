// Package models 提供门禁系统的内部模型实现
package models

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"fmt"
	"github.com/clockworkchen/hikacsuser-go/internal/sdk"
	"github.com/clockworkchen/hikacsuser-go/internal/utils"
	"time"
	"unsafe"
)

// SendISAPIRequest 公开的ISAPI请求发送方法，供外部包调用
func (um *UserManage) SendISAPIRequest(lUserID int, url string, requestData []byte) ([]byte, error) {
	return um.sendISAPIRequest(lUserID, url, requestData)
}

// SendISAPIRequest 公开的ISAPI请求发送方法，供外部包调用
func (cm *CardManage) SendISAPIRequest(lUserID int, url string, requestData []byte) ([]byte, error) {
	return cm.sendISAPIRequest(lUserID, url, requestData)
}

// SendISAPIRequest 公开的ISAPI请求发送方法，供外部包调用
func (fm *FaceManage) SendISAPIRequest(lUserID int, url string, requestData []byte) ([]byte, error) {
	return fm.sendISAPIRequest(lUserID, url, requestData)
}

// SendISAPIRequest 公开的ISAPI请求发送方法，供外部包调用
func (am *ACSManage) SendISAPIRequest(lUserID int, url string, requestData []byte) ([]byte, error) {
	return am.sendISAPIRequest(lUserID, url, requestData)
}

// SendISAPIRequest 公开的ISAPI请求发送方法，供外部包调用
func (es *EventSearch) SendISAPIRequest(lUserID int, url string, requestData []byte) ([]byte, error) {
	return es.sendISAPIRequest(lUserID, url, requestData)
}

// AddUserInfoWithJSON 使用完整的JSON数据添加用户信息
func (um *UserManage) AddUserInfoWithJSON(lUserID int, jsonData string) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID")
	}

	// URL
	url := "POST /ISAPI/AccessControl/UserInfo/Record?format=json"

	// 发送ISAPI请求
	response, err := um.sendISAPIRequest(lUserID, url, []byte(jsonData))
	if err != nil {
		return err
	}

	fmt.Printf("添加用户成功, 响应: %s\n", string(response))
	return nil
}

// DeleteUserInfoWithJSON 使用完整的JSON数据删除用户信息
func (um *UserManage) DeleteUserInfoWithJSON(lUserID int, jsonData string) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID")
	}

	// URL
	url := "PUT /ISAPI/AccessControl/UserInfo/Delete?format=json"

	// 发送ISAPI请求
	response, err := um.sendISAPIRequest(lUserID, url, []byte(jsonData))
	if err != nil {
		return err
	}

	fmt.Printf("删除用户信息成功, 响应: %s\n", string(response))
	return nil
}

// AddCardInfoWithJSON 使用完整的JSON数据添加卡片信息
func (cm *CardManage) AddCardInfoWithJSON(lUserID int, jsonData string) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID")
	}

	// URL
	url := "POST /ISAPI/AccessControl/CardInfo/Record?format=json"

	// 发送ISAPI请求
	response, err := cm.sendISAPIRequest(lUserID, url, []byte(jsonData))
	if err != nil {
		return err
	}

	fmt.Printf("添加卡片成功, 响应: %s\n", string(response))
	return nil
}

// AddFaceByBinaryWithJSON 使用完整的JSON数据通过二进制方式添加人脸
func (fm *FaceManage) AddFaceByBinaryWithJSON(lUserID int, jsonData string, faceFilePath string) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID")
	}

	// 加载人脸图片
	faceData, err := utils.LoadPicture(faceFilePath)
	if err != nil {
		return fmt.Errorf("加载人脸图片失败: %v", err)
	}

	// 构建请求URL
	strInBuffer := "PUT /ISAPI/Intelligent/FDLib/FDSetUp?format=json"

	// 创建ISAPI输入参数
	inputParamPtr := C.malloc(C.size_t(len(strInBuffer) + 1))
	defer C.free(inputParamPtr)
	C.memcpy(inputParamPtr, unsafe.Pointer(&[]byte(strInBuffer)[0]), C.size_t(len(strInBuffer)))

	// 调用远程配置
	lHandle := fm.SDK.NET_DVR_StartRemoteConfig(lUserID, NET_DVR_FACE_DATA_RECORD, inputParamPtr, uint32(len(strInBuffer)), 0, nil)
	if lHandle < 0 {
		return fmt.Errorf("NET_DVR_StartRemoteConfig失败，错误码: %d", fm.SDK.NET_DVR_GetLastError())
	}
	defer fm.SDK.NET_DVR_StopRemoteConfig(lHandle)

	// 定义类似于NET_DVR_JSON_DATA_CFG的结构体
	type NET_DVR_JSON_DATA_CFG struct {
		dwSize                  uint32
		lpJsonData              unsafe.Pointer // JSON报文
		dwJsonDataSize          uint32         // JSON报文大小
		lpPicData               unsafe.Pointer // 图片内容
		dwPicDataSize           uint32         // 图片内容大小
		lpInfraredFacePicBuffer int32          // 红外人脸图片数据缓存
		dwInfraredFacePicSize   unsafe.Pointer // 红外人脸图片数据大小
		byRes                   [248]byte      // 保留
	}

	// 创建并填充NET_DVR_JSON_DATA_CFG结构体
	var struAddFaceDataCfg NET_DVR_JSON_DATA_CFG

	// 使用C.malloc分配内存并复制JSON数据
	jsonDataPtr := C.malloc(C.size_t(len(jsonData)))
	defer C.free(jsonDataPtr)
	C.memcpy(jsonDataPtr, unsafe.Pointer(&[]byte(jsonData)[0]), C.size_t(len(jsonData)))
	struAddFaceDataCfg.lpJsonData = jsonDataPtr
	struAddFaceDataCfg.dwJsonDataSize = uint32(len(jsonData))

	// 使用C.malloc分配内存并复制图片数据
	picDataPtr := C.malloc(C.size_t(len(faceData)))
	defer C.free(picDataPtr)
	C.memcpy(picDataPtr, unsafe.Pointer(&faceData[0]), C.size_t(len(faceData)))
	struAddFaceDataCfg.lpPicData = picDataPtr
	struAddFaceDataCfg.dwPicDataSize = uint32(len(faceData))

	// 设置结构体大小
	struAddFaceDataCfg.dwSize = uint32(unsafe.Sizeof(struAddFaceDataCfg))

	// 准备输出缓冲区
	var outputBuf [1024]byte
	var resultLen uint32

	// 发送数据及接收结果
	result := fm.SDK.NET_DVR_SendWithRecvRemoteConfig(lHandle, unsafe.Pointer(&struAddFaceDataCfg), struAddFaceDataCfg.dwSize, unsafe.Pointer(&outputBuf[0]), uint32(len(outputBuf)), &resultLen)

	// 处理配置结果
	if result == -1 {
		return fmt.Errorf("下发人脸失败: 发送ISAPI请求失败，错误码: %d", fm.SDK.NET_DVR_GetLastError())
	} else if result == sdk.NET_SDK_CONFIG_STATUS_NEED_WAIT {
		// 配置等待，等待一段时间后再次尝试获取结果
		time.Sleep(100 * time.Millisecond)

		// 再次尝试获取结果，但不发送数据
		emptyData := []byte{}
		for i := 0; i < 20; i++ {
			result = fm.SDK.NET_DVR_SendWithRecvRemoteConfig(lHandle, unsafe.Pointer(&emptyData[0]), 0, unsafe.Pointer(&outputBuf[0]), uint32(len(outputBuf)), &resultLen)

			if result != sdk.NET_SDK_CONFIG_STATUS_NEED_WAIT {
				break
			}

			time.Sleep(100 * time.Millisecond)
		}
	}

	// 解析返回结果
	if result == sdk.NET_SDK_CONFIG_STATUS_FAILED {
		return fmt.Errorf("下发人脸失败: 配置失败")
	} else if result == sdk.NET_SDK_CONFIG_STATUS_EXCEPTION {
		return fmt.Errorf("下发人脸失败: 配置异常")
	} else if result == sdk.NET_SDK_CONFIG_STATUS_SUCCESS || result == sdk.NET_SDK_CONFIG_STATUS_FINISH {
		// 提取响应数据
		response := make([]byte, resultLen)
		copy(response, outputBuf[:resultLen])
		fmt.Printf("添加人脸信息成功, 响应: %s\n", string(response))
		return nil
	}

	return fmt.Errorf("下发人脸失败: 未知状态码 %d", result)
}

// AddFaceByBinaryWithJSONAndData 使用完整的JSON数据和二进制图片数据通过二进制方式添加人脸
func (fm *FaceManage) AddFaceByBinaryWithJSONAndData(lUserID int, jsonData string, faceData []byte) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID")
	}

	// 检查图片数据
	if len(faceData) == 0 {
		return fmt.Errorf("人脸图片数据为空")
	}

	// 构建请求URL
	strInBuffer := "PUT /ISAPI/Intelligent/FDLib/FDSetUp?format=json"

	// 创建ISAPI输入参数
	inputParamPtr := C.malloc(C.size_t(len(strInBuffer) + 1))
	defer C.free(inputParamPtr)
	C.memcpy(inputParamPtr, unsafe.Pointer(&[]byte(strInBuffer)[0]), C.size_t(len(strInBuffer)))

	// 调用远程配置
	lHandle := fm.SDK.NET_DVR_StartRemoteConfig(lUserID, NET_DVR_FACE_DATA_RECORD, inputParamPtr, uint32(len(strInBuffer)), 0, nil)
	if lHandle < 0 {
		return fmt.Errorf("NET_DVR_StartRemoteConfig失败，错误码: %d", fm.SDK.NET_DVR_GetLastError())
	}
	defer fm.SDK.NET_DVR_StopRemoteConfig(lHandle)

	// 定义类似于NET_DVR_JSON_DATA_CFG的结构体
	type NET_DVR_JSON_DATA_CFG struct {
		dwSize                  uint32
		lpJsonData              unsafe.Pointer // JSON报文
		dwJsonDataSize          uint32         // JSON报文大小
		lpPicData               unsafe.Pointer // 图片内容
		dwPicDataSize           uint32         // 图片内容大小
		lpInfraredFacePicBuffer int32          // 红外人脸图片数据缓存
		dwInfraredFacePicSize   unsafe.Pointer // 红外人脸图片数据大小
		byRes                   [248]byte      // 保留
	}

	// 创建并填充NET_DVR_JSON_DATA_CFG结构体
	var struAddFaceDataCfg NET_DVR_JSON_DATA_CFG

	// 使用C.malloc分配内存并复制JSON数据
	jsonDataPtr := C.malloc(C.size_t(len(jsonData)))
	defer C.free(jsonDataPtr)
	C.memcpy(jsonDataPtr, unsafe.Pointer(&[]byte(jsonData)[0]), C.size_t(len(jsonData)))
	struAddFaceDataCfg.lpJsonData = jsonDataPtr
	struAddFaceDataCfg.dwJsonDataSize = uint32(len(jsonData))

	// 使用C.malloc分配内存并复制图片数据
	picDataPtr := C.malloc(C.size_t(len(faceData)))
	defer C.free(picDataPtr)
	C.memcpy(picDataPtr, unsafe.Pointer(&faceData[0]), C.size_t(len(faceData)))
	struAddFaceDataCfg.lpPicData = picDataPtr
	struAddFaceDataCfg.dwPicDataSize = uint32(len(faceData))

	// 设置结构体大小
	struAddFaceDataCfg.dwSize = uint32(unsafe.Sizeof(struAddFaceDataCfg))

	// 准备输出缓冲区
	var outputBuf [1024]byte
	var resultLen uint32

	// 发送数据及接收结果
	result := fm.SDK.NET_DVR_SendWithRecvRemoteConfig(lHandle, unsafe.Pointer(&struAddFaceDataCfg), struAddFaceDataCfg.dwSize, unsafe.Pointer(&outputBuf[0]), uint32(len(outputBuf)), &resultLen)

	// 处理配置结果
	if result == -1 {
		return fmt.Errorf("下发人脸失败: 发送ISAPI请求失败，错误码: %d", fm.SDK.NET_DVR_GetLastError())
	} else if result == sdk.NET_SDK_CONFIG_STATUS_NEED_WAIT {
		// 配置等待，等待一段时间后再次尝试获取结果
		time.Sleep(100 * time.Millisecond)

		// 再次尝试获取结果，但不发送数据
		emptyData := []byte{}
		for i := 0; i < 20; i++ {
			result = fm.SDK.NET_DVR_SendWithRecvRemoteConfig(lHandle, unsafe.Pointer(&emptyData[0]), 0, unsafe.Pointer(&outputBuf[0]), uint32(len(outputBuf)), &resultLen)

			if result != sdk.NET_SDK_CONFIG_STATUS_NEED_WAIT {
				break
			}

			time.Sleep(100 * time.Millisecond)
		}
	}

	// 解析返回结果
	if result == sdk.NET_SDK_CONFIG_STATUS_FAILED {
		return fmt.Errorf("下发人脸失败: 配置失败")
	} else if result == sdk.NET_SDK_CONFIG_STATUS_EXCEPTION {
		return fmt.Errorf("下发人脸失败: 配置异常")
	} else if result == sdk.NET_SDK_CONFIG_STATUS_SUCCESS || result == sdk.NET_SDK_CONFIG_STATUS_FINISH {
		// 提取响应数据
		response := make([]byte, resultLen)
		copy(response, outputBuf[:resultLen])
		fmt.Printf("添加人脸信息成功, 响应: %s\n", string(response))
		return nil
	}

	return fmt.Errorf("下发人脸失败: 未知状态码 %d", result)
}

// AddFaceByUrlWithJSON 使用完整的JSON数据通过URL方式添加人脸
func (fm *FaceManage) AddFaceByUrlWithJSON(lUserID int, jsonData string) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID")
	}

	// URL
	url := "PUT /ISAPI/Intelligent/FDLib/FDSetUp?format=json"

	// 发送ISAPI请求
	response, err := fm.sendISAPIRequest(lUserID, url, []byte(jsonData))
	if err != nil {
		return err
	}

	fmt.Printf("添加人脸信息成功, 响应: %s\n", string(response))
	return nil
}
