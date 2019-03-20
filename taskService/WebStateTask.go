package taskService

import (
	"errors"
	"fmt"
	"github.com/Deansquirrel/goMonitorV2/global"
	"github.com/Deansquirrel/goMonitorV2/taskConfigRepository"
	"github.com/Deansquirrel/goMonitorV2/worker"
	"github.com/Deansquirrel/goToolCommon"
	"github.com/robfig/cron"
	"time"
)
import log "github.com/Deansquirrel/goToolLog"

type WebStateTask struct {
}

func (wst *WebStateTask) StartTask() error {
	webStateTaskConfigRepository := taskConfigRepository.WebStateTaskConfig{}
	//获取配置列表
	list, err := webStateTaskConfigRepository.GetWebStateTaskConfigList()
	if err != nil {
		return err
	}
	//清空Config列表
	webStateConfigList = make(map[string]*taskConfigRepository.WebStateTaskConfigData, 0)
	//清空Task列表
	if len(webStateTaskList) > 0 {
		keyList := make([]string, 0)
		for key := range webStateTaskList {
			keyList = append(keyList, key)
		}
		for _, key := range keyList {
			webStateTaskList[key].Stop()
			delete(webStateConfigList, key)
		}
	}
	errMsg := ""
	errMsgFormat := "添加WebState任务[%s]报错：%s；"
	for _, config := range list {
		err = wst.addTask(config)
		if err != nil {
			errMsg = errMsg + fmt.Sprintf(errMsgFormat, config.FId, err.Error())
		}
	}
	if errMsg != "" {
		log.Error(errMsg)
		return errors.New(errMsg)
	}
	wst.startRegularRefresh()
	return nil
}
func (wst *WebStateTask) StartJob(id string) error {
	c, ok := webStateTaskList[id]
	if !ok {
		return errors.New("task is not exist")
	}
	c.Start()
	wst.setTaskRunningState(id, true)
	return nil
}
func (wst *WebStateTask) StopJob(id string) error {
	c, ok := webStateTaskList[id]
	if !ok {
		return errors.New("task is not exist")
	}
	c.Stop()
	wst.setTaskRunningState(id, false)
	return nil
}

func (wst *WebStateTask) startRegularRefresh() {
	c := cron.New()
	var err error
	err = c.AddFunc("0 0/1 * * * ?", wst.refreshConfig)
	if err != nil {
		log.Error("添加WebState配置刷新任务时遇到错误：" + err.Error())
	} else {
		log.Info("添加WebState配置刷新任务完成")
	}
	err = c.AddFunc("0 0 0 * * ?", wst.delHisData)
	if err != nil {
		log.Error("添加删除WebState历史数据任务时遇到错误：" + err.Error())
	} else {
		log.Info("添加删除WebState历史数据任务完成")
	}
	c.Start()
}

func (wst *WebStateTask) refreshConfig() {
	err := wst.RefreshConfig()
	if err != nil {
		log.Error("刷新WebState配置时遇到错误：" + err.Error())
	}
}

func (wst *WebStateTask) delHisData() {
	webStateTaskHis := taskConfigRepository.WebStateTaskHis{}
	d := time.Duration(1000 * 1000 * 1000 * 60 * 60 * 24 * global.SysConfig.TaskConfig.KeepDays)
	_ = webStateTaskHis.ClearWebStateTaskHis(d)
}

func (wst *WebStateTask) RefreshConfig() error {
	rep := taskConfigRepository.WebStateTaskConfig{}
	//获取配置列表
	list, err := rep.GetWebStateTaskConfigList()
	if err != nil {
		return err
	}
	listId := make([]string, 0)
	mapId := make(map[string]*taskConfigRepository.WebStateTaskConfigData, 0)
	for _, config := range list {
		listId = append(listId, config.FId)
		mapId[config.FId] = config
	}
	configId := make([]string, 0)
	for key := range webStateConfigList {
		configId = append(configId, key)
	}

	addList, delList, checkList := goToolCommon.CheckDiff(listId, configId)
	errMsg := ""
	errMsgFormat := "刷新WebState任务[%s]报错：%s；"

	for _, id := range addList {
		err = wst.addTask(mapId[id])
		if err != nil {
			errMsg = errMsg + fmt.Sprintf(errMsgFormat, id, err.Error())
		}
	}
	for _, id := range delList {
		wst.removeTask(id)
	}
	for _, id := range checkList {
		err = wst.checkTask(mapId[id])
		if err != nil {
			errMsg = errMsg + fmt.Sprintf(errMsgFormat, id, err.Error())
		}
	}
	if errMsg != "" {
		return errors.New(errMsg)
	}
	return nil
}

func (wst *WebStateTask) addTask(config *taskConfigRepository.WebStateTaskConfigData) error {
	configStr, err := goToolCommon.GetJsonStr(config)
	if err != nil {
		log.Warn(fmt.Sprintf("Add WebState Task，转换配置内容时遇到错误:%s，configID：%s", configStr, config.FId))
	} else {
		log.Warn(fmt.Sprintf("Add WebState Task:%s", configStr))
	}
	//------------------------------------------------------------------------------------------------------------------
	webStateConfigList[config.FId] = config
	w := worker.NewWebStateWorker(config)
	c := cron.New()
	err = c.AddJob(config.FCron, w)
	if err != nil {
		log.Error(err.Error())
		wst.setTaskRunningState(config.FId, false)
	} else {
		c.Start()
		wst.setTaskRunningState(config.FId, true)
	}
	webStateTaskList[config.FId] = c
	return err
}

func (wst *WebStateTask) checkTask(config *taskConfigRepository.WebStateTaskConfigData) error {
	exConfig, ok := webStateConfigList[config.FId]
	if !ok {
		return wst.addTask(config)
	}
	if exConfig.IsEqual(config) {
		return nil
	}
	wst.removeTask(config.FId)
	return wst.addTask(config)
}

func (wst *WebStateTask) removeTask(id string) {
	wst.clearConfigList(id)
	wst.clearTaskList(id)
	wst.delTaskRunningState(id)
	return
}

func (wst *WebStateTask) clearConfigList(id string) {
	config, ok := webStateConfigList[id]
	if !ok {
		log.Warn(fmt.Sprintf("remove task :config is not exist,taskId[%s]", id))
		return
	}
	configStr, err := goToolCommon.GetJsonStr(config)
	if err != nil {
		log.Warn(fmt.Sprintf("Del WebState Task，转换配置内容时遇到错误:%s，configID：%s", configStr, config.FId))
	} else {
		log.Warn(fmt.Sprintf("Del WebState Task:%s", configStr))
	}
	delete(webStateConfigList, id)
}

func (wst *WebStateTask) clearTaskList(id string) {
	c, ok := webStateTaskList[id]
	if !ok {
		log.Warn(fmt.Sprintf("remove task :task is not exist,taskId[%s]", id))
		return
	}
	c.Stop()
	delete(webStateTaskList, id)
}

func (wst *WebStateTask) setTaskRunningState(id string, s bool) {
	webStateTaskState[id] = s
}

func (wst *WebStateTask) delTaskRunningState(id string) {
	delete(webStateTaskState, id)
}
