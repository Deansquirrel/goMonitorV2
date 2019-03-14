package taskService

import (
	"errors"
	"fmt"
	"github.com/Deansquirrel/goMonitorV2/global"
	"github.com/Deansquirrel/goMonitorV2/taskConfigRepository"
	"github.com/Deansquirrel/goMonitorV2/worker"
	"github.com/Deansquirrel/goToolCommon"
	log "github.com/Deansquirrel/goToolLog"
	"github.com/robfig/cron"
	"time"
)

type IntTask struct {
}

func (it *IntTask) StartTask() error {
	intTaskConfigRepository := taskConfigRepository.IntTaskConfig{}
	//获取配置列表
	list, err := intTaskConfigRepository.GetIntTaskConfigList()
	if err != nil {
		return err
	}
	//清空Config列表
	intConfigList = make(map[string]*taskConfigRepository.IntTaskConfigData)
	//清空Task列表
	if len(intTaskList) > 0 {
		keyList := make([]string, 0)
		for key := range intTaskList {
			keyList = append(keyList, key)
		}
		for _, key := range keyList {
			intTaskList[key].Stop()
			delete(intConfigList, key)
		}
	}
	//缓存配置列表、任务列表
	errMsg := ""
	errMsgFormat := "添加Int任务[%s]报错：%s；"
	for _, config := range list {
		err = it.addTask(config)
		if err != nil {
			errMsg = errMsg + fmt.Sprintf(errMsgFormat, config.FId, err.Error())
		}
	}
	if errMsg != "" {
		log.Error(errMsg)
		return errors.New(errMsg)
	}
	it.startRegularRefresh()
	return nil
}

func (it *IntTask) StartJob(id string) error {
	c, ok := intTaskList[id]
	if !ok {
		return errors.New("task is not exist")
	}
	c.Start()
	it.setTaskRunningState(id, true)
	return nil
}

func (it *IntTask) StopJob(id string) error {
	c, ok := intTaskList[id]
	if !ok {
		return errors.New("task is not exist")
	}
	c.Stop()
	it.setTaskRunningState(id, false)
	return nil
}

func (it *IntTask) startRegularRefresh() {
	time.AfterFunc(time.Second*30, it.startRegularRefresh)
	it.refreshConfig()
	it.delHisData()
}

func (it *IntTask) refreshConfig() {
	err := it.RefreshConfig()
	if err != nil {
		log.Error("刷新Int配置时遇到错误：" + err.Error())
	}
}

func (it *IntTask) delHisData() {
	intTaskHis := taskConfigRepository.IntTaskHis{}
	d := time.Duration(1000 * 1000 * 1000 * 60 * 60 * 24 * global.SysConfig.TaskConfig.KeepDays)
	_ = intTaskHis.ClearIntTaskHis(d)
}

func (it *IntTask) RefreshConfig() error {
	intTaskConfigRepository := taskConfigRepository.IntTaskConfig{}
	//获取配置列表
	list, err := intTaskConfigRepository.GetIntTaskConfigList()
	if err != nil {
		return err
	}
	listId := make([]string, 0)
	mapId := make(map[string]*taskConfigRepository.IntTaskConfigData, 0)
	for _, config := range list {
		listId = append(listId, config.FId)
		mapId[config.FId] = config
	}

	configId := make([]string, 0)
	for key := range intConfigList {
		configId = append(configId, key)
	}

	addList, delList, checkList := goToolCommon.CheckDiff(listId, configId)

	errMsg := ""
	errMsgFormat := "刷新Int任务[%s]报错：%s；"

	for _, id := range addList {
		err = it.addTask(mapId[id])
		if err != nil {
			errMsg = errMsg + fmt.Sprintf(errMsgFormat, id, err.Error())
		}
	}
	for _, id := range delList {
		it.removeTask(id)
	}
	for _, id := range checkList {
		err = it.checkTask(mapId[id])
		if err != nil {
			errMsg = errMsg + fmt.Sprintf(errMsgFormat, id, err.Error())
		}
	}
	if errMsg != "" {
		return errors.New(errMsg)
	}
	return nil
}

func (it *IntTask) addTask(config *taskConfigRepository.IntTaskConfigData) error {
	//------------------------------------------------------------------------------------------------------------------
	configStr, err := goToolCommon.GetJsonStr(config)
	if err != nil {
		log.Warn(fmt.Sprintf("Add Int Task，转换配置内容时遇到错误:%s，configID：%s", configStr, config.FId))
	} else {
		log.Warn(fmt.Sprintf("Add Int Task:%s", configStr))
	}
	//------------------------------------------------------------------------------------------------------------------
	intConfigList[config.FId] = config
	w := worker.NewIntWorker(config)
	c := cron.New()
	err = c.AddJob(config.FCron, w)
	if err != nil {
		log.Error(err.Error())
		it.setTaskRunningState(config.FId, false)
	} else {
		c.Start()
		it.setTaskRunningState(config.FId, true)
	}
	intTaskList[config.FId] = c
	return err
}

func (it *IntTask) checkTask(config *taskConfigRepository.IntTaskConfigData) error {
	exConfig := intConfigList[config.FId]
	if exConfig == nil {
		return it.addTask(config)
	}
	if exConfig.IsEqual(config) {
		return nil
	}
	it.removeTask(config.FId)
	return it.addTask(config)
}

func (it *IntTask) removeTask(id string) {
	config := intConfigList[id]
	if config == nil {
		log.Warn(fmt.Sprintf("remove task :config is not exist,taskId[%s]", id))
	} else {
		configStr, err := goToolCommon.GetJsonStr(config)
		if err != nil {
			log.Warn(fmt.Sprintf("Del Int Task，转换配置内容时遇到错误:%s，configID：%s", configStr, config.FId))
		} else {
			log.Warn(fmt.Sprintf("Del Int Task:%s", configStr))
		}
		delete(intConfigList, id)
	}
	c := intTaskList[id]
	if c == nil {
		log.Warn(fmt.Sprintf("remove task :task is not exist,taskId[%s]", id))
	} else {
		c.Stop()
		delete(intTaskList, id)
	}
	it.delTaskRunningState(id)
}

func (it *IntTask) setTaskRunningState(id string, s bool) {
	intTaskState[id] = s
}

func (it *IntTask) delTaskRunningState(id string) {
	delete(intTaskState, id)
}
