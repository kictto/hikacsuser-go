package models

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/clockworkchen/hikacsuser-go/internal/sdk"
	"github.com/clockworkchen/hikacsuser-go/internal/utils"
	"os"
	"path/filepath"
	"time"
	"unsafe"
)

// FaceManage 人脸管理
type FaceManage struct {
	SDK sdk.HCNetSDK
}

// NewFaceManage 创建人脸管理实例
func NewFaceManage(sdk sdk.HCNetSDK) *FaceManage {
	return &FaceManage{
		SDK: sdk,
	}
}

// NET_DVR_FACE_DATA_RECORD 添加人脸数据到人脸库命令码
const NET_DVR_FACE_DATA_RECORD = 2551

// AddFaceByBinary 通过二进制方式添加人脸图片
func (fm *FaceManage) AddFaceByBinary(lUserID int, employeeNo string) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID")
	}

	// 尝试多个可能的路径加载图片文件
	var filePath string
	var fileExists bool

	// 尝试从resources/pic目录加载
	filePath = utils.GetResourcePath("pic/1.jpg")
	fileExists = utils.FileExists(filePath)

	// 如果resources/pic目录下不存在，尝试从bin/resources/pic目录加载
	if !fileExists {
		execPath, err := os.Executable()
		if err == nil {
			rootDir := filepath.Dir(execPath)
			filePath = filepath.Join(rootDir, "bin", "resources", "pic", "1.jpg")
			fileExists = utils.FileExists(filePath)
		}
	}

	// 如果bin/resources/pic目录下不存在，尝试直接从当前目录的pic子目录加载
	if !fileExists {
		execPath, err := os.Executable()
		if err == nil {
			rootDir := filepath.Dir(execPath)
			filePath = filepath.Join(rootDir, "pic", "1.jpg")
			fileExists = utils.FileExists(filePath)
		}
	}

	if !fileExists {
		return fmt.Errorf("人脸图片文件不存在，已尝试多个路径: %s", filePath)
	}

	fmt.Printf("使用人脸图片: %s\n", filePath)
	faceData, err := utils.LoadPicture(filePath)
	if err != nil {
		return fmt.Errorf("加载人脸图片失败: %v", err)
	}

	// 检查图片大小
	fmt.Printf("人脸图片大小: %d 字节\n", len(faceData))

	// 构建请求URL
	strInBuffer := "PUT /ISAPI/Intelligent/FDLib/FDSetUp?format=json"

	// 创建ISAPI输入参数
	inputParamPtr := C.malloc(C.size_t(len(strInBuffer) + 1))
	defer C.free(inputParamPtr)
	C.memcpy(inputParamPtr, unsafe.Pointer(&[]byte(strInBuffer)[0]), C.size_t(len(strInBuffer)))

	// 调用远程配置 - 使用与Java一致的NET_DVR_FACE_DATA_RECORD命令码
	lHandle := fm.SDK.NET_DVR_StartRemoteConfig(lUserID, NET_DVR_FACE_DATA_RECORD, inputParamPtr, uint32(len(strInBuffer)), 0, nil)
	if lHandle < 0 {
		return fmt.Errorf("NET_DVR_StartRemoteConfig失败，错误码: %d", fm.SDK.NET_DVR_GetLastError())
	}
	defer fm.SDK.NET_DVR_StopRemoteConfig(lHandle)

	// 构建人脸信息JSON
	strJsonData := fmt.Sprintf(`{
		"faceLibType": "blackFD",
		"FDID": "1",
		"FPID": "%s"
	}`, employeeNo)

	fmt.Printf("下发人脸二进制数据,json data: %s\n", strJsonData)
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
	jsonDataPtr := C.malloc(C.size_t(len(strJsonData)))
	defer C.free(jsonDataPtr)
	C.memcpy(jsonDataPtr, unsafe.Pointer(&[]byte(strJsonData)[0]), C.size_t(len(strJsonData)))
	struAddFaceDataCfg.lpJsonData = jsonDataPtr
	struAddFaceDataCfg.dwJsonDataSize = uint32(len(strJsonData))

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
	fmt.Printf("请求数据大小: %d 字节\n", struAddFaceDataCfg.dwSize)
	result := fm.SDK.NET_DVR_SendWithRecvRemoteConfig(lHandle, unsafe.Pointer(&struAddFaceDataCfg), struAddFaceDataCfg.dwSize, unsafe.Pointer(&outputBuf[0]), uint32(len(outputBuf)), &resultLen)

	fmt.Printf("NET_DVR_SendWithRecvRemoteConfig结果: %d, 返回字节数: %d\n", result, resultLen)

	// 处理配置结果
	if result == -1 {
		errCode := fm.SDK.NET_DVR_GetLastError()
		var errMsg string
		switch errCode {
		case sdk.NET_DVR_NETWORK_FAIL_CONNECT:
			errMsg = "连接服务器失败，请检查网络连接和设备状态"
		case sdk.NET_DVR_NETWORK_SEND_ERROR:
			errMsg = "向服务器发送失败，可能是数据包过大"
		case sdk.NET_DVR_NETWORK_RECV_ERROR:
			errMsg = "从服务器接收数据失败"
		case sdk.NET_DVR_NETWORK_RECV_TIMEOUT:
			errMsg = "从服务器接收数据超时"
		case sdk.NET_DVR_NETWORK_ERRORDATA:
			errMsg = "传送的数据有误"
		default:
			errMsg = fmt.Sprintf("未知错误: %d", errCode)
		}
		return fmt.Errorf("下发人脸失败: 发送ISAPI请求失败，错误码: %d, 错误信息: %s", errCode, errMsg)
	} else if result == sdk.NET_SDK_CONFIG_STATUS_NEED_WAIT {
		// 配置等待，等待一段时间后再次尝试获取结果
		fmt.Println("配置等待，等待获取结果...")
		time.Sleep(100 * time.Millisecond)

		// 再次尝试获取结果，但不发送数据
		emptyData := []byte{}
		for i := 0; i < 20; i++ {
			result = fm.SDK.NET_DVR_SendWithRecvRemoteConfig(lHandle, unsafe.Pointer(&emptyData[0]), 0, unsafe.Pointer(&outputBuf[0]), uint32(len(outputBuf)), &resultLen)

			if result != sdk.NET_SDK_CONFIG_STATUS_NEED_WAIT {
				break
			}

			fmt.Printf("等待配置结果，尝试次数: %d\n", i+1)
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

// AddFaceByUrl 通过URL方式添加人脸图片
func (fm *FaceManage) AddFaceByUrl(lUserID int, employeeNo string) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID")
	}

	// 构建人脸信息JSON
	jsonData := fmt.Sprintf(`{
		"faceLibType": "blackFD",
		"FDID": "1",
		"FPID": "%s",
		"faceURL": "https://oss.dev.zswllife.cn/zswl-dev/665fc0ef9083d203896d3549/store/files/06ac68136bd8414e97554d004aff997c_tmp_8dd206c1e7308826f9be0b3e784c6527b503eaea4c1af0c1.jpg/1.jpg"
	}`, employeeNo)

	// URL
	url := "PUT /ISAPI/Intelligent/FDLib/FaceDataRecord?format=json"
	url = "PUT /ISAPI/Intelligent/FDLib/FDSetUp?format=json"

	// 发送ISAPI请求
	response, err := fm.sendISAPIRequest(lUserID, url, []byte(jsonData))
	if err != nil {
		return err
	}

	fmt.Printf("添加人脸信息成功, 响应: %s\n", string(response))
	return nil
}

