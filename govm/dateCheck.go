package main

import (
	"fmt"

	"lyh/SparkEven/govm/db"
)

var dataSupporter *db.MysqlSuport

func init() {
	dataSupporter = new(db.MysqlSuport)
	dataSupporter.Init()
}

func checkDate(date string, no int) {
	faceList, success := dataSupporter.GetCodeFaces(date, no)
	if !success {
		return
	}

	for _, faceItem := range faceList {
		if faceItem.Code == 1912261 {
			continue
		}

		if faceItem.State != 1 {
			continue
		}

		count, err := dataSupporter.GetTimePirceCount(faceItem)
		if err != nil {
			return
		}

		if count == 0 {
			fmt.Printf("dataErr date:%s,no:%d", faceItem.Date, faceItem.Code)
		}
	}
}

func main() {
	checkDate("2021-02-08", -1)
}
