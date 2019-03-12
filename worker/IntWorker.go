package worker

import (
	"database/sql"
	"fmt"
	"github.com/Deansquirrel/goMonitorV2/taskConfigRepository"
	log "github.com/Deansquirrel/goToolLog"
	"github.com/Deansquirrel/goToolMSSql"
	"github.com/kataras/iris/core/errors"
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
		comm.sendMsg(iw.intTaskConfigData.FId, comm.getMsg(iw.intTaskConfigData.FMsgTitle, strings.Replace(iw.intTaskConfigData.FMsgContent, "title", strconv.Itoa(num), -1)))
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
