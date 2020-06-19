package common

type Worker interface {
	Working(chlItems chan WorkItem, exitChl chan bool)
	Onerror() bool
	Onprogress()
	Finished()
	Destory()
}

type Work struct {
	chlItems chan WorkItem
	exitChl  chan bool
	status   int //0 空闲  1 工作中  3 已销毁
	station  WorkStationer
}

func (this *Work) Working(chlItems chan WorkItem, exitChl chan bool) {
	if this.status != 0 {
		return
	}
	this.status = 1
	this.chlItems = chlItems
	this.exitChl = exitChl

	for {
		item, ok := <-chlItems
		if !ok {
			break
		}
		item.Start()
		this.station.AddFinished(1)
	}

	exitChl <- true
	this.Finished()
}

func (this *Work) Onerror() bool {
	return true
}

func (this *Work) Onprogress() {

}

func (this *Work) Finished() {
	this.status = 3
}

func (this *Work) Destory() {

}

func NewWork(station WorkStationer) *Work {
	temp := &Work{}
	temp.station = station
	temp.status = 0
	return temp
}
