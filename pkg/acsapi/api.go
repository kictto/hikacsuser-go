// Package acsapi 提供海康威视门禁系统的公共API接口
// 该包将内部实现封装为公共方法，支持批量和单个数据操作
package acsapi

import (
	"fmt"
	"github.com/clockworkchen/hikacsuser-go/internal/models"
	"github.com/clockworkchen/hikacsuser-go/internal/sdk"
	"time"
)

// ACSClient 门禁系统客户端
type ACSClient struct {
	hcnetsdk    sdk.HCNetSDK
	lUserID     int
	deviceIP    string
	devicePort  uint16
	acsManage   *models.ACSManage
	userManage  *models.UserManage
	cardManage  *models.CardManage
	faceManage  *models.FaceManage
	eventSearch *models.EventSearch
	alarmManage *models.AlarmManage
}

// NewACSClient 创建新的门禁系统客户端
func NewACSClient(logPath string) (*ACSClient, error) {
	// 获取SDK实例
	instance, err := sdk.GetSDKInstance()
	if err != nil {
		return nil, fmt.Errorf("获取SDK实例失败: %v", err)
	}

	// 创建客户端
	client := &ACSClient{
		hcnetsdk: instance,
		lUserID:  -1,
	}

	// 设置日志
	if !client.hcnetsdk.NET_DVR_SetLogToFile(3, logPath, false) {
		return nil, fmt.Errorf("设置日志失败: %d", client.hcnetsdk.NET_DVR_GetLastError())
	}

	// 初始化管理模块
	client.acsManage = models.NewACSManage(client.hcnetsdk)
	client.userManage = models.NewUserManage(client.hcnetsdk)
	client.cardManage = models.NewCardManage(client.hcnetsdk)
	client.faceManage = models.NewFaceManage(client.hcnetsdk)
	client.eventSearch = models.NewEventSearch(client.hcnetsdk)

	return client, nil
}

// Login 登录设备
// 返回登录信息和设备信息供调用方使用
func (c *ACSClient) Login(ip string, port uint16, username, password string) (LoginInfo, DeviceInfo, error) {
	// 创建登录信息
	var loginInfo sdk.NET_DVR_USER_LOGIN_INFO
	var deviceInfo sdk.NET_DVR_DEVICEINFO_V40

	// 设置设备IP地址
	copy(loginInfo.SDeviceAddress[:], []byte(ip))

	// 传输模式设置
	loginInfo.ByUseTransport = 0 // 默认为0

	// 设置端口
	loginInfo.WPort = port

	// 设置用户名和密码
	copy(loginInfo.SUserName[:], []byte(username))
	copy(loginInfo.SPassword[:], []byte(password))

	// 登录模式设置
	loginInfo.BUseAsynLogin = false // 同步登录
	loginInfo.ByLoginMode = 0       // 使用SDK私有协议
	loginInfo.ByHttps = 0           // 不使用HTTPS

	// 执行登录
	userID := c.hcnetsdk.NET_DVR_Login_V40(&loginInfo, &deviceInfo)
	if userID == -1 {
		return LoginInfo{}, DeviceInfo{}, fmt.Errorf("登录失败，错误码为: %d", c.hcnetsdk.NET_DVR_GetLastError())
	}

	// 保存设备信息和用户ID
	c.lUserID = userID
	c.deviceIP = ip
	c.devicePort = port

	// 创建返回给调用方的登录信息
	loginInfoExport := LoginInfo{
		DeviceIP:   ip,
		DevicePort: port,
		Username:   username,
		Password:   password,
	}

	// 创建返回给调用方的设备信息
	deviceInfoExport := DeviceInfo{
		SerialNumber:      sdk.BytesToString(deviceInfo.DeviceInfo.SSerialNumber[:]),
		AlarmInPortNum:    deviceInfo.DeviceInfo.ByAlarmInPortNum,
		AlarmOutPortNum:   deviceInfo.DeviceInfo.ByAlarmOutPortNum,
		DiskNum:           deviceInfo.DeviceInfo.ByDiskNum,
		DeviceType:        deviceInfo.DeviceInfo.ByDVRType,
		ChannelNum:        deviceInfo.DeviceInfo.ByChanNum,
		StartChannel:      deviceInfo.DeviceInfo.ByStartChan,
		IPChannelNum:      deviceInfo.DeviceInfo.ByIPChanNum,
		AudioChannelNum:   deviceInfo.DeviceInfo.ByAudioChanNum,
		IPAlarmInPortNum:  deviceInfo.DeviceInfo.ByIPAlarmInPortNum,
		IPAlarmOutPortNum: deviceInfo.DeviceInfo.ByIPAlarmOutPortNum,
		ZeroChannelNum:    deviceInfo.DeviceInfo.ByZeroChanNum,
		PasswordLevel:     deviceInfo.ByPasswordLevel,
		OEMCode:           deviceInfo.DwOEMCode,
	}

	return loginInfoExport, deviceInfoExport, nil
}

