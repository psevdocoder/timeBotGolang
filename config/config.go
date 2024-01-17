package config

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
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

func (c *Config) String() string {
	return fmt.Sprintf("*Current bot configuration*\nWhitelist: "+
		"%v\nCityURL: `%v`\nUpdate at: `%v`\nSend mins before: %v\n",
		c.Whitelist, c.CityURL, c.UpdateTime, c.TimeTill)
}

func (c *Config) SetCityURL(url string) {
	c.CityURL = url
	c.updateConfig()
}

func (c *Config) SetUpdateTime(updateTime string) {
	c.UpdateTime = updateTime
	c.updateConfig()
}

func (c *Config) SetTimeTill(timeTill int) {
	c.TimeTill = timeTill
	c.updateConfig()
}

func (c *Config) EditWhitelist(whitelist []int64) {
	switch {
	case slices.Contains(whitelist, 0) && len(whitelist) == 1:
		c.Whitelist = []int64{c.AdminID}
	case !slices.Contains(whitelist, 0):
		whitelist = append(whitelist, c.AdminID)
		c.Whitelist = whitelist
	}
	c.updateConfig()
}

func (c *Config) updateConfig() {
	file, _ := json.MarshalIndent(c, "", "  ")
	_ = os.WriteFile("config/config.json", file, 0644)
}
