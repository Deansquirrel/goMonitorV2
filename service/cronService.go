package service

import (
	"github.com/Deansquirrel/goMonitorV2/global"
	"github.com/Deansquirrel/goMonitorV2/taskService"
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
	cs.start()
	select {
	case <-global.Ctx.Done():
		return
	}
}

func (cs *cronService) start() {
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

	//intTaskConfig := taskConfigRepository.TaskMConfig{}
	//msg := "%d times %s"
	//for i:=0;i<1000000;i++ {
	//	list,err := intTaskConfig.GetMConfigList()
	//	if err != nil {
	//		log.Error(fmt.Sprintf(msg,i,err.Error()))
	//	} else {
	//		log.Debug(fmt.Sprintf(msg,i,strconv.Itoa(len(list))))
	//	}
	//	time.Sleep(time.Second * 5)
	//}
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

	//time.AfterFunc(time.Minute,func(){
	//	_ = intTask.StopJob("7235F3B3-E258-4548-9A14-A6B1313B0B49")
	//})
	//
	//time.AfterFunc(time.Minute * 2,func(){
	//	_ = intTask.StartJob("7235F3B3-E258-4548-9A14-A6B1313B0B49")
	//})

	cs.test()
}

//测试消息发送
func (cs *cronService) test() {
	//i := 0
	//c := cron.New()
	//spec := "0/1 * * * * ?"
	//_ = c.AddFunc(spec,func(){
	//	i++
	//	log.Debug(fmt.Sprintf("cron running:%d",i))
	//	for _,j := range c.Entries(){
	//		log.Debug(goToolCommon.GetDateTimeStr(j.Prev))
	//		log.Debug(goToolCommon.GetDateTimeStr(j.Next))
	//	}
	//})
	//time.AfterFunc(time.Second * 5,func(){
	//	c.Stop()
	//})
	//time.AfterFunc(time.Second * 10,func(){
	//	c.Start()
	//})
	//c.Start()
}
