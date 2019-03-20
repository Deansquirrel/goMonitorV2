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

var crmDzXfTestConfigList map[string]*taskConfigRepository.CrmDzXfTestTaskConfigData
var crmDzXfTestTaskList map[string]*cron.Cron
var crmDzXfTestTaskState map[string]bool

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

	crmDzXfTestConfigList = make(map[string]*taskConfigRepository.CrmDzXfTestTaskConfigData)
	crmDzXfTestTaskList = make(map[string]*cron.Cron)
	crmDzXfTestTaskState = make(map[string]bool)
}
