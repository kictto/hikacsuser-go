package models

import (
	"fmt"
	"github.com/hikacsuser-go/internal/sdk"
	"strings"
	"time"
	"unsafe"
)

// NET_DVR_DEVICEINFO_COND 设备参数结构条件
type NET_DVR_DEVICEINFO_COND struct {
	DwSize        uint32               // 结构体大小
	SUserName     [sdk.NAME_LEN]byte   // 用户名
	SPassword     [sdk.PASSWD_LEN]byte // 用户密码
	ByHttps       byte                 // 是否启用HTTPS 0-不启用 1-启用
	IDevInfoIndex byte                 // 设备信息索引
	ByRes         [62]byte             // 保留
}

// NET_DVR_FACE_COND 人脸信息查询条件
type NET_DVR_FACE_COND struct {
	DwSize      uint32                            // 结构体大小
	ByCardNo    [sdk.MAX_CARDNO_LEN]byte          // 人脸关联的卡号
	IEmployeeNo [sdk.NET_SDK_EMPLOYEE_NO_LEN]byte // 工号
	DwFaceNum   uint32                            // 人脸数量，获取全部时为0xffffffff
	STerminalNo [sdk.MAX_CARDNO_LEN]byte          // 终端编号
	ByRes       [188]byte                         // 保留
}

// UserManage 用户管理
type UserManage struct {
	SDK sdk.HCNetSDK
}

// NewUserManage 创建用户管理实例
func NewUserManage(sdk sdk.HCNetSDK) *UserManage {
	return &UserManage{
		SDK: sdk,
	}
}

// AddUserInfo 添加用户信息
func (um *UserManage) AddUserInfo(lUserID int, employeeNo string) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID")
	}

	// 构建用户信息JSON
	jsonData := fmt.Sprintf(`{
		"UserInfo": {
			"employeeNo": "%s",
			"name": "测试用户",
			"userType": "normal",
			"Valid": {
				"enable": true,
				"beginTime": "2021-01-01T00:00:00",
				"endTime": "2031-01-01T00:00:00",
				"timeType": "local"
			},
			"belongGroup": "1",
			"doorRight": "1",
			"password": "",
			"RightPlan": [
				{
					"doorNo": 1,
					"planTemplateNo": "1"
				}
			],
			"maxOpenDoorTime": 0,
			"openDoorTime": 0,
			"localUIRight": false,
			"userVerifyMode": "cardOrFace"
		}
	}`, employeeNo)

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

// SearchUserInfo 查询用户信息
func (um *UserManage) SearchUserInfo(lUserID int) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID")
	}

	// 生成UUID作为searchID，与Java实现保持一致
	uuid := fmt.Sprintf("%d", time.Now().UnixNano())

	// 查询条件JSON，参照Java实现
	jsonData := fmt.Sprintf(`{
		"UserInfoSearchCond": {
			"searchID": "%s",
			"searchResultPosition": 0,
			"maxResults": 30
		}
	}`, uuid)

	// 打印查询JSON
	fmt.Printf("查询的json报文: %s\n", jsonData)

	// URL
	url := "POST /ISAPI/AccessControl/UserInfo/Search?format=json"

	// 发送ISAPI请求
	response, err := um.sendISAPIRequest(lUserID, url, []byte(jsonData))
	if err != nil {
		return err
	}

	fmt.Printf("查询用户信息成功, 响应: %s\n", string(response))
	return nil
}

// DeleteUserInfo 删除用户信息
func (um *UserManage) DeleteUserInfo(lUserID int) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID")
	}

	var input string
	fmt.Print("请选择删除方式 1-全部删除 2-根据工号删除:")
	fmt.Scanln(&input)
	input = strings.TrimSpace(input)

	var url string
	var jsonData string

	if input == "1" {
		// 删除所有用户
		url = "PUT /ISAPI/AccessControl/UserInfo/Delete?format=json"
		jsonData = `{
			"UserInfoDelCond": {
				"OperateType": "clear"
			}
		}`
	} else if input == "2" {
		// 按照工号删除
		var employeeNo string
		fmt.Print("请输入需要删除的工号:")
		fmt.Scanln(&employeeNo)
		employeeNo = strings.TrimSpace(employeeNo)

		url = "PUT /ISAPI/AccessControl/UserInfo/Delete?format=json"
		jsonData = fmt.Sprintf(`{
			"UserInfoDelCond": {
				"OperateType": "delete",
				"EmployeeNoList": [
					{
						"employeeNo": "%s"
					}
				]
			}
		}`, employeeNo)
	} else {
		return fmt.Errorf("无效的操作类型")
	}

	// 发送ISAPI请求
	response, err := um.sendISAPIRequest(lUserID, url, []byte(jsonData))
	if err != nil {
		return err
	}

	fmt.Printf("删除用户信息成功, 响应: %s\n", string(response))
	return nil
}

