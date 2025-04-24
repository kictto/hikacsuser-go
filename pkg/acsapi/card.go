package acsapi

import (
	"fmt"
)

// AddCardWithInfo 使用CardInfo结构体添加卡片
func (c *ACSClient) AddCardWithInfo(cardInfo CardInfo) (ResponseData, error) {
	var response ResponseData
	if c.lUserID < 0 {
		return response, fmt.Errorf("未登录设备")
	}

	// 设置默认值
	if cardInfo.CardType == "" {
		cardInfo.CardType = "normalCard"
	}

	// 构建卡片信息JSON
	jsonData := fmt.Sprintf(`{
		"CardInfo": {
			"employeeNo": "%s",
			"cardNo": "%s",
			"cardType": "%s"
		}
	}`, cardInfo.EmployeeNo, cardInfo.CardNo, cardInfo.CardType)

	respData, err := c.cardManage.AddCardInfoWithJSON(c.lUserID, jsonData)
	if err != nil {
		return response, err
	}

	// 解析响应数据
	response, err = ParseResponseData(respData)
	return response, err
}

// AddCardsWithInfo 批量添加卡片，使用CardInfo结构体数组
func (c *ACSClient) AddCardsWithInfo(cardInfos []CardInfo) ([]ResponseData, []error) {
	if c.lUserID < 0 {
		return nil, []error{fmt.Errorf("未登录设备")}
	}

	responses := make([]ResponseData, 0)
	errors := make([]error, 0)
	for _, cardInfo := range cardInfos {
		resp, err := c.AddCardWithInfo(cardInfo)
		responses = append(responses, resp)
		if err != nil {
			errors = append(errors, fmt.Errorf("添加卡片 %s 失败: %v", cardInfo.CardNo, err))
		}
	}

	return responses, errors
}
