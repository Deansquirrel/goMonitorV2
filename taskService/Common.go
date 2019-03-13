package taskService

import (
	"github.com/Deansquirrel/goMonitorV2/taskConfigRepository"
	"github.com/robfig/cron"
)

var intConfigList map[string]*taskConfigRepository.IntTaskConfigData
var intTaskList map[string]*cron.Cron

var healthConfigList map[string]*taskConfigRepository.HealthTaskConfigData
var healthTaskList map[string]*cron.Cron

func init() {
	intConfigList = make(map[string]*taskConfigRepository.IntTaskConfigData)
	intTaskList = make(map[string]*cron.Cron)

	healthConfigList = make(map[string]*taskConfigRepository.HealthTaskConfigData)
	healthTaskList = make(map[string]*cron.Cron)
}
