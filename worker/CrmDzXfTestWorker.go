package worker

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Deansquirrel/goMonitorV2/global"
	"github.com/Deansquirrel/goMonitorV2/taskConfigRepository"
	"github.com/Deansquirrel/goToolCommon"
	log "github.com/Deansquirrel/goToolLog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type crmDzXfTestWorker struct {
	crmDzTestTaskConfigData *taskConfigRepository.CrmDzXfTestTaskConfigData
}

func NewCrmDzXfTestWorker(crmDzTestTaskConfigData *taskConfigRepository.CrmDzXfTestTaskConfigData) *crmDzXfTestWorker {
	return &crmDzXfTestWorker{
		crmDzTestTaskConfigData: crmDzTestTaskConfigData,
	}
}

func (cw *crmDzXfTestWorker) Run() {
	var msg string
	begTime := time.Now()
	code, err := cw.getHttpCode()
	endTime := time.Now()
	ns := endTime.Sub(begTime).Nanoseconds()
	ms := ns / 1000 / 1000
	if err != nil {
		log.Error(err.Error())
		msg = "Crm定制消费接口测试遇到错误" + "\n" + err.Error()
	} else if code != 200 || ms > 5*1000 {
		msg = cw.getMsgContent(code, int(ms))
		if msg == "" {
			msg = fmt.Sprintf("返回码：%d\n用时：%d", code, ms)
		}
	} else {
		msg = ""
	}
	if msg != "" {
		msg = cw.getMsg(cw.crmDzTestTaskConfigData.FMsgTitle, msg)
		comm.sendMsg(cw.crmDzTestTaskConfigData.FId, msg)
	}
	cw.saveHisResult(code, int(ms), msg)
	return
}

func (cw *crmDzXfTestWorker) getMsg(title, content string) string {
	msg := comm.getMsg(title, content)
	if msg != "" {
		msg = msg + "\n"
	}
	msg = msg + cw.crmDzTestTaskConfigData.FAddress
	return msg
}

func (cw *crmDzXfTestWorker) getMsgContent(code, ms int) string {
	content := cw.crmDzTestTaskConfigData.FMsgContent
	content = strings.Replace(content, "code", strconv.Itoa(code), -1)
	content = strings.Replace(content, "ms", strconv.Itoa(ms), -1)
	return content
}

func (cw *crmDzXfTestWorker) getHttpCode() (int, error) {
	//TODO
	testData, err := goToolCommon.GetJsonStr(cw.getTestData())
	if err != nil {
		return 0, errors.New("构造测试数据时发生错误：" + err.Error())
	}
	req, err := http.NewRequest("POST", cw.crmDzTestTaskConfigData.FAddress, bytes.NewBuffer([]byte(testData)))
	if err != nil {
		return 0, errors.New("构造http请求数据时发生错误：" + err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("passporttype", strconv.Itoa(cw.crmDzTestTaskConfigData.FPassportType))
	req.Header.Set("passport", cw.crmDzTestTaskConfigData.FPassport)

	client := &http.Client{
		Timeout: time.Second * global.HttpTimeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, errors.New("发送http请求时错误：" + err.Error())
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return resp.StatusCode, nil
}

func (cw *crmDzXfTestWorker) saveHisResult(code, ms int, content string) {
	nData := &taskConfigRepository.CrmDzXfTestHisData{
		FId:       strings.ToUpper(goToolCommon.Guid()),
		FConfigId: cw.crmDzTestTaskConfigData.FId,
		FUseTime:  ms,
		FHttpCode: code,
		FContent:  content,
	}
	crmDzXfHis := taskConfigRepository.CrmDzXfTestHis{}
	_ = crmDzXfHis.SetCrmDzXfTestTaskHis(nData)
}

func (cw *crmDzXfTestWorker) getTestData() *crmDzXfRequestData {
	ywCore := crmDzXfRequestDataYwCore{
		Oprtime:   goToolCommon.GetDateTimeStr(time.Now()),
		Oprbrid:   10001,
		Oprbrname: "测试请求",
		Oprxfje:   1000000,
	}
	ywInfo := crmDzXfRequestDataYwInfo{
		Oprywsno:    "YW" + goToolCommon.GetDateStr(time.Now()) + "01",
		Oprppid:     10001,
		Oprppname:   "",
		Oprid:       182,
		Oprname:     "管理员",
		Oprywdate:   goToolCommon.GetDateStr(time.Now()) + " 00:00:00",
		Oprskbrid:   0,
		Oprskbrname: "",
		Oprskppid:   0,
		Oprskppname: "",
		Oprywwindow: "测试请求",
		Oprywbno:    "",
		Oprsummary:  "",
	}
	return &crmDzXfRequestData{
		YwCore: ywCore,
		YwInfo: ywInfo,
	}
}

type crmDzXfRequestData struct {
	YwCore crmDzXfRequestDataYwCore
	YwInfo crmDzXfRequestDataYwInfo
}

type crmDzXfRequestDataYwCore struct {
	Oprtime   string
	Oprbrid   int
	Oprbrname string
	Oprxfje   int
}

type crmDzXfRequestDataYwInfo struct {
	Oprywsno    string
	Oprppid     int
	Oprppname   string
	Oprid       int
	Oprname     string
	Oprywdate   string
	Oprskbrid   int
	Oprskbrname string
	Oprskppid   int
	Oprskppname string
	Oprywwindow string
	Oprywbno    string
	Oprsummary  string
}
