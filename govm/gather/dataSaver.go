package gather

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"lyh/SparkEven/govm/common"
	"lyh/SparkEven/govm/common/csv"
)

type DataSaver struct {
	TaskWorkStation
}

func (this *DataSaver) GetState() int32 {
	if this.chlInput == nil {
		return 1
	}

	v := atomic.LoadInt32(&this.state)
	if v > 0 {
		return v
	}

	return int32(len(this.chlInput))
}

func (this *DataSaver) Start() {
	go func() {
		for {
			faces, ok := <-this.chlInput
			if !ok {
				break
			}
			this.StartDataGet(*faces)
			this.TryFinished()
		}
	}()
}

func (this *DataSaver) TryFinished() {
	if len(this.chlInput) > 0 {
		return
	}

	this.task.EndTask()
}

func (this *DataSaver) StartDataGet(faces map[int]*common.FaceEx) {
	atomic.StoreInt32(&this.state, 1)
	fileStation := common.CreateDefaultStation(4, len(faces))
	go fileStation.Start()
	exitChl := make(chan *getPriceWork, len(faces))

	saveStation := common.CreateDefaultStation(4, len(faces))
	go saveStation.Start()

	var wg sync.WaitGroup
	wg.Add(1)
	inputCount := 0
	go func(chl chan *getPriceWork) {
		for pw := range chl {
			inputCount++
			if pw.face.FileState == 5 || pw.face.DownTimes > 5 {
				saveStation.AddItem(&dataSaveWork{
					pw.face, pw.prices, this.task})
			}
		}
		wg.Done()
	}(exitChl)

	date := ""
	for _, face := range faces {
		wf := &getPriceWork{
			face,
			nil,
			exitChl}
		date = face.Date
		fileStation.AddItem(wf)
	}
	fileStation.EndInput()
	fileStation.WaitEnd()
	fstion, _ := fileStation.(*common.WorkStation)
	fmt.Println(date, "FinishedCount:", fstion.FinishedCount)
	fmt.Println(date, "TotalItemCount:", fstion.TotalItemCount)
	close(exitChl)
	wg.Wait()
	saveStation.EndInput()
	saveStation.WaitEnd()
	sp, _ := saveStation.(*common.WorkStation)
	if sp.StartTime != nil {
		fmt.Println(date, time.Since(*sp.StartTime))
	}

	allSaved := true
	for _, face := range faces {
		allSaved = allSaved && (face.FileState == 1 || face.FileState == 8 || face.DownTimes > 5)
	}
	if allSaved {
		if len(date) > 0 {
			this.task.DB.UpdateFaceState(date, 1912261, 1)
		}
		fmt.Println("data save successed :", date)
	} else {
		this.task.FileDownStation.InputDateItem(&faces)
	}
	atomic.StoreInt32(&this.state, 0)
}

func CreateDataSaver(task *GatherTask) *DataSaver {
	instance := new(DataSaver)
	instance.SetTask(task)
	instance.chlInput = make(chan *map[int]*common.FaceEx, 100)
	return instance
}

type getPriceWork struct {
	face      *common.FaceEx
	prices    []*common.CodePrice
	chlOutput chan *getPriceWork
}

func (this *getPriceWork) Start() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			this.Finished()
		}
	}()

	if this.face.FileState == 3 {
		ok := false
		filePath := common.GetFilePath(this.face.Date, this.face.Code)
		this.prices, ok = GetTimePriceFromFile(filePath)
		if ok {
			this.face.FileState = 5
		} else {
			this.face.FileState = 0
		}

	}

	this.Finished()
}

func (this *getPriceWork) Finished() {
	this.chlOutput <- this
}

type dataSaveWork struct {
	face   *common.FaceEx
	prices []*common.CodePrice
	task   *GatherTask
}

var tempC int = 0

func (this *dataSaveWork) Start() {
	tempC++
	face := this.face
	if face.FileState == 5 {
		if len(this.prices) == 0 {
			panic("数据获取出错")
		}

		t := time.Now()
		ok := this.task.DB.SaveTimePrices(&face.CodeFace, this.prices)
		fmt.Println(face.CodeFace.Date, "DBend:", time.Since(t))
		//	fmt.Println("耗时：", time.Since(t))
		if ok {
			face.FileState = 1
		} else {
			fmt.Println("数据库保存失败！")
		}
	} else if face.DownTimes > 5 {
		this.task.DB.UpdateFaceState(face.Date, face.Code, 8)
		face.FileState = 8
	}

}

func (this *dataSaveWork) Finished() {
}

func GetTimePriceFromFile(filePath string) ([]*common.CodePrice, bool) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			fmt.Println(filePath)
		}
	}()

	fInfo, err := os.Stat(filePath)
	var prices []*common.CodePrice = nil

	if err != nil || fInfo.IsDir() {
		fmt.Println("路径出错")
		return prices, false
	}
	rows, success := csv.GetRowsFromFile(filePath)

	// errr := os.Remove(filePath)
	// if errr != nil {
	// 	fmt.Println(errr)
	// }

	if len(rows) == 0 {
		fmt.Println("文件获取错误")
	}

	if !success || !checkFileData(&rows[0]) {
		return prices, false
	}
	//fmt.Println(rows)

	//成交时间	成交价格	价格变动	成交量(手)	成交额(元)	性质

	prices = make([]*common.CodePrice, 0, len(rows)-1)
	for i, row := range rows {
		if i == 0 {
			continue
		}
		final, ok := true, true
		price := &common.CodePrice{}
		price.Time, ok = common.GetSecondsFromStr(row[0])
		final = final && ok
		temp, err := strconv.ParseFloat(row[1], 32)
		final = final && err == nil
		price.Price = int(temp * 100)
		volume, err := strconv.Atoi(row[3])
		final = final && err == nil
		if !final {
			return prices, false
		}
		price.Volume = volume
		if strings.Index(row[5], "买") >= 0 {
			price.TradeType = 1
		} else if strings.Index(row[5], "卖") >= 0 {
			price.TradeType = -1
		}
		tempLen := len(prices)
		if tempLen > 0 {
			last := prices[tempLen-1]
			if last.Time == price.Time {
				last.Price = (last.Price*last.Volume + price.Volume*price.Price) / (last.Volume + price.Volume)
				last.TradeType = price.TradeType
				last.Volume = last.Volume + price.Volume
				continue
			}
		}

		prices = append(prices, price)
	}

	return prices, true
}

func checkFileData(row *[]string) bool {
	r := *row
	//成交时间	成交价格	价格变动	成交量(手)	成交额(元)	性质
	if len(r) == 6 &&
		r[0] == "成交时间" &&
		r[1] == "成交价格" &&
		r[2] == "价格变动" &&
		r[3] == "成交量(手)" &&
		r[4] == "成交额(元)" &&
		r[5] == "性质" {
		return true
	}

	return false
}
