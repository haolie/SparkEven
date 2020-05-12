package nohelper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"../common"
	"../common/csv"
	"../db"

	//"reflect"
	//"strconv"
	jsoniter "github.com/json-iterator/go"
)

var (
	newSession = ""
	lock       sync.Mutex
)

func StartListenSession() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, r *http.Request) {
		lock.Lock()
		defer lock.Unlock()
		if len(newSession) == 0 {
			newSession = strings.TrimLeft(r.URL.Path, "/")
		}

		//fmt.Println(newSession)
		fmt.Fprintf(res, "%v\n", time.Now().Format("2006-01-02 15:04:05"))
	})
	if err := http.ListenAndServe(":9001", mux); err != nil {
		fmt.Println(err)
	}
	fmt.Println("codeId  监听中……")
}

type tokenObj struct {
	count   int
	token   string
	columns []Cl
}

func GetToken(date string) (tokenObj, bool) {
	newSession = ""
	for {
		if len(newSession) > 0 {
			break
		}
		time.Sleep(time.Second * 2)
	}
	var result tokenObj

	client := &http.Client{}
	var url = "http://www.iwencai.com/stockpick/load-data?typed=0&preParams=&ts=1&f=1&qs=result_original&selfsectsn=&querytype=stock&searchfilter=&tid=stockpick&w=%E6%B6%A8%E8%B7%8C" +
		date +
		"+%E4%BB%B7%E6%A0%BC" +
		date +
		"+%E6%88%90%E4%BA%A4%E9%87%8F" +
		date +
		"+%E6%8D%A2%E6%89%8B" +
		date +
		"+&queryarea="
	fmt.Println(url)
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Host", "www.iwencai.com")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.75 Safari/537.36")
	request.Header.Add("Cookie", newSession)
	resp, err := client.Do(request)
	defer resp.Body.Close()
	if err != nil {
		return result, false
	}
	resultStr, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return result, false
	}
	//fmt.Println(string(result))

	r1 := jsoniter.Get(resultStr, "data", "result")
	count, _ := strconv.Atoi(r1.Get("code_count").ToString())
	result.count = count
	result.token = r1.Get("token").ToString()
	jsonBlob := []byte(r1.Get("columnsIndexID").ToString())
	err = json.Unmarshal(jsonBlob, &result.columns)

	if err != nil {
		return result, false
	}
	// fmt.Println(reflect.TypeOf(r1))
	// c1 := r1.ToString()
	// for i := 0; ; i++ {
	// 	str := r1.Get("columnsIndexID", i).ToString()
	// 	if len(str) == 0 {
	// 		break
	// 	}

	// 	fmt.Println(str)

	// }
	//fmt.Println(count)

	//c1 := gojson.Json(string(result)).Get("data").Get("result")
	return result, true
}

type Col struct {
	name  string
	index int
	id    string
}

type Cl struct {
	Uintn         string `json:"uint"`
	Index_name    string `json:"index_name"`
	Key           string `json:"key"`
	Timestamp     string `json:"timestamp"`
	RealIndexName string `json:"realIndexName"`
}

func getNo(token string) {

}

func createColMap() map[string]*Col {
	m := make(map[string]*Col)
	m["ud"] = &Col{"涨跌", -1, "ud"}
	m["lastprice"] = &Col{"收盘价:不复权", -1, "lastprice"}
	m["startprice"] = &Col{"开盘价:不复权", -1, "startprice"}
	m["max"] = &Col{"最高价:不复权", -1, "max"}
	m["min"] = &Col{"最低价:不复权", -1, "min"}
	m["turnoverRate"] = &Col{"换手率", -1, "turnoverRate"}
	m["turnover"] = &Col{"成交额", -1, "turnover"}
	m["volume"] = &Col{"成交量", -1, "volume"}
	return m
}

