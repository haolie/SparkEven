package codegather

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestCodeDate(t *testing.T) {
	//list, err := getCodeDateList()
	//if err.Exists() {
	//	t.Error(err)
	//} else {
	//	fmt.Println(list)
	//}

	TryMoney()

	//str := "'15:00:00', '667400', '8.890', 'DOWN'"
	//str = strings.Replace(str, "'", "", -1)
	//subValueStr := strings.Split(str, ",")
	//tempPrice, err := strconv.ParseFloat(strings.Trim(subValueStr[2], " "), 32)
	//if err != nil {
	//	log.Fatal(err)
	//} else {
	//	fmt.Println(tempPrice)
	//}
}

func TryPage() (err error) {
	client := &http.Client{}
	url := "https://data.10jqka.com.cn/gzqh/index/order/desc/ajax/1/"

	cookieStr := "Hm_lvt_722143063e4892925903024537075d0d=1736935523,1739257335; HMACCOUNT=847971CC52BC2CB1; Hm_lvt_929f8b362150b1f77b477230541dbbc2=1736935523,1739257335; Hm_lvt_78c58f01938e4d85eaf619eae71b4ed1=1736935523,1739257335; historystock=600839; spversion=20130314; user=MDpoYW9saWU6Ok5vbmU6NTAwOjI3NDUwMjY2MDo3LDExMTExMTExMTExLDQwOzQ0LDExLDQwOzYsMSw0MDs1LDEsNDA7MSwxMDEsNDA7MiwxLDQwOzMsMSw0MDs1LDEsNDA7OCwwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMSw0MDsxMDIsMSw0MDoyNzo6OjI2NDUwMjY2MDoxNzM5NDM2NjAzOjo6MTQzMzc3NjAyMDo2MDQ4MDA6MDoxZTUzY2U1MDdhZGQyMGE3MGZjYmIyMjkwZjQ1ODUyZGY6ZGVmYXVsdF80OjA%3D; userid=264502660; u_name=haolie; escapename=haolie; ticket=9c547d07a178c0acfffa1aa1a200e9d6; utk=71d6ea1e97d4d942da4f56e1f01508ac; log=; Hm_lpvt_722143063e4892925903024537075d0d=1739518962; Hm_lpvt_929f8b362150b1f77b477230541dbbc2=1739518962; __utma=156575163.1376207255.1739518979.1739518979.1739518979.1; __utmc=156575163; __utmz=156575163.1739518979.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); Hm_lvt_f79b64788a4e377c608617fba4c736e2=1739521246; Hm_lvt_60bad21af9c824a4a0530d5dbf4357ca=1739521246; Hm_lpvt_60bad21af9c824a4a0530d5dbf4357ca=1739521355; Hm_lpvt_78c58f01938e4d85eaf619eae71b4ed1=1739521355; Hm_lpvt_f79b64788a4e377c608617fba4c736e2=1739521355; v=A9qyl9D8aXv6x-XX6YhE6MUiK4v5C17r0I_SieRThm04V3Q1zJuu9aAfIpq3"
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	//request.Header.Add("Accept-Encoding", "gzip, deflate, br, zstd")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Host", "data.10jqka.com.cn")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36")
	request.Header.Add("Cookie", cookieStr)
	if err != nil {
		return err
	}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	list := doc.Find("#j_date").Find(".msearch-list").Find("li")

	for _, n := range list.Nodes {
		temp := n.Attr[0].Val

		fmt.Printf("%s-%s-%s", temp[0:4], temp[4:6], temp[6:8])
	}

	//	fmt.Println(list.Text())

	return nil
}

func TryItemPrice() (err error) {
	client := &http.Client{}
	url := "https://vip.stock.finance.sina.com.cn/quotes_service/view/CN_TransListV2.php?num=150000&symbol=sz002582&rn=28950312"

	//cookieStr := "UOR=cn.bing.com,finance.sina.com.cn,; SINAGLOBAL=117.139.247.210_1737010640.808499; SR_SEL=1_511; U_TRS1=000000d2.a7a4925.6788ca60.a0a405ef; SFA_version8.6.0=2025-02-11%2010%3A18; Apache=117.139.247.210_1739240631.705253; SFA_version8.6.0_click=2; U_TRS2=000000d2.5e28208a4.67aab4ca.2b7491b4; ULV=1739240674514:4:2:2:117.139.247.210_1739240631.705253:1739240654409; FINA_V_S_2=sh601318,sh603300,sh000001,sz002582,sz000725,sz399001; FIN_ALL_VISITED=sh601318%2Csh603300%2Csh000001%2Csz002582%2Csz000725%2Csz399001; rotatecount=1; SFA_version8.7.0=2025-02-18%2015%3A06; SFA_version8.7.0_click=1"
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	//request.Header.Add("Accept-Encoding", "gzip, deflate, br, zstd")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Connection", "keep-alive")
	//request.Header.Add("Host", "data.10jqka.com.cn")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36")
	//request.Header.Add("Cookie", cookieStr)
	if err != nil {
		return err
	}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	str := string(data)

	list := strings.Split(str, "\n")
	//fmt.Println(list[0])
	// trade_item_list[0] = new Array('15:00:00', '667400', '8.890', 'DOWN');
	for _, v := range list {
		tempIndex := strings.Index(v, "trade_item_list")
		if tempIndex > -1 && len(v) > 40 {
			tempIndex = strings.IndexRune(v, '(')

			valueStr := v[tempIndex+1:]
			valueStr = valueStr[:len(valueStr)-2]
			fmt.Println(valueStr)
		}
	}

	return err
}

func TryMoney() (err error) {
	client := &http.Client{}
	outUrl := "http://180.184.201.119:18519/WeiXinMoneyApi"

	//var playerId = base.FormData["playerId"].ToString();
	//var outBillNo = base.FormData["outBillNo"].ToString();
	//var orderId = base.FormData["cpOrderId"].ToString();
	//string extra = base.FormData["extralInfo"].ToString();
	//string money = base.FormData["money"].ToString();
	//string t = base.FormData["t"].ToString();
	//string sign = base.FormData["sign"].ToString();
	//
	//fv := url.Values{}
	//fv.Set("", "")

	//cookieStr := "UOR=cn.bing.com,finance.sina.com.cn,; SINAGLOBAL=117.139.247.210_1737010640.808499; SR_SEL=1_511; U_TRS1=000000d2.a7a4925.6788ca60.a0a405ef; SFA_version8.6.0=2025-02-11%2010%3A18; Apache=117.139.247.210_1739240631.705253; SFA_version8.6.0_click=2; U_TRS2=000000d2.5e28208a4.67aab4ca.2b7491b4; ULV=1739240674514:4:2:2:117.139.247.210_1739240631.705253:1739240654409; FINA_V_S_2=sh601318,sh603300,sh000001,sz002582,sz000725,sz399001; FIN_ALL_VISITED=sh601318%2Csh603300%2Csh000001%2Csz002582%2Csz000725%2Csz399001; rotatecount=1; SFA_version8.7.0=2025-02-18%2015%3A06; SFA_version8.7.0_click=1"
	request, err := http.NewRequest("POST", outUrl, nil)
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")

	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36")

	if err != nil {
		return err
	}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return err
}
