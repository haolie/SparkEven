package DBSupport

import (
	"database/sql"
	"fmt"
	"sync"

	"SparkEven/govm/Src/Common/Def"
	"SparkEven/govm/Src/Config"
	_ "github.com/go-sql-driver/mysql"
)

var (
	conStr string
	locker sync.Mutex
)

func CreateConn() (conn *sql.DB, errStr string) {
	conn, err := sql.Open("mysql", getConStr())
	if err != nil {
		errStr = fmt.Sprintf("CreateConn open db fail,err:%v", err)
		return
	}

	return
}

func getConStr() (connStr string) {
	if conStr != "" {
		return conStr
	}

	locker.Lock()
	defer locker.Unlock()

	if conStr != "" {
		return conStr
	}

	//CONSTR       string = "haolie:123456@/vmpark?charset=utf8"
	user, _ := Config.GetValue[string](Def.Config_Db_User)
	pw, _ := Config.GetValue[string](Def.Config_Db_Pw)
	dbName, _ := Config.GetValue[string](Def.Config_Db_Name)

	conStr = fmt.Sprintf("%s:%s@/%s?charset=utf8", user, pw, dbName)
	connStr = conStr
	return
}
