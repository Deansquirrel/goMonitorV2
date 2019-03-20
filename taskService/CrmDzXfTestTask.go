package taskService

import (
	"fmt"
	"github.com/Deansquirrel/goMonitorV2/taskConfigRepository"
	"github.com/Deansquirrel/goToolCommon"
)
import log "github.com/Deansquirrel/goToolLog"

type CrmDzXfTestTask struct {
}

func (ct *CrmDzXfTestTask) StartTask() error {
	//TODO
	return nil
}

func (ct *CrmDzXfTestTask) StartJob(id string) error {
	//TODO
	return nil
}

func (ct *CrmDzXfTestTask) StopJob(id string) error {
	//TODO
	return nil
}

func (ct *CrmDzXfTestTask) startRegularRefresh() {
	//TODO
	return
}

func (ct *CrmDzXfTestTask) refreshConfig() {
	//TODO
	return
}

func (ct *CrmDzXfTestTask) delHisData() {
	//TODO
	return
}

func (ct *CrmDzXfTestTask) RefreshConfig() error {
	//TODO
	return nil
}

func (ct *CrmDzXfTestTask) addTask(config *taskConfigRepository.CrmDzXfTestTaskConfigData) error {
	//TODO
	return nil
}

func (ct *CrmDzXfTestTask) checkTask(config *taskConfigRepository.CrmDzXfTestTaskConfigData) error {
	exConfig, ok := crmDzXfTestConfigList[config.FId]
	if !ok {
		return ct.addTask(config)
	}
	if exConfig.IsEqual(config) {
		return nil
	}
	ct.removeTask(config.FId)
	return ct.addTask(config)
}

func (ct *CrmDzXfTestTask) removeTask(id string) {
	ct.clearConfigList(id)
	ct.clearTaskList(id)
	ct.delTaskRunningState(id)
	return
}

func (ct *CrmDzXfTestTask) clearConfigList(id string) {
	config, ok := crmDzXfTestConfigList[id]
	if !ok {
		log.Warn(fmt.Sprintf("remove task :config is not exist,taskId[%s]", id))
		return
	}
	configStr, err := goToolCommon.GetJsonStr(config)
	if err != nil {
		log.Warn(fmt.Sprintf("Del CrmDzXfTest Task，转换配置内容时遇到错误:%s，configID：%s", configStr, config.FId))
	} else {
		log.Warn(fmt.Sprintf("Del CrmDzXfTest Task:%s", configStr))
	}
	delete(crmDzXfTestConfigList, id)
}

func (ct *CrmDzXfTestTask) clearTaskList(id string) {
	c, ok := crmDzXfTestTaskList[id]
	if !ok {
		log.Warn(fmt.Sprintf("remove task :task is not exist,taskId[%s]", id))
		return
	}
	c.Stop()
	delete(crmDzXfTestTaskList, id)
}

func (ct *CrmDzXfTestTask) setTaskRunningState(id string, s bool) {
	crmDzXfTestTaskState[id] = s
}

func (ct *CrmDzXfTestTask) delTaskRunningState(id string) {
	delete(crmDzXfTestTaskState, id)
}
