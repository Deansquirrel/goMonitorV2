package taskService

import (
	"fmt"
	"github.com/Deansquirrel/goMonitorV2/taskConfigRepository"
	"github.com/robfig/cron"
)
import log "github.com/Deansquirrel/goToolLog"

type taskState struct {
	Health                []*healthTaskSnap
	Int                   []*intTaskSnap
	WebState              []*webStateTaskSnap
	taskMConfigRepository taskConfigRepository.TaskMConfig
}

type healthTaskSnap struct {
	Config    *taskConfigRepository.TaskMConfigData
	IsRunning bool
	C         *cron.Cron
}

type intTaskSnap struct {
	Config    *taskConfigRepository.TaskMConfigData
	IsRunning bool
	C         *cron.Cron
}

type webStateTaskSnap struct {
	Config    *taskConfigRepository.TaskMConfigData
	IsRunning bool
	C         *cron.Cron
}

func NewTaskStateSnap() *taskState {
	ts := &taskState{
		taskMConfigRepository: taskConfigRepository.TaskMConfig{},
	}
	ts.Health = ts.getHealthTaskSnap()
	ts.Int = ts.getIntTaskSnap()
	ts.WebState = ts.getWebStateTaskSnap()
	return ts
}

func (ts *taskState) getIntTaskSnap() []*intTaskSnap {
	int := make([]*intTaskSnap, 0)
	for key, c := range intTaskList {
		i := &intTaskSnap{}
		i.C = c
		config, ok := intConfigList[key]
		if ok {
			c, err := ts.taskMConfigRepository.GetMConfig(config.FId)
			if err != nil {
				log.Error(err.Error())
			} else {
				if len(c) != 1 {
					log.Error(fmt.Sprintf("MConfig 查询数量异常，exp：1，act：%d", len(c)))
				} else {
					i.Config = c[0]
				}
			}
		}
		isRunning, ok := intTaskState[key]
		if ok {
			i.IsRunning = isRunning
		}
		int = append(int, i)
	}
	return int
}

func (ts *taskState) getHealthTaskSnap() []*healthTaskSnap {
	health := make([]*healthTaskSnap, 0)
	for key, c := range healthTaskList {
		h := &healthTaskSnap{}
		h.C = c
		config, ok := healthConfigList[key]
		if ok {
			c, err := ts.taskMConfigRepository.GetMConfig(config.FId)
			if err != nil {
				log.Error(err.Error())
			} else {
				if len(c) != 1 {
					log.Error(fmt.Sprintf("MConfig 查询数量异常，exp：1，act：%d", len(c)))
				} else {
					h.Config = c[0]
				}
			}
		}
		isRunning, ok := healthTaskState[key]
		if ok {
			h.IsRunning = isRunning
		}
		health = append(health, h)
	}
	return health
}

func (ts *taskState) getWebStateTaskSnap() []*webStateTaskSnap {
	webState := make([]*webStateTaskSnap, 0)
	for key, c := range webStateTaskList {
		h := &webStateTaskSnap{}
		h.C = c
		config, ok := webStateConfigList[key]
		if ok {
			c, err := ts.taskMConfigRepository.GetMConfig(config.FId)
			if err != nil {
				log.Error(err.Error())
			} else {
				if len(c) != 1 {
					log.Error(fmt.Sprintf("MConfig 查询数量异常，exp：1，act：%d", len(c)))
				} else {
					h.Config = c[0]
				}
			}
		}
		isRunning, ok := webStateTaskState[key]
		if ok {
			h.IsRunning = isRunning
		}
		webState = append(webState, h)
	}
	return webState
}
