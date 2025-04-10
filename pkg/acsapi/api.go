// Package acsapi 提供海康威视门禁系统的公共API接口
// 该包将内部实现封装为公共方法，支持批量和单个数据操作
package acsapi

import (
	"fmt"
	"github.com/clockworkchen/hikacsuser-go/internal/models"
	"github.com/clockworkchen/hikacsuser-go/internal/sdk"
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
func NewACSClient() (*ACSClient, error) {
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
	if !client.hcnetsdk.NET_DVR_SetLogToFile(3, "./sdklog", false) {
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
func (c *ACSClient) GetACSConfig() error {
	if c.lUserID < 0 {
		return fmt.Errorf("未登录设备")
	}
	return c.acsManage.AcsCfg(c.lUserID)
}

// GetACSStatus 获取门禁状态
func (c *ACSClient) GetACSStatus() error {
	if c.lUserID < 0 {
		return fmt.Errorf("未登录设备")
	}
	return c.acsManage.GetAcsStatus(c.lUserID)
}

// RemoteControlGate 远程控门
func (c *ACSClient) RemoteControlGate() error {
	if c.lUserID < 0 {
		return fmt.Errorf("未登录设备")
	}
	return c.acsManage.RemoteControlGate(c.lUserID)
}

// AddUser 添加单个用户
func (c *ACSClient) AddUser(employeeNo string) error {
	if c.lUserID < 0 {
		return fmt.Errorf("未登录设备")
	}
	return c.userManage.AddUserInfo(c.lUserID, employeeNo)
}

// AddUsers 批量添加用户
func (c *ACSClient) AddUsers(employeeNos []string) []error {
	if c.lUserID < 0 {
		return []error{fmt.Errorf("未登录设备")}
	}

	errors := make([]error, 0)
	for _, employeeNo := range employeeNos {
		err := c.userManage.AddUserInfo(c.lUserID, employeeNo)
		if err != nil {
			errors = append(errors, fmt.Errorf("添加用户 %s 失败: %v", employeeNo, err))
		}
	}
	return errors
}

// SearchUser 查询用户信息
func (c *ACSClient) SearchUser() error {
	if c.lUserID < 0 {
		return fmt.Errorf("未登录设备")
	}
	return c.userManage.SearchUserInfo(c.lUserID)
}

// DeleteUser 方法已移至user.go

// DeleteAllUsers 方法已移至user.go

// AddCard 添加单个卡片
func (c *ACSClient) AddCard(employeeNo, cardNo string) error {
	if c.lUserID < 0 {
		return fmt.Errorf("未登录设备")
	}
	return c.cardManage.AddCardInfo(c.lUserID, employeeNo, cardNo)
}

// AddCards 批量添加卡片
func (c *ACSClient) AddCards(employeeNos, cardNos []string) []error {
	if c.lUserID < 0 {
		return []error{fmt.Errorf("未登录设备")}
	}

	if len(employeeNos) != len(cardNos) {
		return []error{fmt.Errorf("员工号和卡号数量不匹配")}
	}

	errors := make([]error, 0)
	for i, employeeNo := range employeeNos {
		err := c.cardManage.AddCardInfo(c.lUserID, employeeNo, cardNos[i])
		if err != nil {
			errors = append(errors, fmt.Errorf("添加卡片 %s 失败: %v", cardNos[i], err))
		}
	}

	return errors
}

// SearchCard 查询卡片信息
func (c *ACSClient) SearchCard(employeeNo string) error {
	if c.lUserID < 0 {
		return fmt.Errorf("未登录设备")
	}
	return c.cardManage.SearchCardInfo(c.lUserID, employeeNo)
}

// DeleteCard 删除卡片
func (c *ACSClient) DeleteCard(cardNo string) error {
	if c.lUserID < 0 {
		return fmt.Errorf("未登录设备")
	}
	return c.cardManage.DeleteCardInfo(c.lUserID, cardNo)
}

// DeleteCards 批量删除卡片
func (c *ACSClient) DeleteCards(cardNos []string) []error {
	if c.lUserID < 0 {
		return []error{fmt.Errorf("未登录设备")}
	}

	errors := make([]error, 0)
	for _, cardNo := range cardNos {
		err := c.cardManage.DeleteCardInfo(c.lUserID, cardNo)
		if err != nil {
			errors = append(errors, fmt.Errorf("删除卡片 %s 失败: %v", cardNo, err))
		}
	}

	return errors
}

// AddFaceByBinary 通过二进制方式添加人脸
func (c *ACSClient) AddFaceByBinary(employeeNo string) error {
	if c.lUserID < 0 {
		return fmt.Errorf("未登录设备")
	}
	return c.faceManage.AddFaceByBinary(c.lUserID, employeeNo)
}

// AddFacesByBinary 通过二进制方式批量添加人脸
func (c *ACSClient) AddFacesByBinary(employeeNos []string) []error {
	if c.lUserID < 0 {
		return []error{fmt.Errorf("未登录设备")}
	}

	errors := make([]error, 0)
	for _, employeeNo := range employeeNos {
		err := c.faceManage.AddFaceByBinary(c.lUserID, employeeNo)
		if err != nil {
			errors = append(errors, fmt.Errorf("添加人脸 %s 失败: %v", employeeNo, err))
		}
	}

	return errors
}

// AddFaceByUrl 通过URL方式添加人脸
func (c *ACSClient) AddFaceByUrl(employeeNo string) error {
	if c.lUserID < 0 {
		return fmt.Errorf("未登录设备")
	}
	return c.faceManage.AddFaceByUrl(c.lUserID, employeeNo)
}

// AddFacesByUrl 通过URL方式批量添加人脸
func (c *ACSClient) AddFacesByUrl(employeeNos []string) []error {
	if c.lUserID < 0 {
		return []error{fmt.Errorf("未登录设备")}
	}

	errors := make([]error, 0)
	for _, employeeNo := range employeeNos {
		err := c.faceManage.AddFaceByUrl(c.lUserID, employeeNo)
		if err != nil {
			errors = append(errors, fmt.Errorf("添加人脸 %s 失败: %v", employeeNo, err))
		}
	}

	return errors
}

// SearchFace 查询人脸信息
func (c *ACSClient) SearchFace(employeeNo string) error {
	if c.lUserID < 0 {
		return fmt.Errorf("未登录设备")
	}
	return c.faceManage.SearchFaceInfo(c.lUserID, employeeNo)
}

// DeleteFace 删除人脸
func (c *ACSClient) DeleteFace(employeeNo string) error {
	if c.lUserID < 0 {
		return fmt.Errorf("未登录设备")
	}
	return c.faceManage.DeleteFaceInfo(c.lUserID, employeeNo)
}

// DeleteFaces 批量删除人脸
func (c *ACSClient) DeleteFaces(employeeNos []string) []error {
	if c.lUserID < 0 {
		return []error{fmt.Errorf("未登录设备")}
	}

	errors := make([]error, 0)
	for _, employeeNo := range employeeNos {
		err := c.faceManage.DeleteFaceInfo(c.lUserID, employeeNo)
		if err != nil {
			errors = append(errors, fmt.Errorf("删除人脸 %s 失败: %v", employeeNo, err))
		}
	}

	return errors
}

// CaptureFace 采集人脸
func (c *ACSClient) CaptureFace() error {
	if c.lUserID < 0 {
		return fmt.Errorf("未登录设备")
	}
	return c.faceManage.CaptureFaceInfo(c.lUserID)
}

// SearchAllEvents 查询所有门禁历史事件
func (c *ACSClient) SearchAllEvents() error {
	if c.lUserID < 0 {
		return fmt.Errorf("未登录设备")
	}
	return c.eventSearch.SearchAllEvent(c.lUserID)
}

// SetCardTemplate 设置卡片模板
func (c *ACSClient) SetCardTemplate(templateNo int) error {
	if c.lUserID < 0 {
		return fmt.Errorf("未登录设备")
	}
	return c.userManage.SetCardTemplate(c.lUserID, templateNo)
}
