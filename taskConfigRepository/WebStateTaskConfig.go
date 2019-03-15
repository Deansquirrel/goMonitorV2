package taskConfigRepository

import (
	"database/sql"
	log "github.com/Deansquirrel/goToolLog"
)

const SqlGetWebStateTaskConfig = "" +
	"SELECT B.[FId], B.[FUrl], B.[FCron], B.[FMsgTitle], B.[FMsgContent]" +
	" FROM MConfig A" +
	" INNER JOIN WebStateTaskConfig B ON A.FID = B.FId"

const SqlGetWebStateTaskConfigById = "" +
	"SELECT B.[FId], B.[FUrl], B.[FCron], B.[FMsgTitle], B.[FMsgContent]" +
	" FROM MConfig A" +
	" INNER JOIN WebStateTaskConfig B ON A.FID = B.FId" +
	" WHERE B.FId = ?"

type WebStateTaskConfig struct {
}

type WebStateTaskConfigData struct {
	FId         string
	FUrl        string
	FCron       string
	FMsgTitle   string
	FMsgContent string
}

func (config *WebStateTaskConfigData) IsEqual(c *WebStateTaskConfigData) bool {
	if config.FId != c.FId ||
		config.FUrl != c.FUrl ||
		config.FCron != c.FCron ||
		config.FMsgTitle != c.FMsgTitle ||
		config.FMsgContent != c.FMsgContent {
		return false
	}
	return true
}

func (wtc *WebStateTaskConfig) GetWebStateTaskConfigList() ([]*WebStateTaskConfigData, error) {
	rows, err := comm.getRowsBySQL(SqlGetWebStateTaskConfig)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return wtc.getWebStateConfigListByRows(rows)
}

func (wtc *WebStateTaskConfig) GetWebStateTaskConfig(id string) ([]*WebStateTaskConfigData, error) {
	rows, err := comm.getRowsBySQL(SqlGetWebStateTaskConfigById, id)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return wtc.getWebStateConfigListByRows(rows)
}

func (wtc *WebStateTaskConfig) getWebStateConfigListByRows(rows *sql.Rows) ([]*WebStateTaskConfigData, error) {
	defer func() {
		_ = rows.Close()
	}()
	var fId, fUrl, fCron, fMsgTitle, fMsgContent string
	resultList := make([]*WebStateTaskConfigData, 0)
	var err error
	for rows.Next() {
		err = rows.Scan(&fId, &fUrl, &fCron, &fMsgTitle, &fMsgContent)
		if err != nil {
			break
		}
		config := WebStateTaskConfigData{
			FId:         fId,
			FUrl:        fUrl,
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
