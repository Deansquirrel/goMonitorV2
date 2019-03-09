package worker

import (
	"database/sql"
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
	rows, err := iw.getRowsBySQL(iw.intTaskConfigData.FSearch)
	if err != nil {
		comm.sendMsg(iw.intTaskConfigData.FId, comm.getMsg(iw.intTaskConfigData.FMsgTitle, err.Error()))
		return
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
		log.Error(err.Error())
		comm.sendMsg(iw.intTaskConfigData.FId, comm.getMsg(iw.intTaskConfigData.FMsgTitle, err.Error()))
		return
	}
	if rows.Err() != nil {
		log.Error(rows.Err().Error())
		comm.sendMsg(iw.intTaskConfigData.FId, comm.getMsg(iw.intTaskConfigData.FMsgTitle, rows.Err().Error()))
	}
	if len(list) != 1 {
		comm.sendMsg(iw.intTaskConfigData.FId, comm.getMsg(iw.intTaskConfigData.FMsgTitle, fmt.Sprintf("SQL返回数量异常，exp:1,act:%d", len(list))))
		return
	}
	num = list[0]
	if num > iw.intTaskConfigData.FCheckMax || num < iw.intTaskConfigData.FCheckMin {
		comm.sendMsg(iw.intTaskConfigData.FId, comm.getMsg(iw.intTaskConfigData.FMsgTitle, strings.Replace(iw.intTaskConfigData.FMsgContent, "title", strconv.Itoa(num), -1)))
	}
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
