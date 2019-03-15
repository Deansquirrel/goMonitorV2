package taskService

import (
	"github.com/Deansquirrel/goMonitorV2/taskConfigRepository"
	"github.com/robfig/cron"
)

var healthConfigList map[string]*taskConfigRepository.HealthTaskConfigData
var healthTaskList map[string]*cron.Cron
var healthTaskState map[string]bool

var intConfigList map[string]*taskConfigRepository.IntTaskConfigData
var intTaskList map[string]*cron.Cron
var intTaskState map[string]bool

var webStateConfigList map[string]*taskConfigRepository.WebStateTaskConfigData
var webStateTaskList map[string]*cron.Cron
var webStateTaskState map[string]bool

func init() {
	healthConfigList = make(map[string]*taskConfigRepository.HealthTaskConfigData)
	healthTaskList = make(map[string]*cron.Cron)
	healthTaskState = make(map[string]bool)

	intConfigList = make(map[string]*taskConfigRepository.IntTaskConfigData)
	intTaskList = make(map[string]*cron.Cron)
	intTaskState = make(map[string]bool)

	webStateConfigList = make(map[string]*taskConfigRepository.WebStateTaskConfigData)
	webStateTaskList = make(map[string]*cron.Cron)
	webStateTaskState = make(map[string]bool)
}
