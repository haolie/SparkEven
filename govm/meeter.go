package main

import (
	"fmt"
	"strconv"
	"strings"

	"os"

	"./common"
	"./common/csv"
	"./db"
	"./nohelper"
)

var (
	downProgressCount int             = 2
	ds                *db.MysqlSuport = db.Create()
)

func main() {
	fmt.Println("开始了 ")
	go nohelper.StartListenSession()
	//faces := nohelper.GetNocodesFromWeb("2020-01-20")
	//	DownDateFiles("2020-01-20", faces)
	//fmt.Println(nohelper.GetDatesFromWeb("2019-01-20"))
	//fmt.Println("token:",tk)
	//fmt.Println(common.GetFilePath("2019-09-09",1600293))
	//fmt.Println(common.CodeNumToStr(1600222))
	//fmt.Println(common.CodeNumToStr(1300222))
	//DownFile("2020-01-20", 1603922)
	// fs := []*common.CodeFace{}
	// for _, f := range faces {
	// 	fs = append(fs, &f.CodeFace)
	// }

	//nohelper.GetCodeFaces("2020-01-06")
	StartDataSave()
	//ds.SaveCodeFaces(fs)
	// fs, _ := ds.GetCodeFaces("2020-01-06", 1300782)
	// for _, f := range fs {
	// 	f.Println()
	// 	cps, _ := ds.GetTimePrice(f)
	// 	for _, cp := range cps {
	// 		cp.Print()
	// 	}
	// }
	// GetTimePriceFromFile("E:/VM/common/go/src/SparkEven/govm/datefiles/2020-01-20/2020-01-20_1000006.xls")
	select {}
}

func StartDataSave() bool {
	dates := nohelper.GetDatesFromWeb("2020-02-12")
	for _, date := range dates {
		fmt.Println("…………………………")
		fmt.Println("…………………………")
		fmt.Println("…………………………")
		fmt.Printf("start down：%s\n", date)
		fmt.Println("…………………………")
		DownDateFiles(date)
		fmt.Println("…………………………")
		fmt.Println("…………………………")
		fmt.Println("…………………………")
		fmt.Printf("finish down：%s\n", date)
		fmt.Println("…………………………")
	}
	return true
}

func SaveDataToDb() {

}

func DownDateFiles(date string) {
	faces := nohelper.GetCodeFaces(date)
	fmt.Println(len(faces))
	for k, f := range faces {
		if common.CheckFile(common.GetFilePath(f.Date, f.Code)) {
			delete(faces, k)
		}
	}

	count := len(faces)
	if count == 0 {
		// ds.SetFaceState(date, common.GCode)
		return
	}

	var started int = 0
	pcount := common.Min(len(faces), downProgressCount)
	var exitChls = make(chan bool, pcount)
	//var inputchls = make([]chan *common.FaceEx, downProgressCount)
	inputchls := make(chan *common.FaceEx)
	for i := 0; i < pcount; i++ {
		go func(exitchl chan bool, index int) {
			for c := range inputchls {
				DownFile(c.Date, c.Code)
				fmt.Printf("%d:%d\n", index, c.Code)
			}
			exitchl <- true
		}(exitChls, i)
	}

	for _, face := range faces {
		inputchls <- face
		started++
		fmt.Printf("%s started %d/%d\n", date, started, count)
	}
	close(inputchls)
	for i := 0; i < pcount; i++ {
		<-exitChls
	}
	close(exitChls)
	fmt.Printf("%s  下载完成\n", date)
}

func DownFile(date string, code int) bool {
	filePath := common.GetFilePath(date, code)
	url := "http://stock.gtimg.cn/data/index.php?appn=detail&action=download"
	var codestr string
	if code >= 1600000 {
		codestr = fmt.Sprintf("sh%d", code-1000000)
	} else {
		codestr = "sz" + strconv.Itoa(code)[1:]
	}

	url = fmt.Sprintf("%s&c=%s&d=%s", url, codestr, strings.Replace(date, "-", "", 2))

	return common.HttpDown(url, filePath)
}

func GetTimePriceFromFile(filePath string) {
	fInfo, err := os.Stat(filePath)
	if err != nil || fInfo.IsDir() {
		fmt.Println("路径出差")
		return
	}
	rows, success := csv.GetRowsFromFile(filePath)

	if !success {
		return
	}
	//fmt.Println(rows)
	prices := make([]*common.CodePrice, len(rows))
	for i, row := range rows {
		price := &common.CodePrice{}
		price.Time, _ = common.GetSecondsFromStr(row[0])
		temp, _ := strconv.ParseFloat(row[1], 32)
		price.Price = int(temp * 100)
		price.Volume, _ = strconv.Atoi(row[3])
		if strings.Index(row[5], "买") >= 0 {
			price.TradeType = 1
		} else if strings.Index(row[5], "卖") >= 0 {
			price.TradeType = -1
		}
		prices[i] = price
	}

	for _, r := range prices {
		fmt.Println(*r)
	}

}
