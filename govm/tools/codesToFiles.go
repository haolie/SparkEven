package tools

import (
	"errors"
	"fmt"
	"os"

	"../nohelper"
)

type CodeGrouper struct {
	basePath  string
	date      string
	Result    bool
	fileIndex int
	count     int
}

func (this *CodeGrouper) Start() {
	faces, ok := nohelper.GetCodeFaces(this.date)
	this.Result = ok
	if !ok {
		return
	}

	index := 0
	var tempArray []int
	fmt.Println(len(faces))
	for code, _ := range faces {
		if index == 0 {
			tempArray = make([]int, this.count)
		}
		if code > 1680000 {
			continue
		}
		tempArray[index] = code
		index++
		if index == this.count {
			index = 0
			//fmt.Println(tempArray)
			this.saveToFiles(tempArray)
		}
	}

	if index > 0 {
		this.saveToFiles(tempArray)
	}

	fmt.Println("文件保存成功！")
}

func (this *CodeGrouper) saveToFiles(faces []int) (err error) {
	filePath := fmt.Sprintf("%s/%d.txt", this.basePath, this.fileIndex)
	this.fileIndex++
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("保存失败")
			fmt.Println("保存失败：", filePath)
		}
	}()
	f, er := os.Create(filePath)
	err = er
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()

	var offset int64 = 0
	for _, code := range faces {
		codestr := fmt.Sprintf("%d\r\n", code)
		codestr = codestr[1:]
		n, er := f.WriteAt([]byte(codestr), offset)
		if er != nil {
			err = er
			return
		}

		offset += int64(n)
	}
	return
}

func CreateGroupper(date string, basePath string, count int) *CodeGrouper {
	cg := new(CodeGrouper)
	cg.date = date
	cg.basePath = basePath
	cg.count = count
	return cg
}
