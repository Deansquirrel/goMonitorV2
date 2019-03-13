package taskConfigRepository

import (
	"database/sql"
	log "github.com/Deansquirrel/goToolLog"
)

const SqlGetHealthConfig = "" +
	"SELECT B.FId,B.FCron,B.FMsgTitle,B.FMsgContent" +
	" From MConfig A" +
	" INNER JOIN HealthTaskConfig B on A.FId = B.FId"

const SqlGetHealthConfigById = "" +
	"SELECT B.FId,B.FCron,B.FMsgTitle,B.FMsgContent" +
	" From MConfig A" +
	" INNER JOIN HealthTaskConfig B on A.FId = B.FId" +
	" WHERE FId = ?"

type HealthTaskConfig struct {
}

type HealthTaskConfigData struct {
	FId         string
	FCron       string
	FMsgTitle   string
	FMsgContent string
}

func (config *HealthTaskConfigData) IsEqual(c *HealthTaskConfigData) bool {
	if config.FId != c.FId ||
		config.FCron != c.FCron ||
		config.FMsgTitle != c.FMsgTitle ||
		config.FMsgContent != c.FMsgContent {
		return false
	}
	return true
}

func (htc *HealthTaskConfig) GetHealthConfigList() ([]*HealthTaskConfigData, error) {
	rows, err := comm.getRowsBySQL(SqlGetHealthConfig)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return htc.getHealthConfigListByRows(rows)
}

func (htc *HealthTaskConfig) GetHealthConfig(id string) ([]*HealthTaskConfigData, error) {
	rows, err := comm.getRowsBySQL(SqlGetHealthConfigById)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return htc.getHealthConfigListByRows(rows)
}

func (htc *HealthTaskConfig) getHealthConfigListByRows(rows *sql.Rows) ([]*HealthTaskConfigData, error) {
	defer func() {
		_ = rows.Close()
	}()
	var fId, fCron, fMsgTitle, fMsgContent string
	resultList := make([]*HealthTaskConfigData, 0)
	var err error
	for rows.Next() {
		err = rows.Scan(&fId, &fCron, &fMsgTitle, &fMsgContent)
		if err != nil {
			break
		}
		config := HealthTaskConfigData{
			FId:         fId,
			FCron:       fCron,
			FMsgTitle:   fMsgTitle,
			FMsgContent: fMsgContent,
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
