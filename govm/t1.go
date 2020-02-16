package main

import (
	"fmt"
	"./common"
)

func main() {
	common.HttpDown("https://www.baidu.com/img/pc_1c6e30772d5e4103103bd460913332f9.png", "e:/232.jpg")
	fmt.Println("sucess")
}
