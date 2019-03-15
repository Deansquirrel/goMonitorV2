package taskConfigRepository

import (
	"database/sql"
	log "github.com/Deansquirrel/goToolLog"
)

const SqlGetIntTaskConfig = "" +
	"SELECT B.[FId],B.[FServer],B.[FPort],B.[FDbName],B.[FDbUser]," +
	"B.[FDbPwd],B.[FSearch],B.[FCron],B.[FCheckMax],B.[FCheckMin]," +
	"B.[FMsgTitle],B.[FMsgContent]" +
	" FROM [MConfig] A" +
	" INNER JOIN [IntTaskConfig] B ON A.[FId] = B.[FId]"

const SqlGetIntTaskConfigById = "" +
	"SELECT B.[FId],B.[FServer],B.[FPort],B.[FDbName],B.[FDbUser]," +
	"B.[FDbPwd],B.[FSearch],B.[FCron],B.[FCheckMax],B.[FCheckMin]," +
	"B.[FMsgTitle],B.[FMsgContent]" +
	" FROM [MConfig] A" +
	" INNER JOIN [IntTaskConfig] B ON A.[FId] = B.[FId]" +
	" WHERE B.[FId]=?"

type IntTaskConfig struct {
}

type IntTaskConfigData struct {
	FId         string
	FServer     string
	FPort       int
	FDbName     string
	FDbUser     string
	FDbPwd      string
	FSearch     string
	FCron       string
	FCheckMax   int
	FCheckMin   int
	FMsgTitle   string
	FMsgContent string
}

func (config *IntTaskConfigData) IsEqual(c *IntTaskConfigData) bool {
	if config.FId != c.FId {
		return false
	}
	if config.FServer != c.FServer {
		return false
	}
	if config.FPort != c.FPort {
		return false
	}
	if config.FDbName != c.FDbName {
		return false
	}
	if config.FDbUser != c.FDbUser {
		return false
	}
	if config.FDbPwd != c.FDbPwd {
		return false
	}
	if config.FSearch != c.FSearch {
		return false
	}
	if config.FCron != c.FCron {
		return false
	}
	if config.FCheckMax != c.FCheckMax {
		return false
	}
	if config.FCheckMin != c.FCheckMin {
		return false
	}
	if config.FMsgTitle != c.FMsgTitle {
		return false
	}
	if config.FMsgContent != c.FMsgContent {
		return false
	}
	return true
}

func (itc *IntTaskConfig) GetIntTaskConfigList() ([]*IntTaskConfigData, error) {
	rows, err := comm.getRowsBySQL(SqlGetIntTaskConfig)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return itc.getIntTaskConfigListByRows(rows)
}

func (itc *IntTaskConfig) GetIntTaskConfig(id string) ([]*IntTaskConfigData, error) {
	rows, err := comm.getRowsBySQL(SqlGetIntTaskConfigById, id)
	if err != nil {
		return nil, err
	}
	return itc.getIntTaskConfigListByRows(rows)
}

func (itc *IntTaskConfig) getIntTaskConfigListByRows(rows *sql.Rows) ([]*IntTaskConfigData, error) {
	defer func() {
		_ = rows.Close()
	}()
	var fId, fServer, fDbName, fDbUser, fDbPwd, fSearch, fCron, fMsgTitle, fMsgContent string
	var fPort, fCheckMax, fCheckMin int
	resultList := make([]*IntTaskConfigData, 0)
	var err error
	for rows.Next() {
		err = rows.Scan(
			&fId, &fServer, &fPort, &fDbName, &fDbUser,
			&fDbPwd, &fSearch, &fCron, &fCheckMax, &fCheckMin,
			&fMsgTitle, &fMsgContent)
		if err != nil {
			break
		}
		config := IntTaskConfigData{
			FId:         fId,
			FServer:     fServer,
			FPort:       fPort,
			FDbName:     fDbName,
			FDbUser:     fDbUser,
			FDbPwd:      fDbPwd,
			FSearch:     fSearch,
			FCron:       fCron,
			FCheckMax:   fCheckMax,
			FCheckMin:   fCheckMin,
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
