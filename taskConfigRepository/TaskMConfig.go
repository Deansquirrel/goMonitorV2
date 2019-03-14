package taskConfigRepository

import (
	"database/sql"
	log "github.com/Deansquirrel/goToolLog"
)

const SqlGetTaskMConfig = "" +
	"SELECT [FID],[FTitle],[FRemark] " +
	"FROM [MConfig]"
const SqlGetTaskMConfigById = "" +
	"SELECT [FID],[FTitle],[FRemark] " +
	"FROM [MConfig] WHERE [FID] = ?"

type TaskMConfig struct {
}

type TaskMConfigData struct {
	FId     string
	FTitle  string
	FRemark string
}

//func NewTaskMConfig(title string, remark string) *taskMConfigData {
//	return &taskMConfigData{
//		FId:     strings.ToUpper(goToolCommon.Guid()),
//		FTitle:  title,
//		FRemark: remark,
//	}
//}

func (tmc *TaskMConfig) GetMConfigList() ([]*TaskMConfigData, error) {
	rows, err := comm.getRowsBySQL(SqlGetTaskMConfig)
	if err != nil {
		return nil, err
	}
	return tmc.getMConfigListByRows(rows)
}

func (tmc *TaskMConfig) GetMConfig(id string) ([]*TaskMConfigData, error) {
	rows, err := comm.getRowsBySQL(SqlGetTaskMConfigById, id)
	if err != nil {
		return nil, err
	}
	return tmc.getMConfigListByRows(rows)
}

func (tmc *TaskMConfig) getMConfigListByRows(rows *sql.Rows) ([]*TaskMConfigData, error) {
	defer func() {
		_ = rows.Close()
	}()
	var fId, fTitle, fRemark string
	resultList := make([]*TaskMConfigData, 0)
	var err error
	for rows.Next() {
		err = rows.Scan(&fId, &fTitle, &fRemark)
		if err != nil {
			return nil, err
		}
		config := TaskMConfigData{
			FId:     fId,
			FTitle:  fTitle,
			FRemark: fRemark,
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
