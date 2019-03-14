package router

import (
	"github.com/Deansquirrel/goMonitorV2/global"
	"github.com/kataras/iris"
)

type base struct {
	app *iris.Application
	c   common
}

//版本信息返回数据结构（version）
type versionInfo struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Version string `json:"version"`
}

func NewRouterBase(app *iris.Application) *base {
	return &base{
		app: app,
		c:   common{},
	}
}

func (base *base) AddBase() {
	base.app.Get("/version", base.version)
}

func (base *base) version(ctx iris.Context) {
	v := versionInfo{
		ErrCode: 0,
		ErrMsg:  "",
		Version: global.Version,
	}
	base.c.WriteResponse(ctx, v)
}
