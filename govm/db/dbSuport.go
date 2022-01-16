package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"SparkEven/govm/common"
)

const (
	CONSTR       string = "haolie:123456@/vmpark?charset=utf8"
	Mysql_smUInt int    = 65535
)

type MysqlSuport struct {
	State      int
	conn       *sql.DB
	noMap      map[int]int
	tableNames map[string]int
	maxNoId    int
	maxFaceId  int
}

var instense *MysqlSuport

func Create() *MysqlSuport {
	if instense == nil {
		instense = &MysqlSuport{}
		instense.Init()
	}

	return instense
}

func CloseDefault() {
	if instense != nil {
		instense.Close()
		instense = nil
	}
}

func (this *MysqlSuport) Close() {
	if this.State == 1 {
		this.conn.Close()
		this.State = 0
	}
}

func (this *MysqlSuport) Init() bool {
	if this.State == 0 {
		db, _ := this.connectDB()
		this.initNoMap()

		row := db.QueryRow("select max(id) as id from codeface")
		row.Scan(&this.maxFaceId)

		this.State = 1
	}
	return true
}

// *DB
/*
初始化 编号字典
*/
func (this *MysqlSuport) initNoMap() bool {
	if this.noMap == nil {
		this.noMap = map[int]int{}
		rows, err := this.conn.Query("SELECT _no,id from tbl_codes")
		if err != nil {
			fmt.Println("initNoMap 查询 _no,id 失败")
			fmt.Println(err)
			return false
		}
		var no, id int
		for rows.Next() {
			rows.Scan(&no, &id)
			this.noMap[no] = id
			if id > this.maxNoId {
				this.maxNoId = id
			}
		}
	}

	return true
}

func (this *MysqlSuport) connectDB() (*sql.DB, bool) {
	db, err := sql.Open("mysql", CONSTR)
	if err != nil {
		fmt.Println("connectDB 失败")
		fmt.Println(err)
		return db, false
	}

	this.conn = db
	return db, true
}

/*
根据 编号 获取 id
*/
func (this *MysqlSuport) GetIdbyNo(no int) int {
	this.initNoMap()
	if id, ok := this.noMap[no]; ok {
		return id
	} else {
		id := this.AddNewNo(no)
		if id > 0 {
			this.noMap[no] = id
		}
		return id
	}
}

/*
根据 id获取编号
*/
func (this *MysqlSuport) GetNoById(id int) int {
	for no, d := range this.noMap {
		if id == d {
			return no
		}
	}
	return -1
}

/*
添加新的编号
*/
func (this *MysqlSuport) AddNewNo(no int) int {
	temp := this.maxNoId + 1
	ext, _ := this.conn.Prepare("insert into tbl_codes(id,_no,state) values(?,?,0)")
	_, err := ext.Exec(temp, no)
	if err != nil {
		return -1
	}

	this.maxNoId = temp
	return temp
}

