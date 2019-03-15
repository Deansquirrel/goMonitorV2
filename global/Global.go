package global

import (
	"context"
	"github.com/Deansquirrel/goMonitorV2/config"
	"github.com/Deansquirrel/goToolMSSql"
	"time"
)

const (
	//PreVersion = "0.0.3 Build20190315"
	//TestVersion = "0.0.0 Build20190101"
	Version = "0.0.0 Build20190101"
)

var SysConfig *config.SysConfig
var Ctx context.Context
var Cancel func()

func init() {
	goToolMSSql.SetMaxIdleConn(15)
	goToolMSSql.SetMaxOpenConn(15)
	goToolMSSql.SetMaxLifetime(time.Second * 60)
}
