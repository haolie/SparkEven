package common

type WorkItem interface {
	Start()
	//Onerror() bool
	//Onprogress()
	Finished()
	//Destory()
}
