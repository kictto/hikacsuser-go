package acsapi

import (
	"fmt"
	"time"
)

// 默认时间格式
const timeFormat = "2006-01-01T15:04:05"

// 默认的开始和结束时间
var (
	defaultBeginTime = time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local)
	defaultEndTime   = time.Date(2031, 1, 1, 0, 0, 0, 0, time.Local)
)

// AddUserWithInfo 使用UserInfo结构体添加用户
func (c *ACSClient) AddUserWithInfo(userInfo UserInfo) (ResponseData, error) {
	var response ResponseData
	if c.lUserID < 0 {
		return response, fmt.Errorf("未登录设备")
	}

	// 设置默认值
	if userInfo.Name == "" {
		userInfo.Name = "测试用户"
	}
	if userInfo.UserType == "" {
		userInfo.UserType = "normal"
	}
	if userInfo.BelongGroup == "" {
		userInfo.BelongGroup = "1"
	}
	if userInfo.DoorRight == "" {
		userInfo.DoorRight = "1"
	}
	if userInfo.UserVerifyMode == "" {
		userInfo.UserVerifyMode = "cardOrFace"
	}

	// 设置有效期默认值
	if userInfo.Valid.TimeType == "" {
		userInfo.Valid.TimeType = "local"
	}
	if userInfo.Valid.BeginTime.IsZero() {
		userInfo.Valid.BeginTime = defaultBeginTime
	}
	if userInfo.Valid.EndTime.IsZero() {
		userInfo.Valid.EndTime = defaultEndTime
	}

	// 设置权限计划默认值
	if len(userInfo.RightPlan) == 0 {
		userInfo.RightPlan = []DoorPlan{{
			DoorNo:         1,
			PlanTemplateNo: "1",
		}}
	}

	// 构建用户信息JSON
	rightPlanJSON := ""
	for i, plan := range userInfo.RightPlan {
		if i > 0 {
			rightPlanJSON += ","
		}
		rightPlanJSON += fmt.Sprintf(`{
				"doorNo": %d,
				"planTemplateNo": "%s"
			}`, plan.DoorNo, plan.PlanTemplateNo)
	}

	jsonData := fmt.Sprintf(`{
		"UserInfo": {
			"employeeNo": "%s",
			"name": "%s",
			"userType": "%s",
			"Valid": {
				"enable": %t,
				"beginTime": "%s",
				"endTime": "%s",
				"timeType": "%s"
			},
			"belongGroup": "%s",
			"doorRight": "%s",
			"password": "%s",
			"RightPlan": [
				%s
			],
			"maxOpenDoorTime": %d,
			"openDoorTime": %d,
			"localUIRight": %t,
			"userVerifyMode": "%s"
		}
	}`,
		userInfo.EmployeeNo,
		userInfo.Name,
		userInfo.UserType,
		userInfo.Valid.Enable,
		userInfo.Valid.BeginTime.Format(timeFormat),
		userInfo.Valid.EndTime.Format(timeFormat),
		userInfo.Valid.TimeType,
		userInfo.BelongGroup,
		userInfo.DoorRight,
		userInfo.Password,
		rightPlanJSON,
		userInfo.MaxOpenDoorTime,
		userInfo.OpenDoorTime,
		userInfo.LocalUIRight,
		userInfo.UserVerifyMode)

	respData, err := c.userManage.AddUserInfoWithJSON(c.lUserID, jsonData)
	if err != nil {
		return response, err
	}

	// 解析响应数据
	response, err = ParseResponseData(respData)
	return response, err
}

// DeleteUserImpl 删除单个用户的实现
// 这个方法替代了原来交互式的DeleteUser方法
func (c *ACSClient) DeleteUserImpl(employeeNo string) (ResponseData, error) {
	var response ResponseData
	if c.lUserID < 0 {
		return response, fmt.Errorf("未登录设备")
	}

	// 构建删除用户的JSON请求
	jsonData := fmt.Sprintf(`{
		"UserInfoDelCond": {
			"OperateType": "delete",
			"EmployeeNoList": [
				{
					"employeeNo": "%s"
				}
			]
		}
	}`, employeeNo)

	// URL
	url := "PUT /ISAPI/AccessControl/UserInfo/Delete?format=json"

	// 发送ISAPI请求
	respData, err := c.userManage.SendISAPIRequest(c.lUserID, url, []byte(jsonData))
	if err != nil {
		return response, fmt.Errorf("删除用户失败: %v", err)
	}

	// 解析响应数据
	response, err = ParseResponseData(respData)
	fmt.Printf("删除用户信息成功, 响应: %s\n", string(respData))
	return response, err
}

// DeleteAllUsersImpl 删除所有用户的实现
// 这个方法替代了原来交互式的DeleteAllUsers方法
func (c *ACSClient) DeleteAllUsersImpl() (ResponseData, error) {
	var response ResponseData
	if c.lUserID < 0 {
		return response, fmt.Errorf("未登录设备")
	}

	// 构建删除所有用户的JSON请求
	jsonData := `{
		"UserInfoDelCond": {
			"OperateType": "clear"
		}
	}`

	// URL
	url := "PUT /ISAPI/AccessControl/UserInfo/Delete?format=json"

	// 发送ISAPI请求
	respData, err := c.userManage.SendISAPIRequest(c.lUserID, url, []byte(jsonData))
	if err != nil {
		return response, fmt.Errorf("删除所有用户失败: %v", err)
	}

	// 解析响应数据
	response, err = ParseResponseData(respData)
	fmt.Printf("删除所有用户信息成功, 响应: %s\n", string(respData))
	return response, err
}

