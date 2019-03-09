package global

import (
	"context"
	"github.com/Deansquirrel/goMonitorV2/config"
	"github.com/Deansquirrel/goToolMSSql"
	"time"
)

const (
	//PreVersion = "1.0.0 Build20190305"
	//TestVersion = "0.0.0 Build20190101"
	Version = "0.0.0 Build20190101"
)

var SysConfig *config.SysConfig
var Ctx context.Context
var Cancel func()

func init() {
	goToolMSSql.SetMaxIdleConn(3)
	goToolMSSql.SetMaxOpenConn(3)
	goToolMSSql.SetMaxLifetime(time.Second * 15)
}
