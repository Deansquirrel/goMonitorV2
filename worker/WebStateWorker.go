package worker

import (
	"github.com/Deansquirrel/goMonitorV2/global"
	"github.com/Deansquirrel/goMonitorV2/taskConfigRepository"
	"github.com/Deansquirrel/goToolCommon"
	log "github.com/Deansquirrel/goToolLog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type webStateWorker struct {
	webStateTaskConfigData *taskConfigRepository.WebStateTaskConfigData
}

func NewWebStateWorker(webStateTaskConfigData *taskConfigRepository.WebStateTaskConfigData) *webStateWorker {
	return &webStateWorker{
		webStateTaskConfigData: webStateTaskConfigData,
	}
}

func (wsw *webStateWorker) Run() {
	begTime := time.Now()
	log.Debug(goToolCommon.GetDateTimeStr(begTime))
	code, err := wsw.getHttpCode()
	endTime := time.Now()
	log.Debug(goToolCommon.GetDateTimeStr(endTime))
	if err != nil {
		log.Error(err.Error())
		comm.sendMsg(wsw.webStateTaskConfigData.FId, wsw.getMsg(wsw.webStateTaskConfigData.FMsgTitle, err.Error()))
		return
	}
	ns := endTime.Sub(begTime).Nanoseconds()
	ms := ns / 1000 / 1000
	if code == 200 {
		wsw.saveHisResult(200, int(ms), "")
		return
	} else {
		msg := wsw.getMsg(wsw.webStateTaskConfigData.FMsgTitle, wsw.getMsgContent(code))
		comm.sendMsg(wsw.webStateTaskConfigData.FId, msg)
		wsw.saveHisResult(200, int(ms), msg)
		return
	}
}

func (wsw *webStateWorker) getMsg(title, content string) string {
	msg := comm.getMsg(title, content)
	if msg != "" {
		msg = msg + "\n"
	}
	msg = msg + wsw.webStateTaskConfigData.FUrl
	return msg
}

func (wsw *webStateWorker) getMsgContent(code int) string {
	content := wsw.webStateTaskConfigData.FMsgContent
	content = strings.Replace(content, "title", strconv.Itoa(code), -1)
	return content
}

func (wsw *webStateWorker) getHttpCode() (int, error) {
	req, err := http.NewRequest("GET", wsw.webStateTaskConfigData.FUrl, nil)
	if err != nil {
		return -1, err
	}
	client := &http.Client{
		Timeout: time.Second * global.HttpGetTimeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return resp.StatusCode, nil
	//u,err := url.Parse(wsw.webStateTaskConfigData.FUrl)
	//if err != nil {
	//	return 0,err
	//}
	//q := u.Query()
	//u.RawQuery = q.Encode()
	//res,err := http.Get(u.String())
	//if err != nil {
	//	return 0,err
	//}
	//resCode := res.StatusCode
	//_ = res.Body.Close()
	//return resCode,nil
}

func (wsw *webStateWorker) saveHisResult(httpCode int, ms int, content string) {
	nData := &taskConfigRepository.WebStateTaskHisData{
		FId:       strings.ToUpper(goToolCommon.Guid()),
		FConfigId: wsw.webStateTaskConfigData.FId,
		FUseTime:  ms,
		FHttpCode: httpCode,
		FContent:  content,
	}
	webStateHis := taskConfigRepository.WebStateTaskHis{}
	_ = webStateHis.SetWebStateTaskHis(nData)
}
