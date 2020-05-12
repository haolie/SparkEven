package common

import (
	"fmt"
	"sync"
	"time"
)

type WorkStationer interface {
	Init()
	AddItem(item WorkItem)
	Start()
	WaitEnd()
	Finished()
	Destory()
}

type WorkStation struct {
	progressCount  int
	totalItemCount int
	finishedCount  int
	maxItemCount   int
	chlWorkItem    chan WorkItem
	exitChl        chan bool
	wg             sync.WaitGroup
	workers        []Worker
	status         int // -1  已销毁  0 未初始化 1 已初始化  2 工作中 3 等待销毁
	startTime      *time.Time
}

//初始化
func (this *WorkStation) Init() {
	if this.status != 0 {
		return
	}

	this.workers = make([]Worker, this.progressCount)
	for i := 0; i < this.progressCount; i++ {
		this.workers[i] = NewWork()
	}

	this.exitChl = make(chan bool, this.progressCount)
	this.chlWorkItem = make(chan WorkItem, this.maxItemCount)
	this.status = 1
}

//添加工作
func (this *WorkStation) AddItem(item WorkItem) {
	this.chlWorkItem <- item
	this.totalItemCount++
	if this.totalItemCount == this.maxItemCount {
		close(this.chlWorkItem)
	}
	if this.startTime == nil {
		t := time.Now()
		this.startTime = &t
	}
}

func (this *WorkStation) WaitEnd() {
	this.wg.Wait()
}

//开始
func (this *WorkStation) Start() {
	if this.status != 1 {
		return
	}

	for _, work := range this.workers {
		go work.Working(this.chlWorkItem, this.exitChl)
		this.wg.Add(1)
	}

	for i := 0; i < this.progressCount; i++ {
		<-this.exitChl
		this.wg.Done()
	}
	close(this.exitChl)
	fmt.Println("workend:", time.Since(*this.startTime))
	this.Finished()
}

//完成
func (this *WorkStation) Finished() {
	this.status = 3
}

//销毁
func (this *WorkStation) Destory() {

}

func CreateDefaultStation(workerCount int, itemCount int) WorkStationer {
	station := &WorkStation{}
	station.progressCount = workerCount
	station.maxItemCount = itemCount
	station.Init()
	return station
}