// SearchFaceInfo 查询人脸信息
func (fm *FaceManage) SearchFaceInfo(lUserID int, employeeNo string) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID")
	}

	// URL
	url := "POST /ISAPI/Intelligent/FDLib/FDSearch?format=json"

	// 构建查询JSON
	jsonData := fmt.Sprintf(`{
		"searchResultPosition": 0,
		"maxResults": 30,
		"faceLibType": "blackFD",
		"FDID": "1",
		"FPID": "%s"
	}`, employeeNo)

	// 创建ISAPI输入参数
	inputParam := make([]byte, len(url)+1)
	copy(inputParam, []byte(url))

	// 定义类似于NET_DVR_JSON_DATA_CFG的结构体用于处理返回数据
	type NET_DVR_JSON_DATA_CFG struct {
		dwSize                uint32
		lpJsonData            unsafe.Pointer // JSON报文
		dwJsonDataSize        uint32         // JSON报文大小
		lpPicData             unsafe.Pointer // 图片内容
		dwPicDataSize         uint32         // 图片内容大小
		lpInfraredFacePic     unsafe.Pointer // 红外人脸图片数据缓存
		dwInfraredFacePicSize uint32         // 红外人脸图片数据大小
		byRes                 [248]byte      // 保留
	}

	// 调用远程配置 - 使用正确的NET_DVR_FACE_DATA_SEARCH命令码
	lHandle := fm.SDK.NET_DVR_StartRemoteConfig(lUserID, sdk.NET_DVR_FACE_DATA_SEARCH, unsafe.Pointer(&inputParam[0]), uint32(len(inputParam)), 0, nil)
	if lHandle < 0 {
		return fmt.Errorf("NET_DVR_StartRemoteConfig失败，错误码: %d", fm.SDK.NET_DVR_GetLastError())
	}
	defer fm.SDK.NET_DVR_StopRemoteConfig(lHandle)

	// 发送查询数据
	var outputBuf [20 * 1024]byte // 增大缓冲区大小
	var resultLen uint32

	// 创建NET_DVR_JSON_DATA_CFG结构体用于接收数据
	var m_struJsonData NET_DVR_JSON_DATA_CFG

	// 发送请求数据
	jsonDataBytes := []byte(jsonData) // 将字符串转换为字节切片
	result := fm.SDK.NET_DVR_SendWithRecvRemoteConfig(lHandle, unsafe.Pointer(&jsonDataBytes[0]), uint32(len(jsonDataBytes)), unsafe.Pointer(&m_struJsonData), uint32(unsafe.Sizeof(m_struJsonData)), &resultLen)

	// 处理发送结果
	fmt.Printf("NET_DVR_SendWithRecvRemoteConfig结果: %d, 返回字节数: %d\n", result, resultLen)

	// 处理配置结果
	if result == -1 {
		return fmt.Errorf("发送ISAPI请求失败，错误码: %d", fm.SDK.NET_DVR_GetLastError())
	} else if result == sdk.NET_SDK_CONFIG_STATUS_NEED_WAIT {
		// 配置等待，等待一段时间后再次尝试获取结果
		time.Sleep(10 * time.Millisecond)

		// 再次尝试获取结果，但不发送数据
		emptyData := []byte{}
		for i := 0; i < 10; i++ { // 最多尝试10次
			var emptyPtr unsafe.Pointer
			if len(emptyData) > 0 {
				emptyPtr = unsafe.Pointer(&emptyData[0])
			}
			result = fm.SDK.NET_DVR_SendWithRecvRemoteConfig(lHandle, emptyPtr, 0, unsafe.Pointer(&outputBuf[0]), uint32(len(outputBuf)), &resultLen)

			if result != sdk.NET_SDK_CONFIG_STATUS_NEED_WAIT {
				break
			}

			time.Sleep(100 * time.Millisecond)
		}
	}

	// 再次检查结果
	if result == sdk.NET_SDK_CONFIG_STATUS_FAILED {
		return fmt.Errorf("查询人脸信息失败")
	} else if result == sdk.NET_SDK_CONFIG_STATUS_EXCEPTION {
		return fmt.Errorf("查询人脸信息异常")
	} else if result == sdk.NET_SDK_CONFIG_STATUS_SUCCESS || result == sdk.NET_SDK_CONFIG_STATUS_FINISH {
		fmt.Printf("查询人脸信息成功, 数据大小: %d, 图片数据大小: %d, 红外人脸图片数据大小: %d\n",
			m_struJsonData.dwJsonDataSize, m_struJsonData.dwPicDataSize, m_struJsonData.dwInfraredFacePicSize)

		// 处理JSON数据
		var jsonResponse map[string]interface{}
		if m_struJsonData.dwJsonDataSize > 0 && m_struJsonData.lpJsonData != nil {
			// 创建缓冲区接收JSON数据
			jsonBuf := make([]byte, m_struJsonData.dwJsonDataSize)

			// 从lpJsonData指针复制数据到缓冲区
			jsonSlice := (*[1 << 30]byte)(m_struJsonData.lpJsonData)[:m_struJsonData.dwJsonDataSize:m_struJsonData.dwJsonDataSize]
			copy(jsonBuf, jsonSlice)

			// 将字节数据转换为字符串
			jsonStr := string(jsonBuf)
			fmt.Printf("JSON数据: %s\n", jsonStr)

			// 解析JSON响应
			if err := json.Unmarshal(jsonBuf, &jsonResponse); err != nil {
				fmt.Printf("解析响应JSON失败: %v\n", err)
				return nil
			}

			// 获取匹配数量
			numOfMatches, ok := jsonResponse["numOfMatches"].(float64)
			if !ok {
				fmt.Println("未找到匹配数量信息，响应可能不符合预期格式")
				return nil
			}

			fmt.Printf("查询到 %d 条匹配的人脸记录\n", int(numOfMatches))

			// 如果有匹配的人脸
			if numOfMatches > 0 {
				// 获取匹配列表
				matchList, ok := jsonResponse["MatchList"].([]interface{})
				if !ok || len(matchList) == 0 {
					fmt.Println("无法获取匹配列表或列表为空")
					return nil
				}

				// 获取第一个匹配的人脸信息
				matchInfo, ok := matchList[0].(map[string]interface{})
				if !ok {
					fmt.Println("无法获取人脸信息")
					return nil
				}

				// 获取工号
				fpid, ok := matchInfo["FPID"].(string)
				if !ok {
					fpid = employeeNo // 如果获取不到，使用传入的工号
					fmt.Printf("未找到FPID，使用传入的工号: %s\n", fpid)
				}

				// 确保pic目录存在
				err := os.MkdirAll("./pic", 0755)
				if err != nil {
					fmt.Printf("创建图片目录失败: %v\n", err)
					return nil
				}

				// 处理普通人脸图片
				imagePath := fmt.Sprintf("./pic/[%s]_FacePic.jpg", fpid)

				// 保存普通图片数据（如果有）
				if m_struJsonData.dwPicDataSize > 0 && m_struJsonData.lpPicData != nil {
					// 创建缓冲区接收普通人脸图片数据
					picBuf := make([]byte, m_struJsonData.dwPicDataSize)

					// 从lpPicData指针复制数据到缓冲区
					picSlice := (*[1 << 30]byte)(m_struJsonData.lpPicData)[:m_struJsonData.dwPicDataSize:m_struJsonData.dwPicDataSize]
					copy(picBuf, picSlice)

					// 保存普通人脸图片数据到文件
					err = os.WriteFile(imagePath, picBuf, 0644)
					if err != nil {
						fmt.Printf("保存普通人脸图片数据失败: %v\n", err)
					} else {
						fmt.Printf("普通人脸图片已保存到: %s\n", imagePath)
					}
				} else {
					fmt.Println("没有收到普通人脸图片数据，或数据大小为0")

					// 尝试从响应中获取图片数据
					if picData, hasPic := jsonResponse["pic"].(string); hasPic {
						imageBytes, err := base64.StdEncoding.DecodeString(picData)
						if err == nil {
							err = os.WriteFile(imagePath, imageBytes, 0644)
							if err != nil {
								fmt.Printf("从JSON保存普通人脸图片数据失败: %v\n", err)
							} else {
								fmt.Printf("普通人脸图片(Base64)已保存到: %s\n", imagePath)
							}
						} else {
							fmt.Printf("普通人脸图片Base64解码失败: %v\n", err)
						}
					}
				}

				// 处理红外人脸图片数据（如果有）
				if m_struJsonData.dwInfraredFacePicSize > 0 && m_struJsonData.lpInfraredFacePic != nil {
					// 创建缓冲区接收红外人脸图片数据
					irPicBuf := make([]byte, m_struJsonData.dwInfraredFacePicSize)

					// 从lpInfraredFacePic指针复制数据到缓冲区
					irPicSlice := (*[1 << 30]byte)(m_struJsonData.lpInfraredFacePic)[:m_struJsonData.dwInfraredFacePicSize:m_struJsonData.dwInfraredFacePicSize]
					copy(irPicBuf, irPicSlice)

					// 保存红外人脸图片数据到文件
					irImagePath := fmt.Sprintf("./pic/[%s]_InfraredFacePic.jpg", fpid)
					err = os.WriteFile(irImagePath, irPicBuf, 0644)
					if err != nil {
						fmt.Printf("保存红外人脸图片数据失败: %v\n", err)
					} else {
						fmt.Printf("红外人脸图片已保存到: %s\n", irImagePath)
					}
				} else {
					fmt.Println("没有收到红外人脸图片数据，或数据大小为0")
				}

				fmt.Printf("成功查询到人脸信息，工号: %s\n", fpid)
			} else {
				fmt.Println("未找到匹配的人脸信息")
			}
		} else {
			fmt.Println("没有有效的JSON数据")
		}
		return nil
	}

	return fmt.Errorf("查询人脸信息未知错误，结果码: %d", result)
}

