package main

import (
	"fmt"
	"time"

	"SparkEven/govm/gather"
	"SparkEven/govm/nohelper"
)

func main() {
	fmt.Println("开始了 ")
	go nohelper.StartListenSession()
	go work()
	select {}
}

func work() {
	for {
		t := time.Now()
		fmt.Println("任务开始：", t)
		task := gather.CreateGatherTask("2022-01-10")
		task.StartTask()
		task.WaitEnd()
		fmt.Println("任务结束：", time.Now())
		fmt.Println("耗时：", time.Since(t))
		time.Sleep(600 * time.Second)
	}
}
