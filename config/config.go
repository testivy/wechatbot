package config

import (
	"encoding/json"
	"log"
	"os"
)

type ExtraSt struct {
	// 专属群组昵称
	GroupName string `json:"group_name"`
	// 纪念日信息的文件
	MemoDayFile string `json:"memo_day_file"`
}

// Configuration 项目配置
type Configuration struct {
	// gpt apikey
	ApiKey string `json:"api_key"`
	// 是否自动通过好友
	AutoPass bool `json:"auto_pass"`
	// 是否开启gpt聊天功能
	GptChat bool `json:"gpt_chat"`
	// 额外信息
	Extra ExtraSt `json:"extra"`
	//请求API地址
	EndPoint string `json:"endpoint"`
}

var Config *Configuration =&Configuration{}

// LoadConfig 加载配置
func LoadConfig(cfg string) {
	if cfg == "" {
		cfg = "config.json"
		log.Println("cfg: ", cfg)
	}
	b, err := os.ReadFile(cfg)
	if err != nil {
		log.Fatalln("ReadFile error")
	}
	var c Configuration
	err = json.Unmarshal(b, &c)
	if err != nil {
		log.Fatalln("read config file:", cfg, " fail: ", err)
	}
	Config = &c
}
