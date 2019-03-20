package taskConfigRepository

import (
	"database/sql"
	log "github.com/Deansquirrel/goToolLog"
)

const SqlGetCrmDzXfTestTaskConfig = "" +
	"SELECT B.[FId],B.[FCron],B.[FMsgTitle],B.[FMsgContent],B.[FAddress],B.[FPassport],B.[FPassportType]" +
	" FROM MConfig A" +
	" INNER JOIN CrmDzXfTestTaskConfig B ON A.[FId] = B.[FId]"

const SqlGetCrmDzXfTestTaskConfigById = "" +
	"SELECT B.[FId],B.[FCron],B.[FMsgTitle],B.[FMsgContent],B.[FAddress],B.[FPassport],B.[FPassportType]" +
	" FROM MConfig A" +
	" INNER JOIN CrmDzXfTestTaskConfig B ON A.[FId] = B.[FId]" +
	" WHERE B.FId = ?"

type CrmDzXfTestTaskConfig struct {
}

type CrmDzXfTestTaskConfigData struct {
	FId           string
	FCron         string
	FMsgTitle     string
	FMsgContent   string
	FAddress      string
	FPassport     string
	FPassportType int
}

func (config *CrmDzXfTestTaskConfigData) IsEqual(c *CrmDzXfTestTaskConfigData) bool {
	if config.FId != c.FId ||
		config.FCron != c.FCron ||
		config.FMsgTitle != c.FMsgTitle ||
		config.FMsgContent != c.FMsgContent ||
		config.FAddress != c.FAddress ||
		config.FPassport != c.FPassport ||
		config.FPassportType != c.FPassportType {
		return false
	}
	return true
}

func (tc *CrmDzXfTestTaskConfig) GetCrmDzTestTaskConfigList() ([]*CrmDzXfTestTaskConfigData, error) {
	rows, err := comm.getRowsBySQL(SqlGetCrmDzXfTestTaskConfig)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return tc.getCrmDzXfTestTaskConfigListByRows(rows)
}

func (tc *CrmDzXfTestTaskConfig) GetCrmDzTestTaskConfig(id string) ([]*CrmDzXfTestTaskConfigData, error) {
	rows, err := comm.getRowsBySQL(SqlGetCrmDzXfTestTaskConfigById, id)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return tc.getCrmDzXfTestTaskConfigListByRows(rows)
}

func (tc *CrmDzXfTestTaskConfig) getCrmDzXfTestTaskConfigListByRows(rows *sql.Rows) ([]*CrmDzXfTestTaskConfigData, error) {
	defer func() {
		_ = rows.Close()
	}()
	var fId, fCron, fMsgTitle, fMsgContent, fAddress, fPassport string
	var fPassportType int
	resultList := make([]*CrmDzXfTestTaskConfigData, 0)
	var err error
	for rows.Next() {
		err = rows.Scan(&fId, &fCron, &fMsgTitle, &fMsgContent, &fAddress, &fPassport, &fPassportType)
		if err != nil {
			break
		}
		config := CrmDzXfTestTaskConfigData{
			FId:           fId,
			FCron:         fCron,
			FMsgTitle:     fMsgTitle,
			FMsgContent:   fMsgContent,
			FAddress:      fAddress,
			FPassport:     fPassport,
			FPassportType: fPassportType,
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