// Logout 注销设备
func (c *ACSClient) Logout() error {
	if c.lUserID >= 0 {
		if !c.hcnetsdk.NET_DVR_Logout(c.lUserID) {
			return fmt.Errorf("设备注销失败，错误码：%d", c.hcnetsdk.NET_DVR_GetLastError())
		}
		c.lUserID = -1
	}
	return nil
}

// Cleanup 清理SDK资源
func (c *ACSClient) Cleanup() {
	c.hcnetsdk.NET_DVR_Cleanup()
}

// GetACSConfig 获取门禁参数
func (c *ACSClient) GetACSConfig() (ResponseData, error) {
	if c.lUserID < 0 {
		return ResponseData{}, fmt.Errorf("未登录设备")
	}

	// 调用内部方法
	err := c.acsManage.AcsCfg(c.lUserID)
	if err != nil {
		return ResponseData{}, err
	}

	// 获取响应数据
	responseData, err := c.acsManage.SendISAPIRequest(c.lUserID, "GET /ISAPI/AccessControl/AcsCfg?format=json", nil)
	if err != nil {
		return ResponseData{}, err
	}

	// 解析响应数据
	response, err := ParseResponseData(responseData)
	if err != nil {
		return ResponseData{}, err
	}

	return response, nil
}

// GetACSStatus 获取门禁状态
func (c *ACSClient) GetACSStatus() (ResponseData, error) {
	if c.lUserID < 0 {
		return ResponseData{}, fmt.Errorf("未登录设备")
	}

	// 调用内部方法
	err := c.acsManage.GetAcsStatus(c.lUserID)
	if err != nil {
		return ResponseData{}, err
	}

	// 获取响应数据
	responseData, err := c.acsManage.SendISAPIRequest(c.lUserID, "GET /ISAPI/AccessControl/AcsStatus?format=json", nil)
	if err != nil {
		return ResponseData{}, err
	}

	// 解析响应数据
	response, err := ParseResponseData(responseData)
	if err != nil {
		return ResponseData{}, err
	}

	return response, nil
}

// RemoteControlGate 远程控门
func (c *ACSClient) RemoteControlGate() (ResponseData, error) {
	if c.lUserID < 0 {
		return ResponseData{}, fmt.Errorf("未登录设备")
	}

	// 调用内部方法
	err := c.acsManage.RemoteControlGate(c.lUserID)
	if err != nil {
		return ResponseData{}, err
	}

	// 获取响应数据
	responseData, err := c.acsManage.SendISAPIRequest(c.lUserID, "PUT /ISAPI/AccessControl/RemoteControl/door/1?format=json", nil)
	if err != nil {
		return ResponseData{}, err
	}

	// 解析响应数据
	response, err := ParseResponseData(responseData)
	if err != nil {
		return ResponseData{}, err
	}

	return response, nil
}

// AddUser 添加单个用户
func (c *ACSClient) AddUser(employeeNo string) (ResponseData, error) {
	if c.lUserID < 0 {
		return ResponseData{}, fmt.Errorf("未登录设备")
	}

	// 调用内部方法
	err := c.userManage.AddUserInfo(c.lUserID, employeeNo)
	if err != nil {
		return ResponseData{}, err
	}

	// 构建用户信息JSON
	jsonData := fmt.Sprintf(`{
		"UserInfo": {
			"employeeNo": "%s"
		}
	}`, employeeNo)

	// 获取响应数据
	responseData, err := c.userManage.SendISAPIRequest(c.lUserID, "POST /ISAPI/AccessControl/UserInfo/Record?format=json", []byte(jsonData))
	if err != nil {
		return ResponseData{}, err
	}

	// 解析响应数据
	response, err := ParseResponseData(responseData)
	if err != nil {
		return ResponseData{}, err
	}

	return response, nil
}

// AddUsers 批量添加用户
func (c *ACSClient) AddUsers(employeeNos []string) ([]ResponseData, []error) {
	if c.lUserID < 0 {
		return nil, []error{fmt.Errorf("未登录设备")}
	}

	errors := make([]error, 0)
	responses := make([]ResponseData, 0)

	for _, employeeNo := range employeeNos {
		response, err := c.AddUser(employeeNo)
		responses = append(responses, response)

		if err != nil {
			errors = append(errors, fmt.Errorf("添加用户 %s 失败: %v", employeeNo, err))
		}
	}

	return responses, errors
}

