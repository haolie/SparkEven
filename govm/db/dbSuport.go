package db

import (
	"fmt"
	"strconv"
	"strings"
	"database/sql"

	"../common"
	_ "github.com/go-sql-driver/mysql"
)

const (
	CONSTR string = "mysql@/test?charset=utf8"
)

var (
	noMap      map[int]int = nil //map [no] id
	tableNames map[string]int

	conn       *sql.DB
)

type MysqlSuport struct{
	State int =0
	conn  *sql.DB
	noMap map[int]int = nil
	tableNames map[string]int =nil
	maxNoId int =0
}

// *DB
/*
初始化 编号字典
*/
func (this *MysqlSuport) initNoMap() bool {
	if this.noMap == nil {
		this.noMap = map[int]int{}
		rows,err:= this.conn.Query("SELECT _no,id from tbl_codes")
		if err!=nil{
			retrun false
		}
		var no,id int
		for rows.Next(){
			rows.Scan(&no,&id)
			this.noMap[no]=id
			if id>this.maxNoId{
				this.maxNoid=id
			}
		}
	}

	return true
}

func (this *MysqlSuport) connectDB() bool {
	db, err := sql.Open("mysql", constr)
	if err != nil {
		fmt.Println(err)
		return false
	}

	this.conn=db
	return true
}

/*
根据 编号 获取 id
*/
func (this *MysqlSuport) GetIdbyNo(no int) int {
	this.initNoMap()
	if id, ok := this.noMap[no]; ok {
		return this.noMap[no]
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
func (this *) GetNoById(id int) int {
	for no, d := range noMap {
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
    temp:= this.maxNoId+1
	ext, _ := db.Prepare("insert into tbl_codes(id,_no,state) values(?,?,0)")
	_,err:= ext.Exec(temp,no)
	if err==null{
		return -1
	}
	
	this.maxNoId=temp
	return temp
}

/*
获取
*/
func (this *MysqlSuport) GetPriceTableName(date string) string {
	if this.tableNames == nil {
       this.addPriceTable(date)
	}

	strs := strings.Split(date, "-")
	m := strconv.Atoi(strs[1])
	name := strs[0]
	if m > 6 {
		name += "_2"
	} else {
		name += "_1"
	}

	if n, k := this.tableNames[name]; !k {
		this.addPriceTable(name)
	}
	return name
}

func (this *MysqlSuport) addPriceTable(name string) bool {

	//ret2, _ := this.conn.Exec("update product set name= '000' where id > ?", 1)
	return true
}

/*
保存Faces
*/
func (this *MysqlSuport) SaveCodeFaces(faces []*common.CodeFace) bool {
	return true

}

/*
获取faces
*/
func (this *MysqlSuport) GetCodeFaces(date string, no int) ([]*common.CodeFace, bool) {
	faces := []*common.CodeFace{}
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