func getNoFun(token string, index int, perCount int, colConfigs map[string]*Col) ([]*common.FaceEx, bool) {
	client := &http.Client{}
	list := []*common.FaceEx{}
	url := fmt.Sprintf("%s%s%s%d%s%d%s", "http://www.iwencai.com/stockpick/cache?token=", token, "&p=", index, "&perpage=", perCount, "&showType=[%22%22,%22%22,%22onTable%22,%22onTable%22,%22onTable%22,%22onTable%22]")
	fmt.Println(url)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return list, false
	}
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Host", "www.iwencai.com")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.75 Safari/537.36")
	request.Header.Add("Cookie", newSession)
	resp, err := client.Do(request)

	if err != nil {
		return list, false
	}
	defer resp.Body.Close()
	resultStr, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return list, false
	}
	r1 := jsoniter.Get(resultStr, "result")
	for i := 0; ; i++ {
		codestr := r1.Get(i, 0).ToString()
		if len(codestr) == 0 {
			break
		}
		if r1.Get(i, colConfigs["ud"].index).ToString() == "--" {
			continue
		}
		start := time.Now()
		temp := common.FaceEx{}
		temp.Code = common.StrToCodeNum(codestr)
		temp.Change = r1.Get(i, colConfigs["ud"].index).ToFloat32()
		temp.LastPrice = r1.Get(i, colConfigs["lastprice"].index).ToFloat32()
		temp.StartPrice = r1.Get(i, colConfigs["startprice"].index).ToFloat32()
		temp.MaxValue = r1.Get(i, colConfigs["max"].index).ToFloat32()
		temp.MinValue = r1.Get(i, colConfigs["min"].index).ToFloat32()
		temp.TurnoverRate = r1.Get(i, colConfigs["turnoverRate"].index).ToFloat32()
		temp.Turnover = r1.Get(i, colConfigs["turnover"].index).ToFloat32()
		temp.Volume = r1.Get(i, colConfigs["volume"].index).ToInt()
		temp.YestPrice = temp.LastPrice - temp.Change
		temp.Percent = float32(math.Ceil(float64(temp.Change*10000/temp.YestPrice)) / 100)

		list = append(list, &temp)
		fmt.Println("\nspend : ", time.Now().Sub(start))
	}

	return list, true
}

func GetCodeFaces(date string) map[int]*common.FaceEx {
	faces, success := db.Create().GetCodeFaces(date, -1)
	if success && len(faces) > 0 {
		list := map[int]*common.FaceEx{}
		for _, face := range faces {
			temp := common.FaceEx{CodeFace: *face}
			list[face.Code] = &temp
		}
		return list
	}

	fxs := GetNocodesFromWeb(date)
	faces = []*common.CodeFace{}
	for _, f := range fxs {
		faces = append(faces, &f.CodeFace)
	}
	db.Create().SaveCodeFaces(faces)
	return fxs
}

func GetNocodesFromWeb(date string) map[int]*common.FaceEx {
	//tokenObj, _ := GetToken(date)
	tobj, _ := GetToken(date)
	//fmt.Println(tobj)
	perCount := 1000
	pageCount := tobj.count / perCount
	if tobj.count%perCount > 0 {
		pageCount++
	}

	clMaps := createColMap()

	for i, cl := range tobj.columns {
		for k, c := range clMaps {
			if cl.Index_name == c.name {
				c.index = i
				clMaps[k].index = i
			}
		}
	}
	allFace := make(map[int]*common.FaceEx)
	// pageCount = 1
	for i := 0; i < pageCount; i++ {
		list, ok := getNoFun(tobj.token, i+1, perCount, clMaps)
		if ok {
			for _, face := range list {
				face.Date = date
				allFace[face.Code] = face
			}
		}

	}

	// for list := range chl {
	// 	fmt.Print(len(list))

	// }

	fmt.Print(len(allFace))

	// for i, c := range allFace {
	// 	fmt.Println(i)
	// 	c.Println()
	// }
	newSession = ""
	return allFace
}

func GetDatesFromWeb(start string) []string {
	list := []string{}
	client := &http.Client{}
	end := GetLastDateStr()
	var uri = "http://quotes.money.163.com/service/chddata.html?code=0000001&start=" +
		strings.Replace(start, "-", "", 2) +
		"&end=" +
		strings.Replace(end, "-", "", 2) +
		"&fields=TCLOSE;HIGH;LOW;TOPEN;LCLOSE;CHG;PCHG;VOTURNOVER;VATURNOVER"
	fmt.Println(uri)
	request, err := http.NewRequest("GET", uri, nil)
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.75 Safari/537.36")
	resp, err := client.Do(request)
	defer resp.Body.Close()
	if err != nil {
		return list
	}

	body, err := ioutil.ReadAll(resp.Body)
	rows, _ := csv.GetRowsFromBytes(body)

	for i := len(rows) - 1; i >= 0; i-- {
		row := rows[i]
		if i == 0 {
			continue
		}

		list = append(list, strings.Split(row[0], ",")[0])
	}
	// for i, row := range rows {
	// 	if i == 0 {
	// 		continue
	// 	}

	// 	list = append(list, strings.Split(row[0], ",")[0])
	// }
	return list
}

func GetLastDateStr() string {
	client := &http.Client{}
	var url = "http://hq.sinajs.cn/list=sh000001"
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.75 Safari/537.36")
	resp, err := client.Do(request)
	defer resp.Body.Close()
	if err != nil {
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	str := string(body)
	return str[len(str)-26 : len(str)-16]
}
