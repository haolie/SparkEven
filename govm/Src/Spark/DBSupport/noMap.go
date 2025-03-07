package DBSupport

import (
	"database/sql"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"SparkEven/govm/Src/Common/Log"
	"SparkEven/govm/Src/Common/Tools"
)

var (
	noMap map[int]int

	tableNamesMap  map[string]struct{}
	tableNamesLock sync.RWMutex

	rwLocker = &sync.RWMutex{}
	// 最大NoId(tbl_codes 最大id)
	maxNoId int
	// 最大faceId(codeface 最大Id)
	maxFaceId int32
)

func InitSupport(conn *sql.DB) (errStr string) {
	dateStr := Tools.ToDateStr(time.Now())
	errStr = initNoMap(conn)
	if errStr != "" {
		return
	}

	row := conn.QueryRow("select max(id) as id from codeface")
	row.Scan(&maxFaceId)

	//
	initPriceTableName(dateStr, conn)

	Log.Info(fmt.Sprintf("init db support maxNoId=%d,maxFaceId=%d", maxNoId, maxFaceId))

	return
}

func initNoMap(conn *sql.DB) (errStr string) {
	if noMap == nil {
		noMap = map[int]int{}
		rows, err := conn.Query("SELECT _no,id from tbl_codes")
		if err != nil {
			return fmt.Sprintf("initNoMap faild err:%v", err)
		}

		var no, id int
		for rows.Next() {
			err = rows.Scan(&no, &id)
			if err != nil {
				return fmt.Sprintf("initNoMap faild err:%v", err)
			}

			noMap[no] = id
			if id > maxNoId {
				maxNoId = id
			}
		}
	}

	return
}

/*
根据 编号 获取 id
*/
func GetIdbyNo(no int) (id int, exists bool) {
	rwLocker.RLock()
	defer rwLocker.RUnlock()

	id, exists = noMap[no]

	return
}

func GetOrAddId(no int, conn *sql.DB) (id int, errStr string) {
	id, exists := GetIdbyNo(no)
	if exists {
		return
	}

	id, errStr = addNewNo(no, conn)
	return
}

/*
根据 id获取编号
*/
func GetNoById(id int) (no int, exists bool) {
	for no, d := range noMap {
		if id == d {
			return no, true
		}
	}

	return
}

func addNewNo(no int, conn *sql.DB) (id int, errStr string) {
	rwLocker.Lock()
	defer rwLocker.Unlock()

	id, exists := noMap[no]
	if exists {
		return
	}

	temp := maxNoId + 1
	ext, _ := conn.Prepare("insert into tbl_codes(id,_no,state) values(?,?,0)")
	_, err := ext.Exec(temp, no)
	if err != nil {
		errStr = fmt.Sprintf("addNewNo faild err:%v", err)
		return
	}

	noMap[no] = temp
	maxNoId = temp
	id = temp
	return
}

func getNewFaceId() int32 {
	return atomic.AddInt32(&maxFaceId, 1)
}
