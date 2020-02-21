package main

import (
	"./db"
)

func main() {
	ds := db.MysqlSuport{}
	ds.Init()
	//ds.SaveCodeFaces(fs)
	ds.SplitDB()
	//ds.DeleteDate()
}
