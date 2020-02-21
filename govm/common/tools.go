package common

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func StrToCodeNum(str string) int {
	str = strings.Replace(str, ".SH", "", -1)
	str = strings.Replace(str, ".SZ", "", -1)
	temp, _ := strconv.Atoi(str)
	if temp < 1000000 {
		temp += 1000000
	}
	return temp
}

func CodeNumToStr(num int) string {
	if num > 1600000 {
		return fmt.Sprintf("%d.sh", num-1000000)
	} else {
		temp := strconv.Itoa(num)[1:]
		return temp + ".sz"
	}

}

func HttpDown(url string, savePath string) bool {
	res, err := http.Get(url)
	if err != nil {
		return false
	}

	f, err := os.Create(savePath)
	if err != nil {
		panic(err)
	}
	io.Copy(f, res.Body)

	return true
}

func (face CodeFace) Println() {
	fmt.Println("代码：", face.Code)
	fmt.Println("日期：", face.Date)
	fmt.Println("涨跌：", face.Change)
	fmt.Println("收盘价：", face.LastPrice)
	fmt.Println("开盘价：", face.StartPrice)
	fmt.Println("最高价：", face.MaxValue)
	fmt.Println("最低价：", face.MinValue)
	fmt.Println("换手率：", face.TurnoverRate)
	fmt.Println("换手：", face.Turnover)
	fmt.Println("成交量：", face.Volume)
	fmt.Println("昨日收：", face.YestPrice)
	fmt.Println("涨跌率：", face.Percent)
}

func GetFilePath(date string, no int) string {
	fName := fmt.Sprintf("./datefiles/%s/%s_%d.xls", date, date, no)
	path, _ := filepath.Abs(fName)
	return path
}


func GetSecondsFromStr(timestr string) (int, bool) {
	strs := strings.Split(timestr, ":")
	if len(strs) == 3 {
		seconds, err1 := strconv.Atoi(strs[2])
		m, err2 := strconv.Atoi(strs[1])
		h, err3 := strconv.Atoi(strs[0])

		if err1 != nil || err2 != nil || err3 != nil {
			return -1, false
		}

		return seconds + m*60 + h*60*60, true

	}

	return -1, false

}
