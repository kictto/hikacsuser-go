package acsapi

import (
	"fmt"
)

// AddCardWithInfo 使用CardInfo结构体添加卡片
func (c *ACSClient) AddCardWithInfo(cardInfo CardInfo) error {
	if c.lUserID < 0 {
		return fmt.Errorf("未登录设备")
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

	return c.cardManage.AddCardInfoWithJSON(c.lUserID, jsonData)
}

// AddCardsWithInfo 批量添加卡片，使用CardInfo结构体数组
func (c *ACSClient) AddCardsWithInfo(cardInfos []CardInfo) []error {
	if c.lUserID < 0 {
		return []error{fmt.Errorf("未登录设备")}
	}

	errors := make([]error, 0)
	for _, cardInfo := range cardInfos {
		err := c.AddCardWithInfo(cardInfo)
		if err != nil {
			errors = append(errors, fmt.Errorf("添加卡片 %s 失败: %v", cardInfo.CardNo, err))
		}
	}

	return errors
}
