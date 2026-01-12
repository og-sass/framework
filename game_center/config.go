package game_center

import "strings"

type CenterConfig struct {
	Username     string                  `json:"username"`
	Password     string                  `json:"password"`
	RequestURL   string                  `json:"request_url"`
	CurrencyConf []CurrencyItem          `json:"currency_conf"`
	currencyMap  map[string]CurrencyItem // 新增：用于快速查找的映射
}

type CurrencyItem struct {
	Currency string `json:"currency"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Init 初始化配置，构建货币映射
func (c *CenterConfig) Init() {
	c.currencyMap = make(map[string]CurrencyItem)
	for _, item := range c.CurrencyConf {
		c.currencyMap[strings.ToLower(item.Currency)] = item
	}
}

// GetCurrencyConf 获取指定货币的配置（优化版）
func (c *CenterConfig) GetCurrencyConf(currency string) CurrencyItem {
	if item, ok := c.currencyMap[strings.ToLower(currency)]; ok {
		return item
	}

	return CurrencyItem{
		Currency: currency,
		Username: c.Username,
		Password: c.Password,
	}
}