/*
获取
*/
func (this *MysqlSuport) GetPriceTableName(date string) string {
	if this.tableNames == nil {

		rows, err := this.conn.Query("show tables")
		if err != nil {
			fmt.Println("GetPriceTableName 失败")
			fmt.Println(err)
			return ""
		}
		this.tableNames = map[string]int{}
		name := ""
		for rows.Next() {
			rows.Scan(&name)
			this.tableNames[name] = 1
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

	if _, k := this.tableNames[name]; !k {
		this.addPriceTable(name)
	}
	return name
}

func (this *MysqlSuport) addPriceTable(name string) bool {
	sqlStr := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(face_id MEDIUMINT UNSIGNED NOT NULL,time SMALLINT UNSIGNED NOT NULL,price SMALLINT,trade_type TINYINT,volume MEDIUMINT UNSIGNED,primary key (face_id,time));", name)
	_, err := this.conn.Exec(sqlStr)
	if err != nil {
		fmt.Println("addPriceTable 失败")
		fmt.Println(err)
	}
	return err == nil
}

/*
更新交易日状态
*/
func (this *MysqlSuport) UpdateFaceState(date string, code int, state int) {
	faces, ok := this.GetCodeFaces(date, code)
	if !ok {
		return
	}
	if len(faces) > 0 {
		updateStr := fmt.Sprintf("update codeface set state=%d where id=%d;", state, faces[0].ID)
		_, err := this.conn.Exec(updateStr)
		if err != nil {
			fmt.Println("更新交易日期状态失败")
			fmt.Println(err)
		}
	} else {
		face := new(common.CodeFace)
		face.Code = code
		face.Date = date
		face.State = state
		ok = this.SaveCodeFaces([]*common.CodeFace{face})
		if !ok {
			fmt.Println("更新交易日期状态失败")
		}
	}

}

/*
保存Facesgit
*/
func (this *MysqlSuport) SaveCodeFaces(faces []*common.CodeFace) bool {

	i, count := 0, 300
	for {
		index := i + count
		if index > len(faces) {
			index = len(faces)
		}

		array := faces[i:index]
		var str = "INSERT delayed INTO codeface(id,_date,no_id,_min,_max,_change,lastprice,startprice,volume,turnoverRate,turnover,face,dde,dde_b,dde_s,state,per) VALUES"
		var vStr = ""
		for _, face := range array {
			if face.Volume > 16777215 {
				face.Volume = 16777215
			}
			this.maxFaceId++
			face.ID = this.maxFaceId
			turnoverRate := int(face.TurnoverRate * 100)
			if turnoverRate > Mysql_smUInt {
				turnoverRate = Mysql_smUInt
			}
			// fmt.Printf("noId:%d\n", this.GetIdbyNo(face.Code))
			tempStr := fmt.Sprintf("(%d,'%s',%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d)",
				this.maxFaceId,
				face.Date,
				this.GetIdbyNo(face.Code),
				int(face.MinValue*100),
				int(face.MaxValue*100),
				int(face.Change*100),
				int(face.LastPrice*100),
				int(face.StartPrice*100),
				face.Volume,
				turnoverRate,
				int(face.Turnover*100),
				face.Face,
				0,
				face.Dde_b,
				face.Dde_s,
				face.State,
				int(face.Percent*100))

			if len(vStr) > 0 {
				vStr += ","
			}
			vStr += tempStr

		}

		_, err := this.conn.Exec(str + vStr + ";")
		if err != nil {
			fmt.Println("SaveCodeFaces 失败")
			fmt.Println(err)
			return false
		}
		i += count
		if i >= len(faces) {
			break
		}

	}

	return true

}

func (this *MysqlSuport) GetSavedDates(date string) map[string]bool {
	dates := map[string]bool{}

	sql := "select _date from codeface where no_id=0"
	if common.CheckDateStr(date) {
		sql = fmt.Sprintf("%s and _date>'%s'", sql, date)
	}
	sql += ";"
	rows, err := this.conn.Query(sql)
	if err != nil {
		fmt.Println("GetSavedDates 失败")
		fmt.Println(err)
		return dates
	}

	for rows.Next() {
		var str string
		rows.Scan(&str)
		dates[str] = true
	}

	return dates
}

/*
获取faces
*/
func (this *MysqlSuport) GetCodeFaces(date string, no int) ([]*common.CodeFace, bool) {
	faces := []*common.CodeFace{}
	if date == "" && no < 0 {
		fmt.Println("GetCodeFaces 缺少查询条件")
		return faces, false
	}

	sqlStr := "select id,_date,no_id,_min,_max,_change,lastprice,startprice,dde_b,dde_s,face,volume,turnoverRate,turnover,state,per from codeface where"
	filter := ""
	if len(date) > 0 {
		filter = fmt.Sprintf(" _date='%s'", date)
	}

	if no >= 0 {
		if len(filter) > 0 {
			filter += " and"
		}

		filter = fmt.Sprintf("%s no_id=%d", filter, this.GetIdbyNo(no))
	}

	rows, err := this.conn.Query(sqlStr + filter + ";")
	if err != nil {
		fmt.Println("GetCodeFaces 失败")
		fmt.Println(err)
		return faces, false
	}
	// fmt.Println(sqlStr + filter + ";")
	return this.getFaceFromRows(rows)
}

func (this *MysqlSuport) SplitDB() {
	fs, _ := this.GetCodeFaces("", 1912261)
	dcount := len(fs)
	for i, fc := range fs {
		// y, _ := strconv.Atoi(strings.Split(fc.Date, "-")[0])
		// if y > 2018 {
		// 	continue
		// }
		if i < 342 {
			continue
		}
		sql := fmt.Sprintf("select id from codeface where _date='%s' and no_id>0;", fc.Date)
		frows, _ := this.conn.Query(sql)
		//fmt.Println(sql)
		fid := 0
		tbName := this.GetPriceTableName(fc.Date)
		j := 0
		for frows.Next() {
			err := frows.Scan(&fid)
			if err != nil {
				fmt.Println("查询出错！")
				return
			}
			j++
			sql = fmt.Sprintf("replace into %s select * from time_price where face_id=%d;", tbName, fid)
			this.conn.Exec(sql)
			fmt.Printf("正在转存：%d/%d/%d\n", j, i+1, dcount)
		}

		fmt.Printf("正在转存：%d/%d\n", i, dcount)
	}

}

func (this *MysqlSuport) DeleteDate(date string) {
	fs, _ := this.GetCodeFaces("", 1912261)
	dcount := len(fs)
	for i, fc := range fs {

		sql := fmt.Sprintf("select id from codeface where _date='%s' and no_id>0;", fc.Date)
		frows, _ := this.conn.Query(sql)
		//fmt.Println(sql)
		fid := 0
		j := 0
		for frows.Next() {
			err := frows.Scan(&fid)
			if err != nil {
				fmt.Println("查询出错！")
				return
			}
			j++
			sql = fmt.Sprintf("delete from time_price where face_id=%d;", fid)
			this.conn.Exec(sql)
			fmt.Printf("正在删除：%d/%d/%d\n", j, i+1, dcount)
		}

		fmt.Printf("正在删除：%d/%d\n", i, dcount)
	}
}

func (this *MysqlSuport) getFaceFromRows(rows *sql.Rows) ([]*common.CodeFace, bool) {
	faces := []*common.CodeFace{}
	var noId, minValue, maxValue, change, startPrice, lastPrice, turnover, turnoverRate, percent int
	for rows.Next() {
		face := &(common.CodeFace{})
		err := rows.Scan(
			&face.ID,
			&face.Date,
			&noId,
			&minValue,
			&maxValue,
			&change,
			&lastPrice,
			&startPrice,
			&face.Dde_b,
			&face.Dde_s,
			&face.Face,
			&face.Volume,
			&turnover,
			&turnoverRate,
			&face.State,
			&percent)
		if err != nil {
			fmt.Println("getFaceFromRows 失败")
			fmt.Println(err)
			return faces, false
		}

		face.Code = this.GetNoById(noId)
		face.MinValue = float32(minValue) / 100
		face.MaxValue = float32(maxValue) / 100
		face.Change = float32(change) / 100
		face.LastPrice = float32(lastPrice) / 100
		face.StartPrice = float32(startPrice) / 100
		face.Turnover = float32(turnover) / 100
		face.TurnoverRate = float32(turnoverRate) / 100
		face.Percent = float32(percent) / 100
		face.YestPrice = face.LastPrice - face.Change
		faces = append(faces, face)
	}
	return faces, true
}

/*
保存 TimePrices
*/
func (this *MysqlSuport) SaveTimePrices(face *common.CodeFace, prices []*common.CodePrice) bool {
	tableName := this.GetPriceTableName(face.Date)
	sqlStr := fmt.Sprintf("replace into %s(face_id,time,price,trade_type,volume) VALUES", tableName)
	i, count := 0, 2000
	timeset, _ := common.GetSecondsFromStr("09:00:00")
	allCount := len(prices)

	//conn, err := db.Begin()
	conn, err := this.conn.Begin()
	if err != nil {
		fmt.Println("数据库事务开启失败")
		fmt.Println(err)
		return false
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
				common.Min(cp.Volume, 16777215))

		}

		_, err := conn.Exec(sqlStr + valueStr + ";")
		if err != nil {
			fmt.Println(sqlStr + valueStr + ";")
			fmt.Println("SaveTimePrices 失败")
			fmt.Println(err)
			conn.Rollback()
			panic("SaveTimePrices 失败")
		}

		if index == allCount {
			break
		}

		i = index
	}
	sqlStr = fmt.Sprintf("update codeface set state=1 where id=%d;", face.ID)

	_, errr := conn.Exec(sqlStr)
	if errr != nil {
		fmt.Println(sqlStr)
		fmt.Println("SaveTimePrices 失败")
		fmt.Println(errr)
		conn.Rollback()
		panic("SaveTimePrices 失败")
	}
	conn.Commit()
	return true
}

