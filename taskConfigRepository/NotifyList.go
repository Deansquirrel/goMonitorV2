package taskConfigRepository

import (
	"database/sql"
	"github.com/Deansquirrel/goToolCommon"
	log "github.com/Deansquirrel/goToolLog"
	"strings"
)

const SqlGetNotifyList = "" +
	"SELECT [DingTalkRobotId] " +
	"FROM NotifyList " +
	"WHERE [TaskId] = ? or TaskId = '-1'"

type NotifyList struct {
}

type NotifyListData struct {
	DingTalkRobot []string
}

func (nc *NotifyList) GetNotifyList(id string) (*NotifyListData, error) {
	rows, err := comm.getRowsBySQL(SqlGetNotifyList, id)
	if err != nil {
		return nil, err
	}
	return nc.getNotifyListByRows(rows)
}

func (nc *NotifyList) getNotifyListByRows(rows *sql.Rows) (*NotifyListData, error) {
	defer func() {
		_ = rows.Close()
	}()
	var dingTalkRobot string
	dingTalkRobotList := make([]string, 0)
	var err error
	for rows.Next() {
		err = rows.Scan(&dingTalkRobot)
		if err != nil {
			break
		}
		list := strings.Split(dingTalkRobot, ",")
		list = goToolCommon.ClearBlock(list)
		for _, s := range list {
			dingTalkRobotList = append(dingTalkRobotList, s)
		}
	}
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if rows.Err() != nil {
		log.Error(rows.Err().Error())
		return nil, rows.Err()
	}
	dingTalkRobotList = goToolCommon.ClearRepeat(dingTalkRobotList)
	return &NotifyListData{
		DingTalkRobot: dingTalkRobotList,
	}, nil
}
