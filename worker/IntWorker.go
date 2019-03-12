package worker

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Deansquirrel/goMonitorV2/taskConfigRepository"
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
	if num > iw.intTaskConfigData.FCheckMax || num < iw.intTaskConfigData.FCheckMin {
		msg := comm.getMsg(iw.intTaskConfigData.FMsgTitle, strings.Replace(iw.intTaskConfigData.FMsgContent, "title", strconv.Itoa(num), -1))
		dMsg := iw.getDMsg()
		if dMsg != "" {
			msg = msg + "\n" + dMsg
		}
		comm.sendMsg(iw.intTaskConfigData.FId, msg)
	}
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
			break
		} else {
			list = append(list, num)
		}
	}
	if err != nil {
		return 0, err
	}
	if len(list) != 1 {
		return 0, errors.New(fmt.Sprintf("SQL返回数量异常，exp:1,act:%d", len(list)))
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
		return nil, err
	}
	rows, err := conn.Query(sql)
	if err != nil {
		log.Debug(err.Error())
		return nil, err
	}
	return rows, nil
}

//获取详细消息
func (iw *intWorker) getDMsg() string {
	intTaskDConfig := taskConfigRepository.IntTaskDConfig{}
	dConfigList, err := intTaskDConfig.GetIntTaskDConfig(iw.intTaskConfigData.FId)
	if err != nil {
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
	values := make([]sql.RawBytes, len(titleList))
	iList := make([]interface{}, len(titleList))
	for i := range values {
		iList[i] = &values[i]
	}
	result := ""
	for rows.Next() {
		if result != "" {
			result = result + "\n" + "--------------------"
		}
		err = rows.Scan(iList...)
		if err != nil {
			return fmt.Sprintf("读取明细内容时遇到错误：%s，查询语句为：%s", err.Error(), search)
		}
		for i := 0; i < len(titleList); i++ {
			if result != "" {
				result = result + "\n"
			}
			var v string
			if values[i] == nil {
				v = "NULL"
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
