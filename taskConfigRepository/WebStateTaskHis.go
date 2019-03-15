package taskConfigRepository

import (
	"database/sql"
	"github.com/Deansquirrel/goToolCommon"
	"time"
)
import log "github.com/Deansquirrel/goToolLog"

const SqlGetWebStateTaskHis = "" +
	"SELECT FId, FConfigId, FUseTime, FHttpCode, FContent, FOprTime" +
	" FROM WebStateTaskHis"
const SqlGetWebStateTaskHisById = "" +
	"SELECT FId, FConfigId, FUseTime, FHttpCode, FContent, FOprTime" +
	" FROM WebStateTaskHis" +
	" WHERE FId = ?"
const SqlGetWebStateTaskHisByConfigId = "" +
	"SELECT FId, FConfigId, FUseTime, FHttpCode, FContent, FOprTime" +
	" FROM WebStateTaskHis" +
	" WHERE FConfigId = ?"
const SqlSetWebStateTaskHis = "" +
	"INSERT INTO WebStateTaskHis (FId, FConfigId, FUseTime, FHttpCode, FContent)" +
	" VALUES (?,?,?,?,?)"
const SqlDelWebStateTaskHis = "" +
	"DELETE FROM WebStateTaskHis" +
	" WHERE FOprTime < ?"

type WebStateTaskHis struct {
}

type WebStateTaskHisData struct {
	FId       string
	FConfigId string
	FUseTime  int
	FHttpCode int
	FContent  string
	FOprTime  time.Time
}

func (wth *WebStateTaskHis) GetWebStateTaskHisList() ([]*WebStateTaskHisData, error) {
	rows, err := comm.getRowsBySQL(SqlGetWebStateTaskHis)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return wth.getWebStateTaskHisListByRows(rows)
}

func (wth *WebStateTaskHis) GetWebStateTaskHisListById(id string) ([]*WebStateTaskHisData, error) {
	rows, err := comm.getRowsBySQL(SqlGetWebStateTaskHisById, id)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return wth.getWebStateTaskHisListByRows(rows)
}

func (wth *WebStateTaskHis) GetWebStateTaskHisListByConfigId(id string) ([]*WebStateTaskHisData, error) {
	rows, err := comm.getRowsBySQL(SqlGetWebStateTaskHisByConfigId, id)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return wth.getWebStateTaskHisListByRows(rows)
}

func (wth *WebStateTaskHis) SetWebStateTaskHis(data *WebStateTaskHisData) error {
	err := comm.setRowsBySQL(SqlSetWebStateTaskHis,
		data.FId, data.FConfigId, data.FUseTime, data.FHttpCode, data.FContent)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (wth *WebStateTaskHis) ClearWebStateTaskHis(t time.Duration) error {
	dateP := goToolCommon.GetDateTimeStr(time.Now().Add(-t))
	err := comm.setRowsBySQL(SqlDelWebStateTaskHis, dateP)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (wth *WebStateTaskHis) getWebStateTaskHisListByRows(rows *sql.Rows) ([]*WebStateTaskHisData, error) {
	defer func() {
		_ = rows.Close()
	}()
	var fId, fConfigId, fContent string
	var fUseTime, fHttpCode int
	var fOprTime time.Time
	resultList := make([]*WebStateTaskHisData, 0)
	var err error
	for rows.Next() {
		err = rows.Scan(&fId, &fConfigId, &fUseTime, &fHttpCode, &fContent, &fOprTime)
		if err != nil {
			break
		}
		config := WebStateTaskHisData{
			FId:       fId,
			FConfigId: fConfigId,
			FUseTime:  fUseTime,
			FHttpCode: fHttpCode,
			FContent:  fContent,
			FOprTime:  fOprTime,
		}
		resultList = append(resultList, &config)
	}
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if rows.Err() != nil {
		log.Error(rows.Err().Error())
		return nil, rows.Err()
	}
	return resultList, nil
}
