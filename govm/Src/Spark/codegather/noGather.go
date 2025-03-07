package codegather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"SparkEven/govm/Src/Common/Log"
	"SparkEven/govm/Src/Common/Tools"
	"SparkEven/govm/Src/Common/csv"
	"SparkEven/govm/Src/Model"
	"github.com/PuerkitoBio/goquery"
	jsoniter "github.com/json-iterator/go"
)

var (
	columnConfigMap map[string]*Col
)

func init() {
	columnConfigMap = createColMap()
}

func getToken(ctx context.Context, date string) (result *tokenObj, er Model.Err) {
	cookieStr, er := waitCookStr(ctx)
	result = &tokenObj{}

	//cookieStr = "other_uid=Ths_iwencai_Xuangu_e5fmb9qchj3cm6xv3l4nelwsyrp68w1b; ta_random_userid=1mloc8i7wd; cid=481a7d61536bcfeabf2484348407282e1736501850; u_ukey=A10702B8689642C6BE607730E11E6E4A; u_uver=1.0.0; u_dpass=dXNNRcqmggrrHl%2BbDMF6kabLkCS3lvXho4XHaD%2B0IYy5JFpH92JkrfXMX5U3obovHi80LrSsTFH9a%2B6rtRvqGg%3D%3D; u_did=3ADAE2A5B02C4FAC9BF2F0D5C7756EC8; u_ttype=WEB; ttype=WEB; user=MDpoYW9saWU6Ok5vbmU6NTAwOjI3NDUwMjY2MDo3LDExMTExMTExMTExLDQwOzQ0LDExLDQwOzYsMSw0MDs1LDEsNDA7MSwxMDEsNDA7MiwxLDQwOzMsMSw0MDs1LDEsNDA7OCwwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMSw0MDsxMDIsMSw0MDoyNzo6OjI2NDUwMjY2MDoxNzM5NDM2NjAzOjo6MTQzMzc3NjAyMDo2MDQ4MDA6MDoxZTUzY2U1MDdhZGQyMGE3MGZjYmIyMjkwZjQ1ODUyZGY6ZGVmYXVsdF80OjA%3D; userid=264502660; u_name=haolie; escapename=haolie; ticket=9c547d07a178c0acfffa1aa1a200e9d6; user_status=0; utk=71d6ea1e97d4d942da4f56e1f01508ac; PHPSESSID=93ab62e95d8fef90927b10df56e98b25; cid=481a7d61536bcfeabf2484348407282e1736501850; ComputerID=481a7d61536bcfeabf2484348407282e1736501850; WafStatus=0; v=A-uvOxpICJ-pvVRLjndFb9zJegTQAP-CeRTDNl1oxyqB_AX65dCP0onkU55u"
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

	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	//request.Header.Add("Accept-Encoding", "gzip, deflate, br, zstd")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Host", "www.iwencai.com")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36")
	request.Header.Add("Cookie", cookieStr)
	resp, err := client.Do(request)
	if err != nil {
		er = Model.Err(fmt.Sprintf("codegather.getToken client.Do err=%v", err))
		return
	}

	defer resp.Body.Close()
	resultByte, err := io.ReadAll(resp.Body)
	if err != nil {
		er = Model.Err(fmt.Sprintf("codegather.getToken io.ReadAll err=%v", err))
		return
	}

	r1 := jsoniter.Get(resultByte, "data", "result")
	count, _ := strconv.Atoi(r1.Get("code_count").ToString())
	result.count = count
	result.token = r1.Get("token").ToString()
	jsonBlob := []byte(r1.Get("columnsIndexID").ToString())
	err = json.Unmarshal(jsonBlob, &result.columns)

	if err != nil {
		er = Model.Err(fmt.Sprintf("codegather.getToken json.Unmarshal err=%v", err))
	}

	return
}