// DeleteUser 删除单个用户
// 这个方法覆盖了api.go中的同名方法
func (c *ACSClient) DeleteUser(employeeNo string) (ResponseData, error) {
	return c.DeleteUserImpl(employeeNo)
}

// DeleteAllUsers 删除所有用户
// 这个方法覆盖了api.go中的同名方法
func (c *ACSClient) DeleteAllUsers() (ResponseData, error) {
	return c.DeleteAllUsersImpl()
}

// DeleteUsers 批量删除用户
func (c *ACSClient) DeleteUsers(employeeNos []string) ([]ResponseData, []error) {
	if c.lUserID < 0 {
		return nil, []error{fmt.Errorf("未登录设备")}
	}

	responses := make([]ResponseData, 0)
	errors := make([]error, 0)
	for _, employeeNo := range employeeNos {
		resp, err := c.DeleteUserImpl(employeeNo)
		responses = append(responses, resp)
		if err != nil {
			errors = append(errors, fmt.Errorf("删除用户 %s 失败: %v", employeeNo, err))
		}
	}

	return responses, errors
}

// SetupUser 设置用户信息（支持新增和修改）
// 使用PUT /ISAPI/AccessControl/UserInfo/SetUp?format=json接口
func (c *ACSClient) SetupUser(userInfo UserInfo) (ResponseData, error) {
	var response ResponseData
	if c.lUserID < 0 {
		return response, fmt.Errorf("未登录设备")
	}

	// 设置默认值
	if userInfo.Name == "" {
		userInfo.Name = "测试用户"
	}
	if userInfo.UserType == "" {
		userInfo.UserType = "normal"
	}
	if userInfo.BelongGroup == "" {
		userInfo.BelongGroup = "1"
	}
	if userInfo.DoorRight == "" {
		userInfo.DoorRight = "1"
	}
	if userInfo.UserVerifyMode == "" {
		userInfo.UserVerifyMode = "cardOrFace"
	}

	// 设置有效期默认值
	if userInfo.Valid.TimeType == "" {
		userInfo.Valid.TimeType = "local"
	}
	if userInfo.Valid.BeginTime.IsZero() {
		userInfo.Valid.BeginTime = defaultBeginTime
	}
	if userInfo.Valid.EndTime.IsZero() {
		userInfo.Valid.EndTime = defaultEndTime
	}

	// 设置权限计划默认值
	if len(userInfo.RightPlan) == 0 {
		userInfo.RightPlan = []DoorPlan{{
			DoorNo:         1,
			PlanTemplateNo: "1",
		}}
	}

	// 构建用户信息JSON
	rightPlanJSON := ""
	for i, plan := range userInfo.RightPlan {
		if i > 0 {
			rightPlanJSON += ","
		}
		rightPlanJSON += fmt.Sprintf(`{
				"doorNo": %d,
				"planTemplateNo": "%s"
			}`, plan.DoorNo, plan.PlanTemplateNo)
	}

	jsonData := fmt.Sprintf(`{
		"UserInfo": {
			"employeeNo": "%s",
			"name": "%s",
			"userType": "%s",
			"Valid": {
				"enable": %t,
				"beginTime": "%s",
				"endTime": "%s",
				"timeType": "%s"
			},
			"belongGroup": "%s",
			"doorRight": "%s",
			"password": "%s",
			"RightPlan": [
				%s
			],
			"maxOpenDoorTime": %d,
			"openDoorTime": %d,
			"localUIRight": %t,
			"userVerifyMode": "%s"
		}
	}`,
		userInfo.EmployeeNo,
		userInfo.Name,
		userInfo.UserType,
		userInfo.Valid.Enable,
		userInfo.Valid.BeginTime.Format(timeFormat),
		userInfo.Valid.EndTime.Format(timeFormat),
		userInfo.Valid.TimeType,
		userInfo.BelongGroup,
		userInfo.DoorRight,
		userInfo.Password,
		rightPlanJSON,
		userInfo.MaxOpenDoorTime,
		userInfo.OpenDoorTime,
		userInfo.LocalUIRight,
		userInfo.UserVerifyMode)

	// URL - 使用SetUp接口，支持新增和修改
	url := "PUT /ISAPI/AccessControl/UserInfo/SetUp?format=json"

	// 发送ISAPI请求
	respData, err := c.userManage.SendISAPIRequest(c.lUserID, url, []byte(jsonData))
	if err != nil {
		return response, fmt.Errorf("设置用户信息失败: %v", err)
	}

	// 解析响应数据
	response, err = ParseResponseData(respData)
	fmt.Printf("设置用户信息成功, 响应: %s\n", string(respData))
	return response, err
}

// SetupUsers 批量设置用户信息（支持新增和修改）
func (c *ACSClient) SetupUsers(userInfos []UserInfo) ([]ResponseData, []error) {
	if c.lUserID < 0 {
		return nil, []error{fmt.Errorf("未登录设备")}
	}

	responses := make([]ResponseData, 0)
	errors := make([]error, 0)
	for _, userInfo := range userInfos {
		resp, err := c.SetupUser(userInfo)
		responses = append(responses, resp)
		if err != nil {
			errors = append(errors, fmt.Errorf("设置用户 %s 信息失败: %v", userInfo.EmployeeNo, err))
		}
	}

	return responses, errors
}
