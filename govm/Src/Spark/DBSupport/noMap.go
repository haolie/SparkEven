package DBSupport

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

var (
	noMap      map[int]int
	tableNames map[string]int
	rwLocker   = &sync.RWMutex{}
	// 最大 id
	maxFaceId int
)

func InitSupport(conn *sql.DB) (errStr string) {
	errStr = initNoMap(conn)
	if errStr != "" {
		return
	}

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
			if id > maxFaceId {
				maxFaceId = id
			}
		}
	}

	return
}

func getPriceTableName(date string, conn *sql.DB) (errStr string) {
	if tableNames == nil {

		rows, err := conn.Query("show tables")
		if err != nil {
			return fmt.Sprintf("getPriceTableName faild err:%v", err)
		}

		tableNames = map[string]int{}
		name := ""
		for rows.Next() {
			err = rows.Scan(&name)
			if err != nil {
				return fmt.Sprintf("getPriceTableName faild err:%v", err)
			}

			tableNames[name] = 1
		}
	}

	strs := strings.Split(date, "-")
	m, _ := strconv.Atoi(strs[1])
	name := strs[0]
	if m > 6 {
		name += "_2"
	} else {
		name += "_1"
	}
	name = "timePrice" + name

	if _, k := tableNames[name]; !k {
		addPriceTable(name, conn)
	}

	return
}

func addPriceTable(name string, conn *sql.DB) (errStr string) {
	sqlStr := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(face_id MEDIUMINT UNSIGNED NOT NULL,time SMALLINT UNSIGNED NOT NULL,price SMALLINT,trade_type TINYINT,volume MEDIUMINT UNSIGNED,primary key (face_id,time));", name)
	_, err := conn.Exec(sqlStr)
	if err != nil {
		return fmt.Sprintf("addPriceTable faild err:%v", err)
	}

	return ""
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

	temp := maxFaceId + 1
	ext, _ := conn.Prepare("insert into tbl_codes(id,_no,state) values(?,?,0)")
	_, err := ext.Exec(temp, no)
	if err != nil {
		errStr = fmt.Sprintf("addNewNo faild err:%v", err)
		return
	}

	noMap[no] = temp
	maxFaceId = temp
	id = temp
	return
}
