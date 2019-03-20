package taskConfigRepository

import (
	"database/sql"
	"github.com/Deansquirrel/goToolCommon"
	"time"
)
import log "github.com/Deansquirrel/goToolLog"

const SqlGetCrmDzXfTestTaskHis = "" +
	"SELECT [FId], [FConfigId], [FUseTime], [FHttpCode], [FContent], [FOprTime]" +
	" FROM  CrmDzXfTestTaskHis"

const SqlGetCrmDzXfTestTaskHisById = "" +
	"SELECT [FId], [FConfigId], [FUseTime], [FHttpCode], [FContent], [FOprTime]" +
	" FROM  CrmDzXfTestTaskHis" +
	" WHERE FId = ?"

const SqlGetCrmDzXfTestTaskHisByConfigId = "" +
	"SELECT [FId], [FConfigId], [FUseTime], [FHttpCode], [FContent], [FOprTime]" +
	" FROM  CrmDzXfTestTaskHis" +
	" WHERE FConfigId = ?"

const SqlSetCrmDzXfTestTaskHis = "" +
	"INSERT INTO CrmDzXfTestTaskHis (FId, FConfigId, FUseTime, FHttpCode, FContent)" +
	" VALUES (?, ?, ?, ?, ?)"

const SqlDelCrmDzXfTestTaskHis = "" +
	"DELETE FROM CrmDzXfTestTaskHis" +
	" WHERE FOprTime < ?"

type CrmDzXfTestHis struct {
}

type CrmDzXfTestHisData struct {
	FId       string
	FConfigId string
	FUseTime  int
	FHttpCode int
	FContent  string
	FOprTime  time.Time
}

func (ch *CrmDzXfTestHis) GetCrmDzXfTestTaskHisList() ([]*CrmDzXfTestHisData, error) {
	rows, err := comm.getRowsBySQL(SqlGetCrmDzXfTestTaskHis)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return ch.getCrmDzXfTestTaskHisListByRows(rows)
}

func (ch *CrmDzXfTestHis) GetCrmDzXfTestTaskHisListById(id string) ([]*CrmDzXfTestHisData, error) {
	rows, err := comm.getRowsBySQL(SqlGetCrmDzXfTestTaskHisById, id)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return ch.getCrmDzXfTestTaskHisListByRows(rows)
}

func (ch *CrmDzXfTestHis) GetCrmDzXfTestTaskHisListByConfigId(id string) ([]*CrmDzXfTestHisData, error) {
	rows, err := comm.getRowsBySQL(SqlGetCrmDzXfTestTaskHisByConfigId, id)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return ch.getCrmDzXfTestTaskHisListByRows(rows)
}

func (ch *CrmDzXfTestHis) SetCrmDzXfTestTaskHis(data *CrmDzXfTestHisData) error {
	err := comm.setRowsBySQL(SqlSetCrmDzXfTestTaskHis,
		data.FId, data.FConfigId, data.FUseTime, data.FHttpCode, data.FContent)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (ch *CrmDzXfTestHis) ClearCrmDzXfTestTaskHis(t time.Duration) error {
	dateP := goToolCommon.GetDateTimeStr(time.Now().Add(-t))
	err := comm.setRowsBySQL(SqlDelCrmDzXfTestTaskHis, dateP)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (ch *CrmDzXfTestHis) getCrmDzXfTestTaskHisListByRows(rows *sql.Rows) ([]*CrmDzXfTestHisData, error) {
	defer func() {
		_ = rows.Close()
	}()
	var fId, fConfigId, fContent string
	var fUseTime, fHttpCode int
	var fOprTime time.Time
	resultList := make([]*CrmDzXfTestHisData, 0)
	var err error
	for rows.Next() {
		err = rows.Scan(&fId, &fConfigId, &fUseTime, &fHttpCode, &fContent, &fOprTime)
		if err != nil {
			break
		}
		config := CrmDzXfTestHisData{
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
