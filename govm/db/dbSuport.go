package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"../common"
	_ "github.com/go-sql-driver/mysql"
)

const (
	CONSTR string = "mysql:123456@/test?charset=utf8"
)

type MysqlSuport struct {
	State      int
	conn       *sql.DB
	noMap      map[int]int
	tableNames map[string]int
	maxNoId    int
	maxFaceId  int
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
保存Faces
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
			fmt.Printf("noId:%d\n", this.GetIdbyNo(face.Code))
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
				int(face.TurnoverRate*100),
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

/*
获取faces
*/
func (this *MysqlSuport) GetCodeFaces(date string, no int) ([]*common.CodeFace, bool) {
	faces := []*common.CodeFace{}
	if date == "" && no < 0 {
		fmt.Println("GetCodeFaces 缺少查询条件")
		return faces, false
	}

	sqlStr := "select _date,no_id,_min,_max,_change,lastprice,startprice,dde_b,dde_s,face,volume,turnoverRate,turnover,state,per from codeface where"
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

	return this.getFaceFromRows(rows)
}

func (this *MysqlSuport) SplitDB() {
	fs, _ := this.GetCodeFaces("", 1912261)
	dcount := len(fs)
	for i, fc := range fs {
		y, _ := strconv.Atoi(strings.Split(fc.Date, "-")[0])
		if y > 2018 {
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
	sql := fmt.Sprintf("delete from time_price where face_id in(select c.face_id as face_id from tbl_codes as a inner join codeface as b on a.id=b.no_id inner join time_price as c on b.id=c.face_id where b._date='%s');", date)
	r, err := this.conn.Exec(sql)
	if err != nil {
		fmt.Printf("删除失败:%s\n", date)
		fmt.Println(err)
		return
	}

	delNums, _ := r.RowsAffected()
	fmt.Println(delNums)
}

func (this *MysqlSuport) getFaceFromRows(rows *sql.Rows) ([]*common.CodeFace, bool) {
	faces := []*common.CodeFace{}
	var noId, minValue, maxValue, change, startPrice, lastPrice, turnover, turnoverRate, percent int
	for rows.Next() {
		face := &(common.CodeFace{})
		err := rows.Scan(&face.Date,
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
func (this *MysqlSuport) SaveTimePrices(prices []*common.CodePrice) bool {
	return true
}

/*
获取 TimePrice
*/
func (this *MysqlSuport) GetTimePrice(date string, no int) ([]*common.CodePrice, bool) {
	prices := []*common.CodePrice{}
	return prices, true
}
