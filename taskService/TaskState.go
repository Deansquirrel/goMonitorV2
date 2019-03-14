package taskService

import (
	"fmt"
	"github.com/Deansquirrel/goMonitorV2/taskConfigRepository"
	"github.com/robfig/cron"
)
import log "github.com/Deansquirrel/goToolLog"

type taskState struct {
	Health []*healthTaskSnap
	Int    []*intTaskSnap
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

func NewTaskStateSnap() *taskState {
	taskMConfigRepository := taskConfigRepository.TaskMConfig{}
	health := make([]*healthTaskSnap, 0)
	for key, c := range healthTaskList {
		h := &healthTaskSnap{}
		h.C = c
		config, ok := healthConfigList[key]
		if ok {
			c, err := taskMConfigRepository.GetMConfig(config.FId)
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
	int := make([]*intTaskSnap, 0)
	for key, c := range intTaskList {
		i := &intTaskSnap{}
		i.C = c
		config, ok := intConfigList[key]
		if ok {
			c, err := taskMConfigRepository.GetMConfig(config.FId)
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
	return &taskState{
		Health: health,
		Int:    int,
	}
}
