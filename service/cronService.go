package service

import (
	"github.com/Deansquirrel/goMonitorV2/global"
	"github.com/Deansquirrel/goMonitorV2/taskService"
	"github.com/Deansquirrel/goMonitorV2/webService"
	"os"
	"os/signal"
	"syscall"
)
import log "github.com/Deansquirrel/goToolLog"

type cronService struct {
}

func NewCronService() *cronService {
	return &cronService{}
}

func (cs *cronService) Start() {
	log.Debug("CronService starting")
	defer log.Debug("CronService start complete")
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch,
			os.Interrupt,
			syscall.SIGINT,
			os.Kill,
			syscall.SIGKILL,
			syscall.SIGTERM,
		)
		select {
		case <-ch:
			defer global.Cancel()
		case <-global.Ctx.Done():
		}
	}()

	var err error
	intTask := taskService.IntTask{}
	err = intTask.StartTask()
	if err != nil {
		log.Error(err.Error())
	}
	healthTask := taskService.HealthTask{}
	err = healthTask.StartTask()
	if err != nil {
		log.Error(err.Error())
	}

	ws := webService.NewWebServer(global.SysConfig.IrisConfig.Port)
	ws.StartWebService()
}