// SetCardTemplate 设置卡片模板
func (um *UserManage) SetCardTemplate(lUserID int, templateNo int) error {
	if lUserID < 0 {
		return fmt.Errorf("无效的用户ID")
	}

	// 构建模板JSON
	// 这里设置的是全天候24小时均可开门的权限计划
	jsonData := fmt.Sprintf(`{
		"PlanTemplate": {
			"enableCardReader": [1],
			"templateNo": %d,
			"templateName": "全天有效",
			"weekPlanNo": 1,
			"holidayGroupNo": 0
		}
	}`, templateNo)

	// URL
	url := "PUT /ISAPI/AccessControl/PlantTemplate/Template?format=json"

	// 发送ISAPI请求
	response, err := um.sendISAPIRequest(lUserID, url, []byte(jsonData))
	if err != nil {
		return err
	}

	// 设置周计划
	weekPlanJson := `{
		"WeekPlan": {
			"planNo": 1,
			"planName": "全天有效",
			"dayPlans": [
				{
					"dayNo": 1,
					"dayPlanNo": 1
				},
				{
					"dayNo": 2,
					"dayPlanNo": 1
				},
				{
					"dayNo": 3,
					"dayPlanNo": 1
				},
				{
					"dayNo": 4,
					"dayPlanNo": 1
				},
				{
					"dayNo": 5,
					"dayPlanNo": 1
				},
				{
					"dayNo": 6,
					"dayPlanNo": 1
				},
				{
					"dayNo": 7,
					"dayPlanNo": 1
				}
			]
		}
	}`

	// URL
	weekPlanUrl := "PUT /ISAPI/AccessControl/WeekPlan/Template/1?format=json"

	// 发送ISAPI请求
	weekPlanResponse, err := um.sendISAPIRequest(lUserID, weekPlanUrl, []byte(weekPlanJson))
	if err != nil {
		return err
	}

	// 设置日计划
	dayPlanJson := `{
		"DayPlan": {
			"planNo": 1,
			"planName": "全天有效",
			"timeSections": [
				{
					"sectionNo": 1,
					"startTime": "00:00:00",
					"endTime": "23:59:59"
				}
			]
		}
	}`

	// URL
	dayPlanUrl := "PUT /ISAPI/AccessControl/DayPlan/Template/1?format=json"

	// 发送ISAPI请求
	dayPlanResponse, err := um.sendISAPIRequest(lUserID, dayPlanUrl, []byte(dayPlanJson))
	if err != nil {
		return err
	}

	fmt.Printf("设置卡片模板成功, 响应: %s\n", string(response))
	fmt.Printf("设置周计划成功, 响应: %s\n", string(weekPlanResponse))
	fmt.Printf("设置日计划成功, 响应: %s\n", string(dayPlanResponse))
	return nil
}

// sendISAPIRequest 发送ISAPI请求
func (um *UserManage) sendISAPIRequest(lUserID int, url string, requestData []byte) ([]byte, error) {
	// 创建ISAPI输入参数
	inputParam := make([]byte, len(url)+1)
	copy(inputParam, []byte(url))

	var outputBuf [20 * 1024]byte // 增大缓冲区大小，与Java实现保持一致

	// 调用远程配置
	lHandle := um.SDK.NET_DVR_StartRemoteConfig(lUserID, sdk.NET_DVR_JSON_CONFIG, unsafe.Pointer(&inputParam[0]), uint32(len(inputParam)), 0, nil)
	if lHandle < 0 {
		return nil, fmt.Errorf("NET_DVR_StartRemoteConfig失败，错误码: %d", um.SDK.NET_DVR_GetLastError())
	}
	defer um.SDK.NET_DVR_StopRemoteConfig(lHandle)

	// 发送数据
	var resultLen uint32
	var result int

	// 只发送一次请求，不要在循环中重复发送
	if requestData == nil || len(requestData) == 0 {
		// 如果没有请求数据，只发送空数据
		emptyData := []byte{}
		result = um.SDK.NET_DVR_SendWithRecvRemoteConfig(lHandle, unsafe.Pointer(&emptyData[0]), 0, unsafe.Pointer(&outputBuf[0]), uint32(len(outputBuf)), &resultLen)
	} else {
		// 发送请求数据
		result = um.SDK.NET_DVR_SendWithRecvRemoteConfig(lHandle, unsafe.Pointer(&requestData[0]), uint32(len(requestData)), unsafe.Pointer(&outputBuf[0]), uint32(len(outputBuf)), &resultLen)
	}

	// 处理发送结果
	fmt.Printf("NET_DVR_SendWithRecvRemoteConfig结果: %d, 返回字节数: %d\n", result, resultLen)

	// 处理配置结果
	if result == -1 {
		return nil, fmt.Errorf("发送ISAPI请求失败，错误码: %d", um.SDK.NET_DVR_GetLastError())
	} else if result == sdk.NET_SDK_CONFIG_STATUS_NEED_WAIT {
		// 配置等待，等待一段时间后再次尝试获取结果
		time.Sleep(10 * time.Millisecond)

		// 再次尝试获取结果，但不发送数据
		emptyData := []byte{}
		for i := 0; i < 10; i++ { // 最多尝试10次
			result = um.SDK.NET_DVR_SendWithRecvRemoteConfig(lHandle, unsafe.Pointer(&emptyData[0]), 0, unsafe.Pointer(&outputBuf[0]), uint32(len(outputBuf)), &resultLen)

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
