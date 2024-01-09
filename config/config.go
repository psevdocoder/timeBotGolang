package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Whitelist     []int64 `json:"whitelist"`
	AdminID       int64
	TelegramToken string `json:"telegramToken"`
	CityURL       string `json:"cityURL"`
	UpdateTime    string `json:"updateTime"`
	TimeTill      int    `json:"timeTill"`
}

func LoadConfig() (*Config, error) {
	file, err := os.ReadFile("config/config.json")
	var config Config
	if err = json.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Config) SetCityURL(url string) {
	c.CityURL = url
	updateConfig(c)
}

func (c *Config) SetUpdateTime(updateTime string) {
	c.UpdateTime = updateTime
	updateConfig(c)
}

func (c *Config) SetTimeTill(timeTill int) {
	c.TimeTill = timeTill
	updateConfig(c)
}

func (c *Config) EditWhitelist(whitelist []int64) {
	c.Whitelist = append(c.Whitelist, whitelist...)
	updateConfig(c)
}

func updateConfig(config *Config) {
	file, _ := json.MarshalIndent(config, "", "  ")
	_ = os.WriteFile("config/config.json", file, 0644)
}
