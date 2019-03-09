package taskConfigRepository

import (
	"database/sql"
	"fmt"
	log "github.com/Deansquirrel/goToolLog"
)

const SqlGetDingTalkRobot = "" +
	"SELECT [FId],[FWebHookKey],[FAtMobiles],[FIsAtAll] " +
	"FROM [DingTalkRobot]"

const SqlGetDingTalkRobotById = "" +
	"SELECT [FId],[FWebHookKey],[FAtMobiles],[FIsAtAll] " +
	"FROM [DingTalkRobot]" +
	"WHERE [FId]=?"

const SqlGetDingTalkRobotByIdList = "" +
	"SELECT [FId],[FWebHookKey],[FAtMobiles],[FIsAtAll] " +
	"FROM [DingTalkRobot]" +
	"WHERE [FId] in (%s)"

type DingTalkRobot struct {
}

type DingTalkRobotData struct {
	FId         string
	FWebHookKey string
	FAtMobiles  string
	FIsAtAll    int
}

func (dt *DingTalkRobot) GetDingTalkRobotList() ([]*DingTalkRobotData, error) {
	rows, err := comm.getRowsBySQL(SqlGetDingTalkRobot)
	if err != nil {
		return nil, err
	}
	return dt.getDingTalkRobotByRows(rows)
}

func (dt *DingTalkRobot) GetDingTalkRobotByList(idList []string) ([]*DingTalkRobotData, error) {
	rows, err := comm.getRowsBySQL(dt.getSqlGetDingTalkRobotByIdList(len(idList)), idList)
	if err != nil {
		return nil, err
	}
	return dt.getDingTalkRobotByRows(rows)
}

func (dt *DingTalkRobot) GetDingTalkRobot(id string) ([]*DingTalkRobotData, error) {
	rows, err := comm.getRowsBySQL(SqlGetDingTalkRobotById, id)
	if err != nil {
		return nil, err
	}
	return dt.getDingTalkRobotByRows(rows)
}

func (dt *DingTalkRobot) getDingTalkRobotByRows(rows *sql.Rows) ([]*DingTalkRobotData, error) {
	defer func() {
		_ = rows.Close()
	}()
	var fId, fWebHookKey, fAtMobiles string
	var fIsAtAll int
	resultList := make([]*DingTalkRobotData, 0)
	var err error
	for rows.Next() {
		err = rows.Scan(&fId, &fWebHookKey, &fAtMobiles, &fIsAtAll)
		if err != nil {
			break
		}
		config := DingTalkRobotData{
			FId:         fId,
			FWebHookKey: fWebHookKey,
			FAtMobiles:  fAtMobiles,
			FIsAtAll:    fIsAtAll,
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

func (dt *DingTalkRobot) getSqlGetDingTalkRobotByIdList(num int) string {
	appS := "?"
	for i := 1; i < num; i++ {
		appS = appS + ",?"
	}
	return fmt.Sprintf(SqlGetDingTalkRobotByIdList, appS)
}