func (fm *FaceManage) sendISAPIRequest(lUserID int, url string, requestData []byte) ([]byte, error) {
	// 创建ISAPI输入参数
	inputParam := make([]byte, len(url)+1)
	copy(inputParam, []byte(url))

	// 增大缓冲区大小，解决错误码43问题
	var outputBuf [100 * 1024]byte // 增大到100KB，原来是20KB

	// 调用远程配置
	lHandle := fm.SDK.NET_DVR_StartRemoteConfig(lUserID, sdk.NET_DVR_JSON_CONFIG, unsafe.Pointer(&inputParam[0]), uint32(len(inputParam)), 0, nil)
	if lHandle < 0 {
		return nil, fmt.Errorf("NET_DVR_StartRemoteConfig失败，错误码: %d", fm.SDK.NET_DVR_GetLastError())
	}
	defer fm.SDK.NET_DVR_StopRemoteConfig(lHandle)

	// 发送数据
	var resultLen uint32
	var result int

	// 检查请求数据大小
	if requestData != nil {
		fmt.Printf("请求数据大小: %d 字节\n", len(requestData))
	}

	// 只发送一次请求，不要在循环中重复发送
	if requestData == nil || len(requestData) == 0 {
		// 如果没有请求数据，只发送空数据
		emptyData := []byte{}
		result = fm.SDK.NET_DVR_SendWithRecvRemoteConfig(lHandle, unsafe.Pointer(&emptyData[0]), 0, unsafe.Pointer(&outputBuf[0]), uint32(len(outputBuf)), &resultLen)
	} else {
		// 发送请求数据
		result = fm.SDK.NET_DVR_SendWithRecvRemoteConfig(lHandle, unsafe.Pointer(&requestData[0]), uint32(len(requestData)), unsafe.Pointer(&outputBuf[0]), uint32(len(outputBuf)), &resultLen)
	}

	// 处理发送结果
	fmt.Printf("NET_DVR_SendWithRecvRemoteConfig结果: %d, 返回字节数: %d\n", result, resultLen)

	// 处理配置结果
	if result == -1 {
		// 获取详细错误信息
		errCode := fm.SDK.NET_DVR_GetLastError()
		var errMsg string
		switch errCode {
		case sdk.NET_DVR_NETWORK_FAIL_CONNECT:
			errMsg = "连接服务器失败，请检查网络连接和设备状态"
		case sdk.NET_DVR_NETWORK_SEND_ERROR:
			errMsg = "向服务器发送失败，可能是数据包过大"
		case sdk.NET_DVR_NETWORK_RECV_ERROR:
			errMsg = "从服务器接收数据失败"
		case sdk.NET_DVR_NETWORK_RECV_TIMEOUT:
			errMsg = "从服务器接收数据超时"
		case sdk.NET_DVR_NETWORK_ERRORDATA:
			errMsg = "传送的数据有误"
		default:
			errMsg = fmt.Sprintf("未知错误: %d", errCode)
		}
		return nil, fmt.Errorf("发送ISAPI请求失败，错误码: %d, 错误信息: %s", errCode, errMsg)
	} else if result == sdk.NET_SDK_CONFIG_STATUS_NEED_WAIT {
		// 配置等待，等待一段时间后再次尝试获取结果
		fmt.Println("配置等待，等待获取结果...")
		time.Sleep(100 * time.Millisecond) // 增加等待时间

		// 再次尝试获取结果，但不发送数据
		emptyData := []byte{}
		for i := 0; i < 20; i++ { // 增加尝试次数，原来是10次
			result = fm.SDK.NET_DVR_SendWithRecvRemoteConfig(lHandle, unsafe.Pointer(&emptyData[0]), 0, unsafe.Pointer(&outputBuf[0]), uint32(len(outputBuf)), &resultLen)

			if result != sdk.NET_SDK_CONFIG_STATUS_NEED_WAIT {
				break
			}

			fmt.Printf("等待配置结果，尝试次数: %d\n", i+1)
			time.Sleep(200 * time.Millisecond) // 增加等待时间
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

// DeleteFaceInfo 删除人脸信息
func (fm *FaceManage) DeleteFaceInfo(lUserID int, employeeNo string) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID")
	}

	// URL
	url := "PUT /ISAPI/Intelligent/FDLib/FDSearch/Delete?format=json&FDID=1&faceLibType=blackFD"

	// 构建删除人脸JSON
	jsonData := fmt.Sprintf(`{
		"FPID": [{
        	"value": "%s"
    	}]
	}`, employeeNo)

	fmt.Printf("准备删除工号: %s 的人脸信息\n", employeeNo)
	// 发送ISAPI请求
	response, err := fm.sendISAPIRequest(lUserID, url, []byte(jsonData))
	if err != nil {
		return err
	}

	fmt.Printf("删除人脸信息成功, 响应: %s\n", string(response))
	return nil
}

// CaptureFaceInfo 采集人脸信息
func (fm *FaceManage) CaptureFaceInfo(lUserID int) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID")
	}

	// 定义采集人脸条件结构体(与Java实现保持一致)
	type NET_DVR_CAPTURE_FACE_COND struct {
		dwSize uint32
		byRes  [128]byte // 调整为128字节，与Java实现一致
	}

	// 初始化采集人脸条件结构体
	var struCapCond NET_DVR_CAPTURE_FACE_COND
	struCapCond.dwSize = uint32(unsafe.Sizeof(struCapCond))

	// 调用远程配置
	// 注意：Java实现中回调函数和用户数据都为null
	lHandle := fm.SDK.NET_DVR_StartRemoteConfig(lUserID, sdk.NET_DVR_CAPTURE_FACE_INFO, unsafe.Pointer(&struCapCond), struCapCond.dwSize, 0, nil)
	if lHandle < 0 {
		return fmt.Errorf("建立采集人脸长连接失败，错误码: %d", fm.SDK.NET_DVR_GetLastError())
	}
	fmt.Println("建立采集人脸长连接成功!")
	defer fm.SDK.NET_DVR_StopRemoteConfig(lHandle)

	// 定义人脸特征结构体（与Java实现保持一致）
	type NET_VCA_RECT struct {
		fX      float32 // 边界框左上角的X坐标
		fY      float32 // 边界框左上角的Y坐标
		fWidth  float32 // 边界框的宽度
		fHeight float32 // 边界框的高度
	}

	type NET_VCA_POINT struct {
		fX float32 // X坐标
		fY float32 // Y坐标
	}

	type NET_DVR_FACE_FEATURE struct {
		struFace       NET_VCA_RECT  // 人脸子图区域
		struLeftEye    NET_VCA_POINT // 左眼坐标
		struRightEye   NET_VCA_POINT // 右眼坐标
		struLeftMouth  NET_VCA_POINT // 嘴左边坐标
		struRightMouth NET_VCA_POINT // 嘴右边坐标
		struNoseTip    NET_VCA_POINT // 鼻子坐标
	}

	// 定义采集人脸信息结构体（与Java实现保持一致）
	type NET_DVR_CAPTURE_FACE_CFG struct {
		dwSize                   uint32
		dwFaceTemplate1Size      uint32
		pFaceTemplate1Buffer     unsafe.Pointer
		dwFaceTemplate2Size      uint32
		pFaceTemplate2Buffer     unsafe.Pointer
		dwFacePicSize            uint32
		pFacePicBuffer           unsafe.Pointer
		byFaceQuality1           byte                 // 人脸质量，范围1-100
		byFaceQuality2           byte                 // 人脸质量，范围1-100
		byCaptureProgress        byte                 // 采集进度，0-未采集到人脸，100-采集到人脸
		byFacePicQuality         byte                 // 人脸图片中人脸质量
		dwInfraredFacePicSize    uint32               // 红外人脸图片数据大小
		pInfraredFacePicBuffer   unsafe.Pointer       // 红外人脸图片数据缓存
		byInfraredFacePicQuality byte                 // 红外人脸图片中人脸质量
		byRes1                   [3]byte              // 保留字节1
		struFeature              NET_DVR_FACE_FEATURE // 人脸抠图特征信息
		byRes                    [56]byte             // 保留字节，确保与Java实现一致
	}

	// 初始化采集人脸信息结构体
	var struFaceInfo NET_DVR_CAPTURE_FACE_CFG
	struFaceInfo.dwSize = uint32(unsafe.Sizeof(struFaceInfo))

	// 使用C分配的内存
	// 人脸模板1缓冲区
	template1Buf := C.malloc(C.size_t(2500)) // 不大于2.5K
	defer C.free(template1Buf)
	struFaceInfo.pFaceTemplate1Buffer = template1Buf

	// 人脸模板2缓冲区
	template2Buf := C.malloc(C.size_t(2500)) // 不大于2.5K
	defer C.free(template2Buf)
	struFaceInfo.pFaceTemplate2Buffer = template2Buf

	// 人脸图片缓冲区
	facePicBuf := C.malloc(C.size_t(200 * 1024)) // 预留200K
	defer C.free(facePicBuf)
	struFaceInfo.pFacePicBuffer = facePicBuf

	// 红外人脸图片缓冲区
	infraredFacePicBuf := C.malloc(C.size_t(200 * 1024)) // 预留200K
	defer C.free(infraredFacePicBuf)
	struFaceInfo.pInfraredFacePicBuffer = infraredFacePicBuf

	// 采集人脸信息
	for {
		var resultLen uint32
		dwState := fm.SDK.NET_DVR_GetNextRemoteConfig(lHandle, unsafe.Pointer(&struFaceInfo), struFaceInfo.dwSize, &resultLen)

		if dwState == -1 {
			return fmt.Errorf("NET_DVR_GetNextRemoteConfig采集人脸失败，错误码: %d", fm.SDK.NET_DVR_GetLastError())
		} else if dwState == sdk.NET_SDK_CONFIG_STATUS_NEED_WAIT {
			fmt.Println("正在采集中,请等待...")
			time.Sleep(10 * time.Millisecond)
			continue
		} else if dwState == sdk.NET_SDK_CONFIG_STATUS_FAILED {
			return fmt.Errorf("采集人脸失败")
		} else if dwState == sdk.NET_SDK_CONFIG_STATUS_EXCEPTION {
			return fmt.Errorf("采集人脸异常, 网络异常导致连接断开")
		} else if dwState == sdk.NET_SDK_CONFIG_STATUS_SUCCESS {
			if (struFaceInfo.dwFacePicSize > 0) && (struFaceInfo.pFacePicBuffer != nil) {
				// 生成时间戳作为文件名
				timeStamp := time.Now().Format("20060102150405")
				filename := fmt.Sprintf("./pic/%s_capFaceInfo.jpg", timeStamp)

				// 确保pic目录存在
				err := os.MkdirAll("./pic", 0755)
				if err != nil {
					return fmt.Errorf("创建图片目录失败: %v", err)
				}

				// 从C内存复制人脸图片数据到Go切片
				picBuf := make([]byte, struFaceInfo.dwFacePicSize)
				C.memcpy(unsafe.Pointer(&picBuf[0]), struFaceInfo.pFacePicBuffer, C.size_t(struFaceInfo.dwFacePicSize))

				err = os.WriteFile(filename, picBuf, 0644)
				if err != nil {
					return fmt.Errorf("保存人脸图片失败: %v", err)
				}

				fmt.Printf("采集人脸成功, 图片保存路径: %s\n", filename)
			}
			break
		} else {
			fmt.Printf("其他异常, dwState: %d\n", dwState)
			break
		}
	}

	fmt.Println("采集人脸操作完成")
	return nil
}