/*
获取 TimePrice
*/
func (this *MysqlSuport) GetTimePrice(face *common.CodeFace) ([]*common.CodePrice, bool) {
	tableName := this.GetPriceTableName(face.Date)
	prices := []*common.CodePrice{}

	sqlStr := fmt.Sprintf("select face_id,time,price,trade_type,a.volume as volume from %s as a join codeface as b on a.face_id=b.id join tbl_codes as c on b.no_id = c.id where c._no=%d",
		tableName, face.Code)

	rows, err := this.conn.Query(sqlStr)

	if err != nil {
		fmt.Println("GetTimePrice 失败")
		fmt.Println(err)
		fmt.Println(sqlStr)
		return prices, false
	}
	timeset, _ := common.GetSecondsFromStr("09:00:00")
	for rows.Next() {
		cp := common.CodePrice{}
		rows.Scan(&cp.FaceId, &cp.Time, &cp.Price, &cp.TradeType, &cp.Volume)
		cp.Price = cp.Price + int(face.StartPrice*100)
		cp.Time += timeset
		prices = append(prices, &cp)
	}

	return prices, true
}

func (this *MysqlSuport) GetTimePirceCount(face *common.CodeFace) (count int, err error) {
	tableName := this.GetPriceTableName(face.Date)

	sqlStr := fmt.Sprintf("select Count(face_id) from %s as a join codeface as b on a.face_id=b.id join tbl_codes as c on b.no_id = c.id where c._no=%d",
		tableName, face.Code)

	rows, err := this.conn.Query(sqlStr)

	if err != nil {
		fmt.Println("GetTimePrice 失败")
		fmt.Println(err)
		fmt.Println(sqlStr)
		return 0, err
	}

	for rows.Next() {
		rows.Scan(&count)
	}

	return
}

func (this *MysqlSuport) SetFaceState(date string, code int, state int) bool {
	if !common.CheckDateStr(date) && code <= 0 {
		return false
	}

	baseSql := fmt.Sprintf("update codeface set state=%d", state)
	filter := ""
	if common.CheckDateStr(date) {
		filter = fmt.Sprintf(" where _date='%s'", date)
	}
	if id := this.GetIdbyNo(code); id > 0 {
		if len(filter) > 0 {
			filter += " and"
		}

		filter = fmt.Sprintf("%s no_id=%d", filter, id)
	}

	if len(filter) == 0 {
		return false
	}

	_, err := this.conn.Exec(baseSql + filter + ";")
	return err == nil
}