// SearchUser 查询用户信息
func (c *ACSClient) SearchUser() (ResponseData, error) {
	if c.lUserID < 0 {
		return ResponseData{}, fmt.Errorf("未登录设备")
	}

	// 调用内部方法
	err := c.userManage.SearchUserInfo(c.lUserID)
	if err != nil {
		return ResponseData{}, err
	}

	// 获取响应数据
	responseData, err := c.userManage.SendISAPIRequest(c.lUserID, "POST /ISAPI/AccessControl/UserInfo/Search?format=json", nil)
	if err != nil {
		return ResponseData{}, err
	}

	// 解析响应数据
	response, err := ParseResponseData(responseData)
	if err != nil {
		return ResponseData{}, err
	}

	return response, nil
}

// DeleteUser 方法已移至user.go

// DeleteAllUsers 方法已移至user.go

// AddCard 添加单个卡片
func (c *ACSClient) AddCard(employeeNo, cardNo string) (ResponseData, error) {
	if c.lUserID < 0 {
		return ResponseData{}, fmt.Errorf("未登录设备")
	}

	// 调用内部方法
	err := c.cardManage.AddCardInfo(c.lUserID, employeeNo, cardNo)
	if err != nil {
		return ResponseData{}, err
	}

	// 构建卡片信息JSON
	jsonData := fmt.Sprintf(`{
		"CardInfo": {
			"employeeNo": "%s",
			"cardNo": "%s",
			"cardType": "normalCard"
		}
	}`, employeeNo, cardNo)

	// 获取响应数据
	responseData, err := c.cardManage.SendISAPIRequest(c.lUserID, "POST /ISAPI/AccessControl/CardInfo/Record?format=json", []byte(jsonData))
	if err != nil {
		return ResponseData{}, err
	}

	// 解析响应数据
	response, err := ParseResponseData(responseData)
	if err != nil {
		return ResponseData{}, err
	}

	return response, nil
}

// AddCards 批量添加卡片
func (c *ACSClient) AddCards(employeeNos, cardNos []string) ([]ResponseData, []error) {
	if c.lUserID < 0 {
		return nil, []error{fmt.Errorf("未登录设备")}
	}

	if len(employeeNos) != len(cardNos) {
		return nil, []error{fmt.Errorf("员工号和卡号数量不匹配")}
	}

	errors := make([]error, 0)
	responses := make([]ResponseData, 0)

	for i, employeeNo := range employeeNos {
		response, err := c.AddCard(employeeNo, cardNos[i])
		responses = append(responses, response)

		if err != nil {
			errors = append(errors, fmt.Errorf("添加卡片 %s 失败: %v", cardNos[i], err))
		}
	}

	return responses, errors
}

// SearchCard 查询卡片信息
func (c *ACSClient) SearchCard(employeeNo string) (ResponseData, error) {
	if c.lUserID < 0 {
		return ResponseData{}, fmt.Errorf("未登录设备")
	}

	// 调用内部方法
	err := c.cardManage.SearchCardInfo(c.lUserID, employeeNo)
	if err != nil {
		return ResponseData{}, err
	}

	// 构建查询JSON
	uuid := fmt.Sprintf("%d", time.Now().UnixNano())
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

	// 获取响应数据
	responseData, err := c.cardManage.SendISAPIRequest(c.lUserID, "POST /ISAPI/AccessControl/CardInfo/Search?format=json", []byte(jsonData))
	if err != nil {
		return ResponseData{}, err
	}

	// 解析响应数据
	response, err := ParseResponseData(responseData)
	if err != nil {
		return ResponseData{}, err
	}

	return response, nil
}

// DeleteCard 删除卡片
func (c *ACSClient) DeleteCard(cardNo string) (ResponseData, error) {
	if c.lUserID < 0 {
		return ResponseData{}, fmt.Errorf("未登录设备")
	}

	// 调用内部方法
	err := c.cardManage.DeleteCardInfo(c.lUserID, cardNo)
	if err != nil {
		return ResponseData{}, err
	}

	// 构建删除JSON
	jsonData := fmt.Sprintf(`{
		"CardInfoDelCond": {
			"CardNoList": [
				{
					"cardNo": "%s"
				}
			]
		}
	}`, cardNo)

	// 获取响应数据
	responseData, err := c.cardManage.SendISAPIRequest(c.lUserID, "PUT /ISAPI/AccessControl/CardInfo/Delete?format=json", []byte(jsonData))
	if err != nil {
		return ResponseData{}, err
	}

	// 解析响应数据
	response, err := ParseResponseData(responseData)
	if err != nil {
		return ResponseData{}, err
	}

	return response, nil
}

