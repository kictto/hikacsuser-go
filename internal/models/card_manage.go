package models

import (
	"fmt"
	"github.com/hikacsuser-go/internal/sdk"
	"time"
	"unsafe"
)

// CardManage 卡片管理
type CardManage struct {
	SDK sdk.HCNetSDK
}

// NewCardManage 创建卡片管理实例
func NewCardManage(sdk sdk.HCNetSDK) *CardManage {
	return &CardManage{
		SDK: sdk,
	}
}

// AddCardInfo 添加卡片信息
func (cm *CardManage) AddCardInfo(lUserID int, employeeNo, cardNo string) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID")
	}

	// 构建卡片信息JSON
	jsonData := fmt.Sprintf(`{
		"CardInfo": {
			"employeeNo": "%s",
			"cardNo": "%s",
			"cardType": "normalCard"
		}
	}`, employeeNo, cardNo)

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

// SearchCardInfo 查询卡片信息
func (cm *CardManage) SearchCardInfo(lUserID int, employeeNo string) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID")
	}

	// 生成UUID作为searchID，与Java实现保持一致
	uuid := fmt.Sprintf("%d", time.Now().UnixNano())

	// 查询条件JSON
	jsonData := fmt.Sprintf(`{
		"CardInfoSearchCond": {
			"searchID": "%s",
			"searchResultPosition": 0,
			"maxResults": 30,
			"EmployeeNoList": [
				{
					"employeeNo": "%s"
				}
			]
		}
	}`, uuid, employeeNo)

	// 打印查询JSON
	fmt.Printf("查询的json报文: %s\n", jsonData)

	// URL
	url := "POST /ISAPI/AccessControl/CardInfo/Search?format=json"

	// 发送ISAPI请求
	response, err := cm.sendISAPIRequest(lUserID, url, []byte(jsonData))
	if err != nil {
		return err
	}

	fmt.Printf("查询卡片信息成功, 响应: %s\n", string(response))
	return nil
}

// DeleteCardInfo 删除卡片信息
func (cm *CardManage) DeleteCardInfo(lUserID int, cardNo string) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID")
	}

	// 删除条件JSON
	jsonData := fmt.Sprintf(`{
		"CardInfoDelCond": {
			"CardNoList": [
				{
					"cardNo": "%s"
				}
			]
		}
	}`, cardNo)

	// URL
	url := "PUT /ISAPI/AccessControl/CardInfo/Delete?format=json"

	// 发送ISAPI请求
	response, err := cm.sendISAPIRequest(lUserID, url, []byte(jsonData))
	if err != nil {
		return err
	}

	fmt.Printf("删除卡片信息成功, 响应: %s\n", string(response))
	return nil
}

// sendISAPIRequest 发送ISAPI请求
func (cm *CardManage) sendISAPIRequest(lUserID int, url string, requestData []byte) ([]byte, error) {
	// 创建ISAPI输入参数
	inputParam := make([]byte, len(url)+1)
	copy(inputParam, []byte(url))

	var outputBuf [20 * 1024]byte // 增大缓冲区大小，与Java实现保持一致

	// 调用远程配置
	lHandle := cm.SDK.NET_DVR_StartRemoteConfig(lUserID, sdk.NET_DVR_JSON_CONFIG, unsafe.Pointer(&inputParam[0]), uint32(len(inputParam)), 0, nil)
	if lHandle < 0 {
		return nil, fmt.Errorf("NET_DVR_StartRemoteConfig失败，错误码: %d", cm.SDK.NET_DVR_GetLastError())
	}
	defer cm.SDK.NET_DVR_StopRemoteConfig(lHandle)

	// 发送数据
	var resultLen uint32
	var result int

	// 只发送一次请求，不要在循环中重复发送
	if requestData == nil || len(requestData) == 0 {
		// 如果没有请求数据，只发送空数据
		emptyData := []byte{}
		result = cm.SDK.NET_DVR_SendWithRecvRemoteConfig(lHandle, unsafe.Pointer(&emptyData[0]), 0, unsafe.Pointer(&outputBuf[0]), uint32(len(outputBuf)), &resultLen)
	} else {
		// 发送请求数据
		result = cm.SDK.NET_DVR_SendWithRecvRemoteConfig(lHandle, unsafe.Pointer(&requestData[0]), uint32(len(requestData)), unsafe.Pointer(&outputBuf[0]), uint32(len(outputBuf)), &resultLen)
	}

	// 处理发送结果
	fmt.Printf("NET_DVR_SendWithRecvRemoteConfig结果: %d, 返回字节数: %d\n", result, resultLen)

	// 处理配置结果
	if result == -1 {
		return nil, fmt.Errorf("发送ISAPI请求失败，错误码: %d", cm.SDK.NET_DVR_GetLastError())
	} else if result == sdk.NET_SDK_CONFIG_STATUS_NEED_WAIT {
		// 配置等待，等待一段时间后再次尝试获取结果
		time.Sleep(10 * time.Millisecond)

		// 再次尝试获取结果，但不发送数据
		emptyData := []byte{}
		for i := 0; i < 10; i++ { // 最多尝试10次
			result = cm.SDK.NET_DVR_SendWithRecvRemoteConfig(lHandle, unsafe.Pointer(&emptyData[0]), 0, unsafe.Pointer(&outputBuf[0]), uint32(len(outputBuf)), &resultLen)

			if result != sdk.NET_SDK_CONFIG_STATUS_NEED_WAIT {
				break
			}

			time.Sleep(100 * time.Millisecond)
		}
	}

	// 再次检查结果
	if result == sdk.NET_SDK_CONFIG_STATUS_FAILED {
		return nil, fmt.Errorf("配置失败")
	} else if result == sdk.NET_SDK_CONFIG_STATUS_EXCEPTION {
		return nil, fmt.Errorf("配置异常")
	} else if result == sdk.NET_SDK_CONFIG_STATUS_SUCCESS {
		// 提取响应数据
		response := make([]byte, resultLen)
		copy(response, outputBuf[:resultLen])
		return response, nil
	} else if result == sdk.NET_SDK_CONFIG_STATUS_FINISH {
		// 配置完成
		response := make([]byte, resultLen)
		copy(response, outputBuf[:resultLen])
		return response, nil
	}

	// 如果没有明确的成功或失败状态，返回当前获取的数据
	response := make([]byte, resultLen)
	copy(response, outputBuf[:resultLen])
	return response, nil
}
