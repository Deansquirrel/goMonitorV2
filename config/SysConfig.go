package config

import (
	"github.com/Deansquirrel/goToolCommon"
)

type SysConfig struct {
	Total          Total          `toml:"total"`
	ServiceConfig  serviceConfig  `toml:"serviceConfig"`
	DingTalkConfig dingTalkConfig `toml:"dingTalkConfig"`
	ConfigDBConfig configDBConfig `toml:"configDBConfig"`
	TaskConfig     taskConfig     `toml:"taskConfig"`
	IrisConfig     irisConfig     `toml:"irisConfig"`
}

//返回配置字符串
func (sc *SysConfig) GetConfigStr() (string, error) {
	return goToolCommon.GetJsonStr(sc)
}

//配置检查并格式化
func (sc *SysConfig) FormatConfig() {
	sc.Total.FormatConfig()
	sc.ServiceConfig.FormatConfig()
	sc.DingTalkConfig.FormatConfig()
	sc.TaskConfig.FormatConfig()
	sc.IrisConfig.FormatConfig()
}