// DeleteCards 批量删除卡片
func (c *ACSClient) DeleteCards(cardNos []string) ([]ResponseData, []error) {
	if c.lUserID < 0 {
		return nil, []error{fmt.Errorf("未登录设备")}
	}

	errors := make([]error, 0)
	responses := make([]ResponseData, 0)

	for _, cardNo := range cardNos {
		response, err := c.DeleteCard(cardNo)
		responses = append(responses, response)

		if err != nil {
			errors = append(errors, fmt.Errorf("删除卡片 %s 失败: %v", cardNo, err))
		}
	}

	return responses, errors
}

// AddFaceByBinary 通过二进制方式添加人脸
func (c *ACSClient) AddFaceByBinary(employeeNo string) (ResponseData, error) {
	if c.lUserID < 0 {
		return ResponseData{}, fmt.Errorf("未登录设备")
	}

	// 调用内部方法
	err := c.faceManage.AddFaceByBinary(c.lUserID, employeeNo)
	if err != nil {
		return ResponseData{}, err
	}

	// 获取响应数据
	responseData, err := c.faceManage.SendISAPIRequest(c.lUserID, "PUT /ISAPI/Intelligent/FDLib/FDSetUp?format=json", nil)
	if err != nil {
		return ResponseData{}, err
	}

	// 解析响应数据
	response, err := ParseResponseData(responseData)
	if err != nil {
		return ResponseData{}, err
	}

	return response, nil
}

// AddFacesByBinary 通过二进制方式批量添加人脸
func (c *ACSClient) AddFacesByBinary(employeeNos []string) ([]ResponseData, []error) {
	if c.lUserID < 0 {
		return nil, []error{fmt.Errorf("未登录设备")}
	}

	errors := make([]error, 0)
	responses := make([]ResponseData, 0)

	for _, employeeNo := range employeeNos {
		response, err := c.AddFaceByBinary(employeeNo)
		responses = append(responses, response)

		if err != nil {
			errors = append(errors, fmt.Errorf("添加人脸 %s 失败: %v", employeeNo, err))
		}
	}

	return responses, errors
}

// AddFaceByUrl 通过URL方式添加人脸
func (c *ACSClient) AddFaceByUrl(employeeNo string) (ResponseData, error) {
	if c.lUserID < 0 {
		return ResponseData{}, fmt.Errorf("未登录设备")
	}

	// 调用内部方法
	err := c.faceManage.AddFaceByUrl(c.lUserID, employeeNo)
	if err != nil {
		return ResponseData{}, err
	}

	// 获取响应数据
	responseData, err := c.faceManage.SendISAPIRequest(c.lUserID, "PUT /ISAPI/Intelligent/FDLib/FaceDataRecord?format=json", nil)
	if err != nil {
		return ResponseData{}, err
	}

	// 解析响应数据
	response, err := ParseResponseData(responseData)
	if err != nil {
		return ResponseData{}, err
	}

	return response, nil
}

// AddFacesByUrl 通过URL方式批量添加人脸
func (c *ACSClient) AddFacesByUrl(employeeNos []string) ([]ResponseData, []error) {
	if c.lUserID < 0 {
		return nil, []error{fmt.Errorf("未登录设备")}
	}

	errors := make([]error, 0)
	responses := make([]ResponseData, 0)

	for _, employeeNo := range employeeNos {
		response, err := c.AddFaceByUrl(employeeNo)
		responses = append(responses, response)

		if err != nil {
			errors = append(errors, fmt.Errorf("添加人脸 %s 失败: %v", employeeNo, err))
		}
	}

	return responses, errors
}

// SearchFace 查询人脸信息
func (c *ACSClient) SearchFace(employeeNo string) (ResponseData, error) {
	if c.lUserID < 0 {
		return ResponseData{}, fmt.Errorf("未登录设备")
	}

	// 调用内部方法
	err := c.faceManage.SearchFaceInfo(c.lUserID, employeeNo)
	if err != nil {
		return ResponseData{}, err
	}

	// 构建查询JSON
	jsonData := fmt.Sprintf(`{
		"FaceInfoSearchCond": {
			"searchResultPosition": 0,
			"maxResults": 30,
			"faceLibType": "blackFD",
			"FDID": "1",
			"employeeNo": "%s"
		}
	}`, employeeNo)

	// 获取响应数据
	responseData, err := c.faceManage.SendISAPIRequest(c.lUserID, "POST /ISAPI/Intelligent/FDLib/FDSearch?format=json", []byte(jsonData))
	if err != nil {
		return ResponseData{}, err
	}

	// 解析响应数据
	response, err := ParseResponseData(responseData)
	if err != nil {
		return ResponseData{}, err
	}

	return response, nil
}

