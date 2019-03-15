package worker

import (
	"fmt"
	"github.com/Deansquirrel/goMonitorV2/global"
	"github.com/Deansquirrel/goMonitorV2/notify"
	"github.com/Deansquirrel/goMonitorV2/taskConfigRepository"
	"github.com/Deansquirrel/goToolCommon"
	log "github.com/Deansquirrel/goToolLog"
	"strings"
	"time"
)

var comm common
var notifyList taskConfigRepository.NotifyList
var dingTalkRobot taskConfigRepository.DingTalkRobot

func init() {
	comm = common{}
	notifyList = taskConfigRepository.NotifyList{}
	dingTalkRobot = taskConfigRepository.DingTalkRobot{}
}

type common struct {
}

//获取待发送消息
func (c *common) getMsg(title, content string) string {
	msg := ""
	titleList := strings.Split(title, "###")
	if len(titleList) > 0 {
		for _, t := range titleList {
			if strings.Trim(t, " ") != "" {
				if msg != "" {
					msg = msg + "\n"
				}
				msg = msg + strings.Trim(t, " ")
			}
		}
	}
	contentList := strings.Split(content, "###")
	if len(contentList) > 0 {
		for _, t := range contentList {
			if strings.Trim(t, " ") != "" {
				if msg != "" {
					msg = msg + "\n"
				}
				msg = msg + strings.Trim(t, " ")
			}
		}
	}
	msg = goToolCommon.GetDateTimeStr(time.Now()) + "\n" + msg
	return msg
}

//发送消息
func (c *common) sendMsg(taskId, msg string) {
	if msg == "" {
		return
	}
	notifyFlag := false
	defer func() {
		if !notifyFlag {
			log.Warn(fmt.Sprintf("消息未发送，taskId：%s，msg：%s", taskId, msg))
		}
	}()
	nData, err := notifyList.GetNotifyList(taskId)
	if err != nil {
		log.Error(fmt.Sprintf("获取通知信息时发送错误:%s", err.Error()))
	}
	//发送钉钉消息
	for _, s := range nData.DingTalkRobot {
		dingTalkRobotConfig, err := dingTalkRobot.GetDingTalkRobot(s)
		if err != nil {
			log.Error(fmt.Sprintf("获取DingTalkRobot时发生错误，taskId：%s，err：%s", s, err.Error()))
			continue
		}

		for _, config := range dingTalkRobotConfig {
			err = c.sendDingTalkRobotMsg(config, msg)
			if err != nil {
				log.Error(fmt.Sprintf("DingTalkRobot发送消息时发生错误，taskId：%s，dingTalkRobotId：%s，err：%s", s, config.FId, err.Error()))
				continue
			}
			notifyFlag = true
		}
	}
	return
}

func (c *common) sendDingTalkRobotMsg(config *taskConfigRepository.DingTalkRobotData, msg string) error {
	dingTalkRobotNotify := notify.NewDingTalkRobot(global.SysConfig.DingTalkConfig.Address)
	if config.FIsAtAll == 1 {
		return dingTalkRobotNotify.SendTextMsgWithAtAll(config.FWebHookKey, msg)
	} else {
		atList := strings.Split(config.FAtMobiles, ",")
		atList = goToolCommon.ClearBlock(atList)
		if len(atList) > 0 {
			return dingTalkRobotNotify.SendTextMsgWithAtList(config.FWebHookKey, msg, atList)
		} else {
			return dingTalkRobotNotify.SendTextMsg(config.FWebHookKey, msg)
		}
	}
}
