package gather

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	"lyh/SparkEven/govm/common"
)

type FileDowner struct {
	TaskWorkStation
}

type downWork struct {
	face *common.FaceEx
}

func (this *FileDowner) GetState() int32 {
	if this.chlInput == nil {
		return 1
	}

	v := atomic.LoadInt32(&this.state)
	if v > 0 {
		return v
	}

	return int32(len(this.chlInput))
}

func (this *FileDowner) InputDateItem(item *map[int]*common.FaceEx) {
	if this.chlInput != nil {
		this.chlInput <- item
		fmt.Println(len(this.chlInput))
	}
}

func (this *downWork) Start() {
	filePath := common.GetFilePath(this.face.Date, this.face.Code)
	if this.face.State == 1 {
		this.face.FileState = 1
	} else if common.CheckFile(filePath) {
		this.face.FileState = 3
	} else if this.face.FileState == 0 {
		this.face.DownTimes++
		url := "http://stock.gtimg.cn/data/index.php?appn=detail&action=download"
		var codestr string
		if this.face.Code >= 1600000 {
			codestr = fmt.Sprintf("sh%d", this.face.Code-1000000)
		} else {
			codestr = "sz" + strconv.Itoa(this.face.Code)[1:]
		}

		url = fmt.Sprintf("%s&c=%s&d=%s", url, codestr, strings.Replace(this.face.Date, "-", "", 2))
		success := common.HttpDown(url, filePath)
		if success {
			this.face.FileState = 3
			fmt.Println("文件下载成功")
		} else {
			this.face.FileState = 0
			fmt.Println("文件下载失败")
		}
	}
	this.Finished()
}

func (this *downWork) Finished() {
	//this.chlOutput <- &this.face
}

func (this *FileDowner) Start() {
	go func() {

		for face := range this.chlInput {
			this.StartDown(*face)
		}

	}()
}

func (this *FileDowner) StartDown(faces map[int]*common.FaceEx) {
	atomic.StoreInt32(&this.state, 1)
	station := common.CreateDefaultStation(4, len(faces))
	go station.Start()
	for _, face := range faces {
		station.AddItem(&downWork{face})
	}
	station.EndInput()
	station.WaitEnd()
	this.task.DataSaveStation.InputDateItem(&faces)
	atomic.StoreInt32(&this.state, 0)
}

func CreateFileDowner(task *GatherTask) *FileDowner {
	instance := new(FileDowner)
	instance.chlInput = make(chan *map[int]*common.FaceEx, 100)
	instance.SetTask(task)
	return instance
}
