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

const (
	con_pagePerCount = 70
)

var (
	columnConfigMap map[string]*Col
)

func init() {
	columnConfigMap = createColMap()
}

func getToken(ctx context.Context, date string) (result *tokenObj, er Model.Err) {
	i := 0
	Tools.LoopCtx(ctx, func() bool {
		if i >= 10 {

			Log.Warn(fmt.Sprintf("getToken loopOut failed times %v", i))
			return false
		}

		i++
		var cookieStr string
		cookieStr, er = waitCookStr(ctx)
		result = &tokenObj{}

		var url = "http://www.iwencai.com/stockpick/load-data?typed=0&preParams=&ts=1&f=1&qs=result_original&selfsectsn=&querytype=stock&searchfilter=&tid=stockpick&w=%E6%B6%A8%E8%B7%8C" +
			date +
			"+%E4%BB%B7%E6%A0%BC" +
			date +
			"+%E6%88%90%E4%BA%A4%E9%87%8F" +
			date +
			"+%E6%8D%A2%E6%89%8B" +
			date +
			"+&queryarea="

		headers := make(map[string]string, 8)

		//	GET /stockpick/load-data?typed=0&preParams=&ts=1&f=1&qs=result_original&selfsectsn=&querytype=stock&searchfilter=&tid=stockpick&w=%E6%B6%A8%E8%B7%8C2026-02-26+%E4%BB%B7%E6%A0%BC2026-02-26+%E6%88%90%E4%BA%A4%E9%87%8F2026-02-26+%E6%8D%A2%E6%89%8B2026-02-26+&queryarea= HTTP/1.1
		//Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7
		headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"
		//	Accept-Encoding: gzip, deflate, br, zstd
		//headers["Accept-Encoding"] = "gzip, deflate, br, zstd"
		//	Accept-Language: zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7
		headers["Accept-Language"] = "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7"
		//	Cache-Control: max-age=0
		headers["Cache-Control"] = "max-age=0"
		//	Connection: keep-alive
		headers["Connection"] = "keep-alive"
		//	Cookie: cid=481a7d61536bcfeabf2484348407282e1736501850; ComputerID=481a7d61536bcfeabf2484348407282e1736501850; WafStatus=0; u_ukey=A10702B8689642C6BE607730E11E6E4A; u_uver=1.0.0; u_dpass=bqtFa7hjOl%2BiRh7cUC5dpQw1UX9TpSjuF%2FjicNkkheHOSQgrgVkmI%2B%2BmvYNIeqp6Hi80LrSsTFH9a%2B6rtRvqGg%3D%3D; u_did=5C977C9FB2F140429E51D2D6C31A810A; u_ttype=WEB; other_uid=Ths_iwencai_Xuangu_phqw1wjmevgbyhvt57e7pz93gucdmdea; ta_random_userid=xeaawsvxop; user=MDpoYW9saWU6Ok5vbmU6NTAwOjI3NDUwMjY2MDo3LDExMTExMTExMTExLDQwOzQ0LDExLDQwOzYsMSw0MDs1LDEsNDA7MSwxMDEsNDA7MiwxLDQwOzMsMSw0MDs1LDEsNDA7OCwwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMSw0MDsxMDIsMSw0MDoyNzo6OjI2NDUwMjY2MDoxNzcyMDA0OTU1Ojo6MTQzMzc3NjAyMDo2MDQ4MDA6MDoxZTQwYzZlYzRlODZkNzU5YzQ1NjM4MDg5Y2U3NjYzZGM6ZGVmYXVsdF81OjA%3D; userid=264502660; u_name=haolie; escapename=haolie; ticket=c84c7cbeff979ff7c1d198ee03e4ae26; user_status=0; utk=bb96658365d4f5e236d6dc7c9a926438; sess_tk=eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzI1NiIsImtpZCI6InNlc3NfdGtfMSIsImJ0eSI6InNlc3NfdGsifQ.eyJqdGkiOiJkYzYzNzZjZTg5ODA2MzQ1OWM3NTZkZThjNDZlMGNlNDEiLCJpYXQiOjE3NzIwMDQ5NTUsImV4cCI6MTc3MjYwOTc1NSwic3ViIjoiMjY0NTAyNjYwIiwiaXNzIjoidXBhc3MuaXdlbmNhaS5jb20iLCJhdWQiOiIyMDIwMTExODUyODg5MDcyIiwiYWN0Ijoib2ZjIiwiY3VocyI6ImE2M2IyOWVmYWY5Y2EyMmYxNDJhM2Q3ZDhmOTA3YzYxODNhNWNjN2Q4ZmE5ZGQ2MDk0MzFkYjk4NmY4YTNhZWYifQ.6BpLWnmm02urYx6BuBIZ0fz6HnWYIPMeK2nG09Lz9wDbOPIsaC0nw8xObv1wBcqVuxQNmsGJfZ87Rjd9BgzK1Q; cuc=pv1exa0r25ol; THSSESSID=ac6584e3b2859910ff6df5b038; PHPSESSID=e1e429b65daae73c9dbc46ce58f987d6; _clck=bllxsm%7C2%7Cg3w%7C0%7C0; v=AyFlmXQG7G_Ji0CbCgofNaQcMOY-zpXAv0I51IP2HSiH6k8Yyx6lkE-SSaoQ; _clsk=zz6imo1qgsq6%7C1772091822130%7C27%7C1%7C; _clsk=zz6imo1qgsq6%7C1772091822130%7C27%7C1%7C; _clsk=zz6imo1qgsq6%7C1772091822130%7C27%7C1%7C
		headers["Cookie"] = cookieStr
		//	Host: www.iwencai.com
		headers["Host"] = "www.iwencai.com"
		//	Sec-Fetch-Dest: document
		headers["Sec-Fetch-Dest"] = "document"
		//	Sec-Fetch-Mode: navigate
		headers["Sec-Fetch-Mode"] = "navigate"
		//	Sec-Fetch-Site: none
		headers["Sec-Fetch-Site"] = "none"
		//	Sec-Fetch-User: ?1
		headers["Sec-Fetch-User"] = "?1"
		//	Upgrade-Insecure-Requests: 1
		headers["Upgrade-Insecure-Requests"] = "1"
		//	User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36
		headers["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36"
		//	sec-ch-ua: "Not(A:Brand";v="8", "Chromium";v="144", "Google Chrome";v="144"
		headers["sec-ch-ua"] = "\"Not(A:Brand\";v=\"8\", \"Chromium\";v=\"144\", \"Google Chrome\";v=\"144\""
		//	sec-ch-ua-mobile: ?0
		headers["sec-ch-ua-mobile"] = "?0"
		//	sec-ch-ua-platform: "Windows"
		headers["sec-ch-ua-platform"] = "Windows"

		data, err := Tools.HttpRequest(url, "GET", headers)
		if err != nil {
			er = Model.Err(fmt.Sprintf("codegather.getToken client.Do err=%v url=%v cookie=%v", err, url, cookieStr))
			Log.Error(string(er))
			resetCook(true)

			Tools.Wait(context.Background(), 30, "codegather.getTokenErr")
			return true
		}

		r1 := jsoniter.Get(data, "data", "result")
		count, _ := strconv.Atoi(r1.Get("code_count").ToString())
		result.count = count
		result.token = r1.Get("token").ToString()
		jsonBlob := []byte(r1.Get("columnsIndexID").ToString())
		err = json.Unmarshal(jsonBlob, &result.columns)

		if err != nil {
			er = Model.Err(fmt.Sprintf("codegather.getToken json.Unmarshal err=%v content=%s cookie=%v", err, string(data), cookieStr))
			Log.Error(string(er))
			Log.Error("----------")
			Log.Error("----------")
			Log.Error("----------")
			Log.Info(url)
			Log.Error("")
			Log.Error("")
			resetCook(false)

			Tools.Wait(context.Background(), 30, "codegather.getToken json.Unmarshal ")
			return true
		}

		return false
	}, func() {
		Log.Info(fmt.Sprintf("getToken ctx.Down"))
	})

	return
}

func getNoFun(token string, gCtx *codeGatherCtx) (list []*Model.CodeFace, errStr Model.Err) {

	newCookStr, errStr := waitCookStr(context.Background())
	if errStr.Exists() {
		return
	}

	client := &http.Client{}
	list = []*Model.CodeFace{}
	url := fmt.Sprintf("%s%s%s%d%s%d%s", "http://www.iwencai.com/stockpick/cache?token=", token, "&p=", gCtx.pageIndex, "&perpage=", con_pagePerCount, "&showType=[%22%22,%22%22,%22onTable%22,%22onTable%22,%22onTable%22,%22onTable%22]")
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
			errStr = Model.Err(fmt.Sprintf("getNoFun try %d times", i))
			return
		}

		codeStr := r1.Get(i, 0).ToString()
		if len(codeStr) == 0 {
			// 如果已有数据行 标识该页爬取已结束
			if !hasRows {
				Tools.Wait(context.Background(), 2, "getNoFunErr")
				es := fmt.Sprintf("getNoFun err: bodyRead fail,url=%s body=%s", url, string(resultStr))
				errStr = Model.Err(es)
				Log.Error(es)
				return
			}
			break
		}

		colConfigs := gCtx.clMap
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

	Log.Info(fmt.Sprintf("page=%d;count=%d \r\n", gCtx.pageIndex, len(list)))
	Tools.Wait(context.Background(), 20, "getNoFunAfter")

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

func createCtx(ctx context.Context, date string) (gCtx *codeGatherCtx, errStr Model.Err) {

	tObj, errStr := getToken(ctx, date)
	if errStr.Exists() {
		return
	}

	perCount := con_pagePerCount
	pageCount := tObj.count / perCount
	if tObj.count%perCount > 0 {
		pageCount++
	}

	gCtx = &codeGatherCtx{
		dateStr:    date,
		pageIndex:  1,
		pageCount:  int32(pageCount),
		totalCount: int32(tObj.count),
		clMap:      createColMap(),
		faceMap:    make(map[int]struct{}, tObj.count),
		faceList:   make([]*Model.CodeFace, 0, tObj.count),
	}

	for k, c := range gCtx.clMap {
		var find bool
		for i, cl := range tObj.columns {
			if cl.Index_name == c.name {
				find = true
				c.index = i
				gCtx.clMap[k].index = i

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

	return
}

func fromWeb(ctx context.Context, gCtx *codeGatherCtx) (errStr Model.Err) {
	var tObj *tokenObj
	tObj, errStr = getToken(ctx, gCtx.dateStr)
	if errStr.Exists() {
		resetCook(false)
		return
	}

	Tools.LoopCtx(ctx, func() bool {
		if gCtx.pageIndex >= gCtx.pageCount {
			return false
		}

		// 因cookie 可能失效  尝试5次
		for t := 0; t < 5; t++ {
			var list []*Model.CodeFace
			list, errStr = getNoFun(tObj.token, gCtx)
			if errStr.Exists() {
				resetCook(false)
				tObj, errStr = getToken(ctx, gCtx.dateStr)
				continue
			}

			for _, face := range list {
				face.Date = gCtx.dateStr
				_, ex := gCtx.faceMap[face.Code]
				if ex {
					fmt.Printf("repead:%d \r\n", face.Code)
				}

				gCtx.faceMap[face.Code] = struct{}{}
				gCtx.faceList = append(gCtx.faceList, face)
			}

			Log.Info(fmt.Sprintf("date:%s page:%d  %d/%d", gCtx.dateStr, gCtx.pageIndex, tObj.count, len(gCtx.faceList)))
			gCtx.pageIndex += 1
			//return false
			break

		}

		return !errStr.Exists()
	}, func() {
		Log.Info(fmt.Sprintf("fromWeb ctx.Down"))
	})

	return
}

func GetNocodesFromWeb(ctx context.Context, date string) (faceList []*Model.CodeFace, errStr Model.Err) {

	gCtx, es := createCtx(ctx, date)
	if es.Exists() {
		errStr = es
		return
	}

	Tools.LoopCtx(ctx, func() bool {
		errStr = fromWeb(ctx, gCtx)
		if !errStr.Exists() {
			faceList = gCtx.faceList
			return false
		}

		Tools.Wait(ctx, 50, "GetNocodesFromWeb")
		return true
	}, func() {
		Log.Info(fmt.Sprintf("GetNocodesFromWeb ctx.Down"))
	})

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
