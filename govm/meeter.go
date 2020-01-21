package main

import(
	"./nohelper"
	"fmt"
)

func main(){
	fmt.Println("开始了 ")
	go nohelper.StartListenSession()
    nohelper.GetNocodesFromWeb("2020-01-20")
	//fmt.Println("token:",tk)
	select{}
}