// DeleteFace 删除人脸
func (c *ACSClient) DeleteFace(employeeNo string) (ResponseData, error) {
	if c.lUserID < 0 {
		return ResponseData{}, fmt.Errorf("未登录设备")
	}

	// 调用内部方法
	err := c.faceManage.DeleteFaceInfo(c.lUserID, employeeNo)
	if err != nil {
		return ResponseData{}, err
	}

	// 构建删除JSON
	jsonData := fmt.Sprintf(`{
		"FaceInfoDelCond": {
			"FDID": "1",
			"faceLibType": "blackFD",
			"employeeNo": "%s"
		}
	}`, employeeNo)

	// 获取响应数据
	responseData, err := c.faceManage.SendISAPIRequest(c.lUserID, "PUT /ISAPI/Intelligent/FDLib/FDDelete?format=json", []byte(jsonData))
	if err != nil {
		return ResponseData{}, err
	}

	// 解析响应数据
	response, err := ParseResponseData(responseData)
	if err != nil {
		return ResponseData{}, err
	}

	return response, nil
}

// DeleteFaces 批量删除人脸
func (c *ACSClient) DeleteFaces(employeeNos []string) ([]ResponseData, []error) {
	if c.lUserID < 0 {
		return nil, []error{fmt.Errorf("未登录设备")}
	}

	errors := make([]error, 0)
	responses := make([]ResponseData, 0)

	for _, employeeNo := range employeeNos {
		response, err := c.DeleteFace(employeeNo)
		responses = append(responses, response)

		if err != nil {
			errors = append(errors, fmt.Errorf("删除人脸 %s 失败: %v", employeeNo, err))
		}
	}

	return responses, errors
}

// CaptureFace 采集人脸
func (c *ACSClient) CaptureFace() (ResponseData, error) {
	if c.lUserID < 0 {
		return ResponseData{}, fmt.Errorf("未登录设备")
	}

	// 调用内部方法
	err := c.faceManage.CaptureFaceInfo(c.lUserID)
	if err != nil {
		return ResponseData{}, err
	}

	// 获取响应数据
	responseData, err := c.faceManage.SendISAPIRequest(c.lUserID, "PUT /ISAPI/Intelligent/FDLib/FDSearch?format=json", nil)
	if err != nil {
		return ResponseData{}, err
	}

	// 解析响应数据
	response, err := ParseResponseData(responseData)
	if err != nil {
		return ResponseData{}, err
	}

	return response, nil
}

// SearchAllEvents 查询所有门禁历史事件
func (c *ACSClient) SearchAllEvents() (ResponseData, error) {
	if c.lUserID < 0 {
		return ResponseData{}, fmt.Errorf("未登录设备")
	}

	// 调用内部方法
	err := c.eventSearch.SearchAllEvent(c.lUserID)
	if err != nil {
		return ResponseData{}, err
	}

	// 获取响应数据
	responseData, err := c.eventSearch.SendISAPIRequest(c.lUserID, "POST /ISAPI/AccessControl/AcsEvent/Search?format=json", nil)
	if err != nil {
		return ResponseData{}, err
	}

	// 解析响应数据
	response, err := ParseResponseData(responseData)
	if err != nil {
		return ResponseData{}, err
	}

	return response, nil
}

// SetCardTemplate 设置卡片模板
func (c *ACSClient) SetCardTemplate(templateNo int) (ResponseData, error) {
	if c.lUserID < 0 {
		return ResponseData{}, fmt.Errorf("未登录设备")
	}

	// 调用内部方法
	err := c.userManage.SetCardTemplate(c.lUserID, templateNo)
	if err != nil {
		return ResponseData{}, err
	}

	// 构建模板JSON
	jsonData := fmt.Sprintf(`{
		"CardTemplate": {
			"templateNo": %d
		}
	}`, templateNo)

	// 获取响应数据
	responseData, err := c.userManage.SendISAPIRequest(c.lUserID, "PUT /ISAPI/AccessControl/CardTemplate/Parameter?format=json", []byte(jsonData))
	if err != nil {
		return ResponseData{}, err
	}

	// 解析响应数据
	response, err := ParseResponseData(responseData)
	if err != nil {
		return ResponseData{}, err
	}

	return response, nil
}
