package csv

import (
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func GetRowsFromFile(file string) ([][]string, bool) {
	f, err := os.Open(file)
	if err != nil {
		return [][]string{}, false
	}

	utfReader := transform.NewReader(f, simplifiedchinese.GB18030.NewDecoder())
	all, err := ioutil.ReadAll(utfReader)
	if err != nil {
		return [][]string{}, false
	}

	return GetRowsFromBytes(all)
}

func GetRowsFromBytes(data []byte) ([][]string, bool) {
	//fmt.Println(string(data))
	items := strings.Split(string(data), "\n")
	list := [][]string{}
	for _, item := range items {
		if len(item) > 0 {
			list = append(list, strings.Split(item, "\t"))
		}
	}

	// for i, item := range list {
	// 	fmt.Println(i)
	// 	fmt.Printf("\n")
	// 	fmt.Println(strings.Split(item[0], ",")[0])
	// }

	return list, true
}
