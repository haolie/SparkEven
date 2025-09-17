package Tools

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// GetShDateTime
//
//	通关web 获取当前上证指数时间
//
// 参数：
// 返回:
//
//	t:
//	err:
func GetShDateTime() (t time.Time, err error) {
	client := &http.Client{}
	url := "https://w.sinajs.cn/list=sh000001"

	//cookieStr := "UOR=cn.bing.com,finance.sina.com.cn,; SINAGLOBAL=117.139.247.210_1737010640.808499; SR_SEL=1_511; U_TRS1=000000d2.a7a4925.6788ca60.a0a405ef; SFA_version8.6.0=2025-02-11%2010%3A18; Apache=117.139.247.210_1739240631.705253; SFA_version8.6.0_click=2; U_TRS2=000000d2.5e28208a4.67aab4ca.2b7491b4; ULV=1739240674514:4:2:2:117.139.247.210_1739240631.705253:1739240654409; FINA_V_S_2=sh601318,sh603300,sh000001,sz002582,sz000725,sz399001; FIN_ALL_VISITED=sh601318%2Csh603300%2Csh000001%2Csz002582%2Csz000725%2Csz399001; rotatecount=1; SFA_version8.7.0=2025-02-18%2015%3A06; SFA_version8.7.0_click=1"
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	//request.Header.Add("Accept-Encoding", "gzip, deflate, br, zstd")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")
	request.Header.Add("Cache-Control", "no-cache")

	//referer
	request.Header.Add("referer", "https://quotes.sina.cn/hs/company/quotes/view/sh000001?from=redirect")
	//request.Header.Add("Host", "data.10jqka.com.cn")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Mobile Safari/537.36")
	//request.Header.Add("Cookie", cookieStr)
	if err != nil {
		return
	}

	resp, err := client.Do(request)
	if err != nil {
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	str := string(data)
	str = strings.Replace(str, "var hq_str_sh000001=", "", -1)
	fmt.Println(str)

	str = strings.Replace(str, "\"", "", -1)
	fmt.Println(str)
	strList := strings.Split(str, ",")
	l := len(strList)
	if l < 5 {
		err = fmt.Errorf("sina request failed")
		return
	}

	str = fmt.Sprintf("%s %s", strList[l-4], strList[l-3])
	fmt.Println(str)
	t, err = time.Parse(fmt.Sprintf("%s %s", time.DateOnly, time.TimeOnly), str)

	return

}
