package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"SparkEven/govm/common"
	"SparkEven/govm/common/csv"
	"SparkEven/govm/db"
	"SparkEven/govm/nohelper"
)

var (
	downProgressCount int             = 2
	ds                *db.MysqlSuport = db.Create()
	chlDataGet        chan *map[int]*common.FaceEx
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
	chlDataGet = make(chan *map[int]*common.FaceEx, 100)
	go SaveDataToDb()
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
	dates := nohelper.GetDatesFromWeb("2020-03-21")
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

type getPriceWork struct {
	face      *common.FaceEx
	prices    []*common.CodePrice
	chlOutput chan *getPriceWork
}

func (this *getPriceWork) Start() {
	ok := false
	filePath := common.GetFilePath(this.face.Date, this.face.Code)
	this.prices, ok = GetTimePriceFromFile(filePath)
	if ok {
		this.face.FileState = 5
	}
	this.Finished()
}

func (this *getPriceWork) Finished() {
	this.chlOutput <- this
}

func SaveDataToDb() {
	//station:=common.CreateDefaultStation

	for {
		f, ok := <-chlDataGet
		if !ok {
			break
		}
		faces := *f
		station := common.CreateDefaultStation(4, len(faces))
		exitChl := make(chan *getPriceWork, len(faces))

		dbStatio := common.CreateDefaultStation(8, len(faces))
		go station.Start()
		go dbStatio.Start()
		for _, face := range faces {
			wf := &getPriceWork{
				face,
				nil,
				exitChl}

			if face.FileState == 3 {
				station.AddItem(wf)
			}
		}

		var totalCount = 0
		var tempTime *time.Time
		for i := 0; i < len(faces); i++ {
			fw := <-exitChl
			//	fmt.Println("文件数据：", time.Since(tempTime))
			if tempTime == nil {
				t := time.Now()
				tempTime = &t
			}

			face := fw.face
			dbStatio.AddItem(&dataSaveWork{
				fw.face,
				fw.prices})
			if face.FileState == 5 {
				totalCount += len(fw.prices)

				//t := time.Now()
				//ok = ds.SaveTimePrices(&face.CodeFace, fw.prices)
				//	fmt.Println("耗时：", time.Since(t))
				// if !ok {
				// 	fmt.Println("数据库保存失败！")
				// 	continue
				// }
			}
		}
		dbStatio.WaitEnd()
		fmt.Println("文件数据：", time.Since(*tempTime))
		fmt.Println("totalCount：", totalCount)
		fmt.Println("完成数据文件解析：")
		break

	}
}

type dataSaveWork struct {
	face   *common.FaceEx
	prices []*common.CodePrice
}

func (this *dataSaveWork) Start() {
	face := this.face
	if face.FileState == 5 {
		ok := ds.SaveTimePrices(&face.CodeFace, this.prices)
		//	fmt.Println("耗时：", time.Since(t))
		if !ok {
			fmt.Println("数据库保存失败！")
		}
	}

}

func (this *dataSaveWork) Finished() {
}

func SaveDataToDb_1() {
	for {
		f, ok := <-chlDataGet
		if !ok {
			break
		}
		faces := *f
		for _, face := range faces {
			if face.FileState == 3 {
				filePath := common.GetFilePath(face.Date, face.Code)
				prices, ok := GetTimePriceFromFile(filePath)
				if !ok {
					fmt.Println("文件数据获取失败")
					continue
				}

				ok = ds.SaveTimePrices(&face.CodeFace, prices)
				if ok {
					fmt.Println("数据库保存失败！")
					continue
				}

			}
		}
		// d := "ces"
		// fmt.Println("开始数据获取：", d, len(*c))
	}
}

func DownDateFiles(date string) {
	faces, _ := nohelper.GetCodeFaces(date)
	fmt.Println(len(faces))
	count := 0
	for _, f := range faces {
		if common.CheckFile(common.GetFilePath(f.Date, f.Code)) {
			f.FileState = 3
		} else {
			count++
		}
	}

	// if count == 0 {
	// 	// ds.SetFaceState(date, common.GCode)
	// 	return
	// }

	if count > 0 {
		var started int = 0
		pcount := common.Min(len(faces), downProgressCount)
		var exitChls = make(chan bool, pcount)
		//var inputchls = make([]chan *common.FaceEx, downProgressCount)
		inputchls := make(chan *common.FaceEx)
		for i := 0; i < pcount; i++ {
			go func(exitchl chan bool, index int) {
				for c := range inputchls {
					DownFile(c.Date, c.Code)
					c.FileState = 3
					fmt.Printf("%d:%d\n", index, c.Code)
				}
				exitchl <- true
			}(exitChls, i)
		}

		for _, face := range faces {
			if face.FileState == 0 {
				inputchls <- face
				face.FileState = 2
				started++
				fmt.Printf("%s started %d/%d\n", date, started, count)
			}
		}
		close(inputchls)
		for i := 0; i < pcount; i++ {
			<-exitChls
		}
		close(exitChls)
	}
	chlDataGet <- &faces
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

	success := common.HttpDown(url, filePath)
	// CheckFile(filePath)
	return success
}

func CheckFile(path string) bool {
	f1, _ := os.Stat(path)
	fmt.Println(f1.Size())
	return f1.Size() > 0
}

func GetTimePriceFromFile(filePath string) ([]*common.CodePrice, bool) {
	fInfo, err := os.Stat(filePath)
	var prices []*common.CodePrice = nil

	if err != nil || fInfo.IsDir() {
		fmt.Println("路径出差")
		return prices, false
	}
	rows, success := csv.GetRowsFromFile(filePath)

	if !success {
		return prices, false
	}
	//fmt.Println(rows)
	prices = make([]*common.CodePrice, 0, len(rows)-1)
	for i, row := range rows {
		if i == 0 {
			continue
		}
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
		tempLen := len(prices)
		if tempLen > 0 {
			last := prices[tempLen-1]
			if last.Time == price.Time {

				// Price     int
				// TradeType int
				// Volume    int

				last.Price = (last.Price*last.Volume + price.Volume*price.Price) / (last.Volume + price.Volume)
				last.TradeType = price.TradeType
				last.Volume = last.Volume + price.Volume
				continue
			}
		}

		prices = append(prices, price)
		//prices[i-1] = price
	}

	// for _, r := range prices {
	// 	fmt.Println(*r)
	// }

	return prices, true
}
