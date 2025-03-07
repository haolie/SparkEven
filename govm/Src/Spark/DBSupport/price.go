package DBSupport

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"SparkEven/govm/Src/Common/Def"
	"SparkEven/govm/Src/Common/Tools"
	"SparkEven/govm/Src/Model"
)

func addPriceTable(name string, conn *sql.DB) (errStr Model.Err) {
	sqlStr := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(face_id MEDIUMINT UNSIGNED NOT NULL,time SMALLINT UNSIGNED NOT NULL,price SMALLINT NOT NULL,trade_type TINYINT,volume MEDIUMINT UNSIGNED,primary key (face_id,time));", name)
	_, err := conn.Exec(sqlStr)
	if err != nil {
		return Model.Err(fmt.Sprintf("addPriceTable faild err:%v", err))
	}

	tableNamesLock.Lock()
	defer tableNamesLock.Unlock()
	tableNamesMap[name] = struct{}{}

	return
}

/*
初始化价格表名
*/
func initPriceTableName(date string, conn *sql.DB) (errStr string) {
	if tableNamesMap == nil {
		rows, err := conn.Query("show tables")
		if err != nil {
			return fmt.Sprintf("initPriceTableName faild err:%v", err)
		}

		tableNamesMap = make(map[string]struct{}, 8)
		tempName := ""
		for rows.Next() {
			err = rows.Scan(&tempName)
			if err != nil {
				return fmt.Sprintf("initPriceTableName faild err:%v", err)
			}

			tableNamesMap[tempName] = struct{}{}
		}
	}

	name := createTableName(date)
	if _, k := tableNamesMap[name]; !k {
		addPriceTable(name, conn)
	}

	return
}

/*
创建表名
*/
func createTableName(date string) string {
	strList := strings.Split(date, "-")
	m, _ := strconv.Atoi(strList[1])
	name := "timePrice" + strList[0]
	if m > 6 {
		name += "_2"
	} else {
		name += "_1"
	}

	return name
}

/*
返回表名
*/
func getOrAddTableName(conn *sql.DB, date string) (tableName string, errStr Model.Err) {
	tableName = createTableName(date)
	isExists := func() bool {
		tableNamesLock.RLock()
		defer tableNamesLock.RUnlock()

		_, e := tableNamesMap[tableName]
		return e
	}

	if isExists() {
		return
	}

	errStr = addPriceTable(tableName, conn)
	if errStr.Exists() {
		return
	}

	tableNamesLock.Lock()
	defer tableNamesLock.Unlock()
	tableNamesMap[tableName] = struct{}{}

	return
}

/*
保存 TimePrices
*/
func SaveTimePrices(cconn *sql.DB, face *Model.CodeFace, prices []*Model.CodePrice) (errStr Model.Err) {
	tableName, errStr := getOrAddTableName(cconn, face.Date)
	if errStr.Exists() {
		return
	}

	sqlStr := fmt.Sprintf("replace into %s(face_id,time,price,trade_type,volume) VALUES", tableName)
	i, count := 0, 2000
	timeset, _ := Tools.GetSecondsFromStr("09:00:00")
	allCount := len(prices)

	//conn, err := db.Begin()
	sqlCtx, err := cconn.Begin()
	if err != nil {
		return Model.Err(fmt.Sprintf("saveTimePrices faild err:%v", err))
	}

	for {
		index := i + count
		if index > allCount {
			index = allCount
		}

		list := prices[i:index]
		valueStr := ""
		if len(list) == 0 {
			panic("错误")
		}
		for _, cp := range list {
			if len(valueStr) > 0 {
				valueStr += ","
			}

			valueStr = fmt.Sprintf("%s(%d,%d,%d,%d,%d)", valueStr,
				face.ID,
				cp.Time-timeset,
				cp.Price-int(face.StartPrice*100),
				cp.TradeType,
				Tools.Min(cp.Volume, 16777215))

		}

		_, err := sqlCtx.Exec(sqlStr + valueStr + ";")
		if err != nil {
			fmt.Println(sqlStr + valueStr + ";")
			fmt.Println("SaveTimePrices 失败")
			fmt.Println(err)
			sqlCtx.Rollback()
			panic("SaveTimePrices 失败")
		}

		if index == allCount {
			break
		}

		i = index
	}
	sqlStr = fmt.Sprintf("update codeface set state=1 where id=%d;", face.ID)

	_, errr := sqlCtx.Exec(sqlStr)
	if errr != nil {
		fmt.Println(sqlStr)
		fmt.Println("SaveTimePrices 失败")
		fmt.Println(errr)
		sqlCtx.Rollback()
		panic("SaveTimePrices 失败")
	}
	sqlCtx.Commit()
	return
}

func searchPrice(conn *sql.DB, date string, code int) (prices []*Model.CodePrice, errStr Model.Err) {
	face, errStr := getDateCodeFace(conn, code, date)
	if errStr.Exists() {
		return
	}

	if face == nil {
		return
	}

	tbName, errStr := getOrAddTableName(conn, date)
	if errStr.Exists() {
		return
	}

	sqlStr := fmt.Sprintf("SELECT `time`,`price`,`trade_type`,`volume` FROM `%s` WHERE face_id=%d;", tbName, face.ID)
	rows, err := conn.Query(sqlStr)
	if err != nil {
		errStr = Model.Err(fmt.Sprintf(" searchPrice 失败%v", err))
		return
	}

	prices = make([]*Model.CodePrice, 0, 1024)
	for rows.Next() {
		item := &Model.CodePrice{}
		err = rows.Scan(&item.Time, &item.Price, &item.TradeType, &item.Volume)
		if err != nil {
			errStr = Model.Err(fmt.Sprintf(" searchPrice rows.Scan 失败%v", err))
			return
		}

		item.Time += Def.StartTick
		item.Price += int(face.StartPrice * 100)
		prices = append(prices, item)
	}

	return
}
