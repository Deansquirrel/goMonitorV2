package worker

import (
	"github.com/Deansquirrel/goMonitorV2/taskConfigRepository"
	"github.com/Deansquirrel/goToolEnvironment"
	log "github.com/Deansquirrel/goToolLog"
)

type healthWorker struct {
	healthTaskConfigData *taskConfigRepository.HealthTaskConfigData
}

func NewHealthWorker(healthTaskConfigData *taskConfigRepository.HealthTaskConfigData) *healthWorker {
	return &healthWorker{
		healthTaskConfigData: healthTaskConfigData,
	}
}

//检查
func (hw *healthWorker) Run() {
	if hw.healthTaskConfigData == nil {
		return
	}
	comm.sendMsg(hw.healthTaskConfigData.FId, hw.getMsg())
}

func (hw *healthWorker) getMsg() string {
	iAddr, err := goToolEnvironment.GetInternetAddr()
	if err != nil {
		log.Error("获取外网地址时遇到错误：" + err.Error())
	} else {
		if hw.healthTaskConfigData.FMsgTitle != "" {
			iAddr = iAddr + "\n" + hw.healthTaskConfigData.FMsgTitle
		}
	}
	msg := comm.getMsg(iAddr, hw.healthTaskConfigData.FMsgContent)
	return msg
}
