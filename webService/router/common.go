package router

import (
	"encoding/json"
	"fmt"
	"github.com/Deansquirrel/goDingtalkRobot/object"
	log "github.com/Deansquirrel/goToolLog"
	"github.com/kataras/iris"
	"io/ioutil"
)

const (
	TranErrStr = "{\"errcode\":-1,\"errmsg\":\"构造返回结果时发生错误 [%s]\"}"
)

type common struct {
}

func (c *common) GetRequestBody(ctx iris.Context) string {
	body := ctx.Request().Body
	defer func() {
		_ = body.Close()
	}()
	b, err := ioutil.ReadAll(body)
	if err != nil {
		log.Error("获取Http请求文本时发生错误：" + err.Error())
		return ""
	}
	return string(b)
}

func (c *common) GetOKReturn(msg string) *object.SimpleResponse {
	return c.GetMsgReturn("OK")
}

func (c *common) GetMsgReturn(msg string) *object.SimpleResponse {
	return &object.SimpleResponse{
		ErrCode: 0,
		ErrMsg:  msg,
	}
}

func (c *common) GetErrReturn(err error) *object.SimpleResponse {
	return &object.SimpleResponse{
		ErrCode: -1,
		ErrMsg:  err.Error(),
	}
}

//func (c *common) getReturn(code int, msg string) string {
//	rd := object.SimpleResponse{
//		ErrCode: code,
//		ErrMsg:  msg,
//	}
//	rb, err := json.Marshal(rd)
//	if err != nil {
//		return fmt.Sprintf(TranErrStr, "err:"+err.Error()+",code:"+strconv.Itoa(code)+",msg:"+msg)
//	} else {
//		return string(rb)
//	}
//}

//向ctx中添加返回内容
func (c *common) WriteResponse(ctx iris.Context, v interface{}) {
	str, err := json.Marshal(v)
	if err != nil {
		_, _ = ctx.WriteString(fmt.Sprintf(TranErrStr, "err:"+err.Error()))
		return
	}
	_, _ = ctx.WriteString(string(str))
	return
}
