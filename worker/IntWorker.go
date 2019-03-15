package worker

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Deansquirrel/goMonitorV2/taskConfigRepository"
	"github.com/Deansquirrel/goToolCommon"
	log "github.com/Deansquirrel/goToolLog"
	"github.com/Deansquirrel/goToolMSSql"
	"strconv"
	"strings"
)

type intWorker struct {
	intTaskConfigData *taskConfigRepository.IntTaskConfigData
}

func NewIntWorker(intTaskConfigData *taskConfigRepository.IntTaskConfigData) *intWorker {
	return &intWorker{
		intTaskConfigData: intTaskConfigData,
	}
}

//检查
func (iw *intWorker) Run() {
	if iw.intTaskConfigData == nil {
		return
	}
	num, err := iw.getCheckNum()
	if err != nil {
		comm.sendMsg(iw.intTaskConfigData.FId, comm.getMsg(iw.intTaskConfigData.FMsgTitle, err.Error()))
		return
	}
	var msg string
	if num > iw.intTaskConfigData.FCheckMax || num < iw.intTaskConfigData.FCheckMin {
		//msg = comm.getMsg(iw.intTaskConfigData.FMsgTitle, strings.Replace(iw.intTaskConfigData.FMsgContent, "title", strconv.Itoa(num), -1))
		//dMsg := iw.getDMsg()
		//if dMsg != "" {
		//	msg = msg + "\n" + dMsg
		//}
		msg = iw.getTMsg(num)
		dMsg := iw.getDMsg()
		if dMsg != "" {
			msg = msg + "\n" + dMsg
		}
		comm.sendMsg(iw.intTaskConfigData.FId, msg)
	}
	iw.saveSearchResult(num, msg)
}

//获取待检测值
func (iw *intWorker) getCheckNum() (int, error) {
	rows, err := iw.getRowsBySQL(iw.intTaskConfigData.FSearch)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = rows.Close()
	}()
	list := make([]int, 0)
	var num int
	for rows.Next() {
		err = rows.Scan(&num)
		if err != nil {
			log.Error(err.Error())
			break
		} else {
			list = append(list, num)
		}
	}
	if err != nil {
		return 0, err
	}
	if len(list) != 1 {
		errMsg := fmt.Sprintf("SQL返回数量异常，exp:1,act:%d", len(list))
		log.Error(errMsg)
		return 0, errors.New(errMsg)
	}
	return list[0], nil
}

//获取DB配置
func (iw *intWorker) getDBConfig() *goToolMSSql.MSSqlConfig {
	return &goToolMSSql.MSSqlConfig{
		Server: iw.intTaskConfigData.FServer,
		Port:   iw.intTaskConfigData.FPort,
		DbName: iw.intTaskConfigData.FDbName,
		User:   iw.intTaskConfigData.FDbUser,
		Pwd:    iw.intTaskConfigData.FDbPwd,
	}
}

//查询数据
func (iw *intWorker) getRowsBySQL(sql string) (*sql.Rows, error) {
	conn, err := goToolMSSql.GetConn(iw.getDBConfig())
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	rows, err := conn.Query(sql)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return rows, nil
}

//获取主信息
func (iw *intWorker) getTMsg(num int) string {
	return comm.getMsg(iw.intTaskConfigData.FMsgTitle, iw.getTMsgContent(num))
}

func (iw *intWorker) getTMsgContent(num int) string {
	msgContent := iw.intTaskConfigData.FMsgContent
	msgContent = strings.Replace(msgContent, "title", strconv.Itoa(num), -1)
	return msgContent
}

//获取详细消息
func (iw *intWorker) getDMsg() string {
	intTaskDConfig := taskConfigRepository.IntTaskDConfig{}
	dConfigList, err := intTaskDConfig.GetIntTaskDConfig(iw.intTaskConfigData.FId)
	if err != nil {
		log.Error("获取明细信息查询配置时遇到错误：" + err.Error())
		return err.Error()
	}
	var msg, result string
	for _, dConfig := range dConfigList {
		msg = iw.getSingleDMsg(dConfig.FMsgSearch)
		if msg != "" {
			if result != "" {
				result = result + "\n"
			}
			result = result + "--------------------"
			result = result + "\n" + msg
		}
	}
	return result
}

func (iw *intWorker) getSingleDMsg(search string) string {
	if search == "" {
		return ""
	}
	rows, err := iw.getRowsBySQL(search)
	if err != nil {
		return fmt.Sprintf("查询明细内容时遇到错误：%s，查询语句为：%s", err.Error(), search)
	}
	defer func() {
		_ = rows.Close()
	}()
	titleList, err := rows.Columns()
	if err != nil {
		return fmt.Sprintf("获取明细内容表头时遇到错误：%s，查询语句为：%s", err.Error(), search)
	}
	counter := len(titleList)
	values := make([]sql.RawBytes, counter)
	valuePointers := make([]interface{}, counter)
	for i := 0; i < counter; i++ {
		valuePointers[i] = &values[i]
	}
	var result string
	for rows.Next() {
		err = rows.Scan(valuePointers...)
		if err != nil {
			return fmt.Sprintf("读取明细内容时遇到错误：%s，查询语句为：%s", err.Error(), search)
		}
		if result != "" {
			result = result + "\n" + "--------------------"
		}
		for i := 0; i < counter; i++ {
			if result != "" {
				result = result + "\n"
			}
			var v string
			if values[i] == nil {
				v = "Null"
			} else {
				v = string(values[i])
			}
			result = result + titleList[i] + " - " + v
		}
	}
	if rows.Err() != nil {
		return fmt.Sprintf("读取明细内容时遇到错误：%s，查询语句为：%s", err.Error(), search)
	}
	return result
}

//保存查询结果
func (iw *intWorker) saveSearchResult(num int, content string) {
	nData := &taskConfigRepository.IntTaskHisData{
		FId:       strings.ToUpper(goToolCommon.Guid()),
		FConfigId: iw.intTaskConfigData.FId,
		FNum:      num,
		FContent:  content,
	}
	intTaskHis := taskConfigRepository.IntTaskHis{}
	_ = intTaskHis.SetIntTaskHis(nData)
}
