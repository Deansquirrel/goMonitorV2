package taskConfigRepository

import (
	"database/sql"
	log "github.com/Deansquirrel/goToolLog"
)

//const SqlGetIntTaskDConfig = "" +
//	"SELECT [FID],[FMsgSearch] " +
//	"FROM [IntTaskDConfig]"

const SqlGetIntTaskDConfigById = "" +
	"SELECT [FID],[FMsgSearch] " +
	"FROM [IntTaskDConfig] " +
	"WHERE [FId]=?"

type IntTaskDConfig struct {
}

type intTaskDConfigData struct {
	FId        string
	FMsgSearch string
}

//func (itc *IntTaskDConfig) GetIntTaskDConfigList() ([]intTaskDConfigData, error) {
//	rows, err := comm.getRowsBySQL(SqlGetIntTaskDConfig)
//	if err != nil {
//		return nil, err
//	}
//	return itc.getIntTaskDConfigListByRows(rows)
//}

func (itc *IntTaskDConfig) GetIntTaskDConfig(id string) ([]*intTaskDConfigData, error) {
	rows, err := comm.getRowsBySQL(SqlGetIntTaskDConfigById, id)
	if err != nil {
		return nil, err
	}
	return itc.getIntTaskDConfigListByRows(rows)
}

func (itc *IntTaskDConfig) getIntTaskDConfigListByRows(rows *sql.Rows) ([]*intTaskDConfigData, error) {
	defer func() {
		_ = rows.Close()
	}()
	var fId, fMsgSearch sql.NullString
	resultList := make([]*intTaskDConfigData, 0)
	var err error
	for rows.Next() {
		err := rows.Scan(&fId, &fMsgSearch)
		if err != nil {
			break
		}
		config := intTaskDConfigData{}
		config.FId = "Null"
		if fId.Valid {
			config.FId = fId.String
		}
		config.FMsgSearch = "Null"
		if fMsgSearch.Valid {
			config.FMsgSearch = fMsgSearch.String
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
