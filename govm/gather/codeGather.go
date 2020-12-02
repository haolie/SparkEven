package gather

import (
	"fmt"

	"../nohelper"
)

type CodeGather struct {
	TaskWorkStation
}

func (this *CodeGather) Start() {
	this.state = 1
	go this.StartWork()

}

func (this *CodeGather) StartWork() {
	dates := nohelper.GetDatesFromWeb(this.task.StartDate)
	faces, ok := this.task.DB.GetCodeFaces("", 1912261)
	var list []string
	if ok {
	LOOP:
		for _, date := range dates {
			for _, face := range faces {
				if face.Date == date && face.State == 1 {
					continue LOOP
				}

			}
			list = append(list, date)
		}
	}

	if len(list) == 0 {
		this.state = 0
		this.task.EndTask()
	}

	// this.chlOutput = make(chan *map[int]*common.FaceEx, len(dates))
	for _, date := range list {
		fmt.Println("…………………………")
		fmt.Println("…………………………")
		fmt.Println("…………………………")
		fmt.Printf("start down：%s\n", date)
		fmt.Println("…………………………")
		faces, ok := nohelper.GetCodeFaces(date)
		if ok {
			this.task.FileDownStation.InputDateItem(&faces)
			fmt.Println("…………………………")
			fmt.Println("…………………………")
			fmt.Println("…………………………")
			fmt.Printf("finish down：%s\n", date)
			fmt.Println("…………………………")
		} else {
			fmt.Println("…………………………")
			fmt.Println("…………………………")
			fmt.Println("…………………………")
			fmt.Printf("down faild：%s\n", date)
			fmt.Println("…………………………")
		}

	}
	this.state = 0
}

func CreateCodeGather(task *GatherTask) *CodeGather {
	instance := new(CodeGather)
	instance.SetTask(task)
	return instance
}
