package config

import (
	"github.com/Deansquirrel/goToolCommon"
	"strings"
)

type dingTalkConfig struct {
	Address string `tom;:"address"`
}

func (dt *dingTalkConfig) FormatConfig() {
	dt.Address = strings.Trim(dt.Address, " ")
	dt.Address = strings.ToLower(dt.Address)
	dt.Address = goToolCommon.CheckAndDeleteLastChar(dt.Address, "/")
	dt.Address = goToolCommon.CheckAndDeleteLastChar(dt.Address, "\\")
}