func getNoFun(token string, index int, perCount int, colConfigs map[string]*Col) (list []*Model.CodeFace, errStr Model.Err) {

	newCookStr, errStr := waitCookStr(context.Background())
	if errStr.Exists() {
		return
	}

	client := &http.Client{}
	list = []*Model.CodeFace{}
	url := fmt.Sprintf("%s%s%s%d%s%d%s", "http://www.iwencai.com/stockpick/cache?token=", token, "&p=", index, "&perpage=", perCount, "&showType=[%22%22,%22%22,%22onTable%22,%22onTable%22,%22onTable%22,%22onTable%22]")
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errStr = Model.Err(fmt.Sprintf("getNoFun err: NewRequest fail, %s", err.Error()))
		return
	}

	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Host", "www.iwencai.com")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.75 Safari/537.36")
	request.Header.Add("Cookie", newCookStr)
	resp, err := client.Do(request)

	if err != nil {
		errStr = Model.Err(fmt.Sprintf("getNoFun err: client.Do fail, %s,url=%s", err.Error(), url))
		return
	}

	defer resp.Body.Close()
	resultStr, err := io.ReadAll(resp.Body)

	if err != nil {
		errStr = Model.Err(fmt.Sprintf("getNoFun err: client.Do fail, %s,url=%s", err.Error(), url))
		return
	}

	r1 := jsoniter.Get(resultStr, "result")
	hasRows := false
	for i := 0; ; i++ {
		if i >= 100 {
			errStr = "getNoFun try 100 times"
			return
		}

		codeStr := r1.Get(i, 0).ToString()
		if len(codeStr) == 0 {
			// 如果已有数据行 标识该页爬取已结束
			if !hasRows {
				time.Sleep(time.Second * 2)
				Log.Error(url)
				errStr = Model.Err(fmt.Sprintf("getNoFun err: bodyRead fail,body=%s", string(resultStr)))

				return
			}
			break
		}

		// 标记有数据行
		hasRows = true
		if r1.Get(i, colConfigs[Col_Ud].index).ToString() == "--" ||
			r1.Get(i, colConfigs[Col_LastPrice].index).ToString() == "--" ||
			r1.Get(i, colConfigs[Col_StartPrice].index).ToString() == "--" ||
			r1.Get(i, colConfigs[Col_Min].index).ToString() == "--" ||
			r1.Get(i, colConfigs[Col_Max].index).ToString() == "--" {
			continue
		}

		temp := Model.CodeFace{}
		temp.Code = Tools.StrToCodeNum(codeStr)
		temp.Change = r1.Get(i, colConfigs[Col_Ud].index).ToFloat32()
		temp.LastPrice = r1.Get(i, colConfigs[Col_LastPrice].index).ToFloat32()
		temp.StartPrice = r1.Get(i, colConfigs[Col_StartPrice].index).ToFloat32()
		temp.MaxValue = r1.Get(i, colConfigs[Col_Max].index).ToFloat32()
		temp.MinValue = r1.Get(i, colConfigs[Col_Min].index).ToFloat32()
		temp.TurnoverRate = r1.Get(i, colConfigs[Col_TurnoverRate].index).ToFloat32()
		temp.Turnover = r1.Get(i, colConfigs[Col_Turnover].index).ToFloat32()
		temp.Volume = r1.Get(i, colConfigs[Col_Volume].index).ToInt()
		temp.YestPrice = temp.LastPrice - temp.Change
		temp.Percent = float32(math.Ceil(float64(temp.Change*10000/temp.YestPrice)) / 100)

		list = append(list, &temp)
	}

	Log.Info(fmt.Sprintf("page=%d;count=%d \r\n", index, len(list)))
	time.Sleep(time.Second * 1)

	return
}

func createColMap() map[string]*Col {
	m := make(map[string]*Col)
	m[Col_Ud] = &Col{"涨跌", -1, "ud", true}
	m[Col_LastPrice] = &Col{"收盘价:不复权", -1, "lastprice", true}
	m[Col_StartPrice] = &Col{"开盘价:不复权", -1, "startprice", true}
	m[Col_Max] = &Col{"最高价:不复权", -1, "max", true}
	m[Col_Min] = &Col{"最低价:不复权", -1, "min", true}
	m[Col_TurnoverRate] = &Col{"换手率", -1, "turnoverRate", true}
	m[Col_Turnover] = &Col{"成交额", -1, "turnover", true}
	m[Col_Volume] = &Col{"成交量", -1, "volume", true}
	return m
}

func GetNocodesFromWeb(ctx context.Context, date string) (faceList []*Model.CodeFace, errStr Model.Err) {

	tObj, errStr := getToken(ctx, date)
	if errStr.Exists() {
		return
	}

	perCount := 70
	pageCount := tObj.count / perCount
	if tObj.count%perCount > 0 {
		pageCount++
	}

	clMaps := createColMap()
	for k, c := range clMaps {
		var find bool
		for i, cl := range tObj.columns {
			if cl.Index_name == c.name {
				find = true
				c.index = i
				clMaps[k].index = i

				if c.dateCheck && cl.Timestamp != strings.Replace(date, "-", "", -1) {
					errStr = Model.Err(fmt.Sprintf("数据时间不匹配  列名:%s  数据时间:%s  目标时间:%s  ", c.name, cl.Timestamp, date))
				}
			}
		}

		if !find {
			errStr = Model.Err(fmt.Sprintf("GetNocodesFromWeb 未找到数据列【%s】  ", c.name))
			return
		}
	}

	faceList = make([]*Model.CodeFace, 0, tObj.count)
	existsMap := make(map[int]struct{})
	for i := 0; i < pageCount; i++ {
		// 每请求15次后刷新cookie
		if (i+1)%15 == 0 {
			// 清空cookie 后重启获取token
			resetCook()
			tObj, errStr = getToken(ctx, date)
			if errStr.Exists() {
				return
			}
		}

		// 因cookie 可能失效  尝试5次
		for t := 0; t < 5; t++ {
			var list []*Model.CodeFace
			list, errStr = getNoFun(tObj.token, i+1, perCount, clMaps)
			if errStr.Exists() {
				resetCook()
				continue
			}

			for _, face := range list {
				face.Date = date
				_, ex := existsMap[face.Code]
				if ex {
					fmt.Printf("repead:%d \r\n", face.Code)
				}

				existsMap[face.Code] = struct{}{}
				faceList = append(faceList, face)
			}

			Log.Info(fmt.Sprintf("date:%s page:%d  %d/%d", date, i+1, tObj.count, len(faceList)))
			break
		}
	}

	Log.Info(fmt.Sprintf("gatherFinished date=%s  num=%d  total=%d", date, len(faceList), tObj.count))
	resetCook()
	return
}

