package gather

import (
	"sync"

	"SparkEven/govm/common"
	"SparkEven/govm/db"
)

type GatherTask struct {
	Status            int // 0 未开始  // 1 正在运行中  // 3 已结束
	StartDate         string
	CodeGatherStation ITask
	FileDownStation   ITask
	DataSaveStation   ITask
	DB                *db.MysqlSuport
	wg                sync.WaitGroup
}

func (this *GatherTask) StartTask() {
	if this.Status != 0 {
		return
	}
	this.Status = 1
	this.wg.Add(1)
	this.CodeGatherStation.Start()
	this.DataSaveStation.Start()
	this.FileDownStation.Start()
}

func (this *GatherTask) EndTask() {
	if this.Status != 1 {
		return
	}

	if st := this.CodeGatherStation.GetState(); st > 0 {
		return
	}

	if st := this.FileDownStation.GetState(); st > 0 {
		return
	}

	if st := this.DataSaveStation.GetState(); st > 0 {
		return
	}

	// if this.CodeGatherStation.GetState() > 0 ||
	// 	this.FileDownStation.GetState() > 0 ||
	// 	this.DataSaveStation.GetState() > 0 {
	// 	return
	// }

	this.CodeGatherStation.Close()
	this.FileDownStation.Close()
	this.DataSaveStation.Close()
	this.Status = 3
	db.CloseDefault()
	this.DB = nil
	this.wg.Done()
}

func (this *GatherTask) WaitEnd() {
	this.wg.Wait()
}

func CreateGatherTask(startDate string) *GatherTask {
	task := new(GatherTask)
	task.DB = db.Create()
	task.StartDate = startDate
	task.CodeGatherStation = CreateCodeGather(task)
	task.FileDownStation = CreateFileDowner(task)
	task.DataSaveStation = CreateDataSaver(task)
	return task
}

type ITask interface {
	SetId(id int)
	GetId() int
	GetState() int32
	SetState(state int32)
	SetTask(task *GatherTask)
	InputDateItem(item *map[int]*common.FaceEx)
	SetInputChl(input chan *map[int]*common.FaceEx)
	GetOutputChl() chan *map[int]*common.FaceEx
	Start()
	Close()
}

type CodeFacePrices struct {
	common.FaceEx
	prices []*common.CodePrice
}

type TaskWorkStation struct {
	id        int
	state     int32
	task      *GatherTask
	chlInput  chan *map[int]*common.FaceEx
	chlOutput chan *map[int]*common.FaceEx
}

// SetId(id int)
// GetId() int
func (this *TaskWorkStation) SetId(id int) {
	this.id = id
}

func (this *TaskWorkStation) GetId() int {
	return this.id
}

func (this *TaskWorkStation) Close() {
	if this.chlInput != nil {
		close(this.chlInput)
	}
}

func (this *TaskWorkStation) InputDateItem(item *map[int]*common.FaceEx) {
	if this.chlInput != nil {
		this.chlInput <- item
	}
}

func (this *TaskWorkStation) GetOutputChl() chan *map[int]*common.FaceEx {
	return this.chlOutput
}

func (this *TaskWorkStation) SetInputChl(input chan *map[int]*common.FaceEx) {
	this.chlInput = input
}

func (this *TaskWorkStation) SetTask(task *GatherTask) {
	this.task = task
}

func (this *TaskWorkStation) GetState() int32 {
	return this.state
}

func (this *TaskWorkStation) SetState(state int32) {
	this.state = state
}
