package taskConfigRepository

import (
	"database/sql"
	"github.com/Deansquirrel/goToolCommon"
	"time"
)
import log "github.com/Deansquirrel/goToolLog"

const SqlGetIntTaskHis = "" +
	"SELECT [FId],[FConfigId],[FNum],[FContent],[FOprTime]" +
	" FROM [IntTaskHis]"

const SqlGetIntTaskHisById = "" +
	"SELECT [FId],[FConfigId],[FNum],[FContent],[FOprTime]" +
	" FROM [IntTaskHis]" +
	" WHERE [FId]=?"

const SqlGetIntTaskHisByConfigId = "" +
	"SELECT [FId],[FConfigId],[FNum],[FContent],[FOprTime]" +
	" FROM [IntTaskHis]" +
	" WHERE [FConfigId] = ?" +
	" Order By [FOprTime] Asc"

const SqlSetIntTaskHis = "" +
	"INSERT INTO [IntTaskHis]([FId],[FConfigId],[FNum],[FContent])" +
	" SELECT ?,?,?,?"

const SqlDelIntTaskHisByOprTime = "" +
	"DELETE FROM [IntTaskHis]" +
	" WHERE [FOprTime] < ?"

type IntTaskHis struct {
}

type IntTaskHisData struct {
	FId       string
	FConfigId string
	FNum      int
	FContent  string
	FOprTime  time.Time
}

func (ith *IntTaskHis) GetIntTaskHisList() ([]*IntTaskHisData, error) {
	rows, err := comm.getRowsBySQL(SqlGetIntTaskHis)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return ith.getIntTaskHisListByRows(rows)
}

func (ith *IntTaskHis) GetIntTaskHisListById(id string) ([]*IntTaskHisData, error) {
	rows, err := comm.getRowsBySQL(SqlGetIntTaskHisById, id)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return ith.getIntTaskHisListByRows(rows)
}

func (ith *IntTaskHis) GetIntTaskHisListByConfigId(id string) ([]*IntTaskHisData, error) {
	rows, err := comm.getRowsBySQL(SqlGetIntTaskHisByConfigId, id)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return ith.getIntTaskHisListByRows(rows)
}

func (ith *IntTaskHis) SetIntTaskHis(data *IntTaskHisData) error {
	err := comm.setRowsBySQL(SqlSetIntTaskHis, data.FId, data.FConfigId, data.FNum, data.FContent)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (ith *IntTaskHis) ClearIntTaskHis(t time.Duration) error {
	dateP := goToolCommon.GetDateTimeStr(time.Now().Add(-t))
	err := comm.setRowsBySQL(SqlDelIntTaskHisByOprTime, dateP)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (ith *IntTaskHis) getIntTaskHisListByRows(rows *sql.Rows) ([]*IntTaskHisData, error) {
	defer func() {
		_ = rows.Close()
	}()
	var fId, fConfigId, fContent string
	var fNum int
	var fOprTime time.Time
	resultList := make([]*IntTaskHisData, 0)
	var err error
	for rows.Next() {
		err = rows.Scan(&fId, &fConfigId, &fNum, &fContent, &fOprTime)
		if err != nil {
			break
		}
		config := IntTaskHisData{
			FId:       fId,
			FConfigId: fConfigId,
			FNum:      fNum,
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