// 弃用
func GetDatesFromWeb(start string) []string {
	list := []string{}
	client := &http.Client{}
	end := "2025-02-10"
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
	if err != nil {
		return list
	}

	defer resp.Body.Close()

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

/*
查询最后交易日  （不准确，当日收盘以后获取结果依然是前一交易日）
*/
func GetLastDateStr() (date string, errStr Model.Err) {
	client := &http.Client{}
	var url = "https://finance.sina.com.cn/realstock/lastfive/sh000001.js"
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.75 Safari/537.36")
	resp, err := client.Do(request)

	if err != nil {
		errStr = Model.Err(fmt.Sprintf("GetLastDateStr request url=%s err=%v \r\n", url, err))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	str := string(body)

	tempList := strings.Split(str, "\n")
	if len(tempList) != 2 {
		errStr = Model.Err(fmt.Sprintf("GetLastDateStr 数据解析失败 url=%s data=%s ", url, str))
		return
	}

	tempList = strings.Split(tempList[0], "=")
	if len(tempList) != 2 {
		errStr = Model.Err(fmt.Sprintf("GetLastDateStr 数据解析失败 url=%s data=%s ", url, str))
		return
	}

	lastfive := jsoniter.Get([]byte(tempList[1]), "lastfive")
	d := lastfive.Get(0)
	date = d.Get("d").ToString()

	return
}

func getCodeDateList() (dateList []string, errStr Model.Err) {
	client := &http.Client{}
	url := "https://data.10jqka.com.cn/gzqh/index/order/desc/ajax/1/"

	cookieStr := "Hm_lvt_722143063e4892925903024537075d0d=1736935523,1739257335; HMACCOUNT=847971CC52BC2CB1; Hm_lvt_929f8b362150b1f77b477230541dbbc2=1736935523,1739257335; Hm_lvt_78c58f01938e4d85eaf619eae71b4ed1=1736935523,1739257335; historystock=600839; spversion=20130314; user=MDpoYW9saWU6Ok5vbmU6NTAwOjI3NDUwMjY2MDo3LDExMTExMTExMTExLDQwOzQ0LDExLDQwOzYsMSw0MDs1LDEsNDA7MSwxMDEsNDA7MiwxLDQwOzMsMSw0MDs1LDEsNDA7OCwwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMSw0MDsxMDIsMSw0MDoyNzo6OjI2NDUwMjY2MDoxNzM5NDM2NjAzOjo6MTQzMzc3NjAyMDo2MDQ4MDA6MDoxZTUzY2U1MDdhZGQyMGE3MGZjYmIyMjkwZjQ1ODUyZGY6ZGVmYXVsdF80OjA%3D; userid=264502660; u_name=haolie; escapename=haolie; ticket=9c547d07a178c0acfffa1aa1a200e9d6; utk=71d6ea1e97d4d942da4f56e1f01508ac; log=; Hm_lpvt_722143063e4892925903024537075d0d=1739518962; Hm_lpvt_929f8b362150b1f77b477230541dbbc2=1739518962; __utma=156575163.1376207255.1739518979.1739518979.1739518979.1; __utmc=156575163; __utmz=156575163.1739518979.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); Hm_lvt_f79b64788a4e377c608617fba4c736e2=1739521246; Hm_lvt_60bad21af9c824a4a0530d5dbf4357ca=1739521246; Hm_lpvt_60bad21af9c824a4a0530d5dbf4357ca=1739521355; Hm_lpvt_78c58f01938e4d85eaf619eae71b4ed1=1739521355; Hm_lpvt_f79b64788a4e377c608617fba4c736e2=1739521355; v=A9qyl9D8aXv6x-XX6YhE6MUiK4v5C17r0I_SieRThm04V3Q1zJuu9aAfIpq3"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errStr = Model.Err(fmt.Sprintf("getCodeDateList http.NewRequest err: %v", err))
		return
	}

	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	//request.Header.Add("Accept-Encoding", "gzip, deflate, br, zstd")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Host", "data.10jqka.com.cn")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36")
	request.Header.Add("Cookie", cookieStr)

	resp, err := client.Do(request)
	if err != nil {
		errStr = Model.Err(fmt.Sprintf("getCodeDateList client.Do err: %v", err))
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	list := doc.Find("#j_date").Find(".msearch-list").Find("li")
	dateList = make([]string, 0, len(list.Nodes))
	for _, n := range list.Nodes {
		temp := n.Attr[0].Val
		dateList = append(dateList, fmt.Sprintf("%s-%s-%s", temp[0:4], temp[4:6], temp[6:8]))
	}

	return
}

func getCurCodeDate() (dateStr string) {
	t := time.Now()
	if t.Hour() > 6 {
		t.AddDate(0, 0, -1)
	}

	return t.Format("2006-01-02")
}
