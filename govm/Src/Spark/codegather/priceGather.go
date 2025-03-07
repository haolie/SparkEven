package codegather

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"SparkEven/govm/Src/Common/Tools"
	"SparkEven/govm/Src/Model"
)

const (
	price_change_up   = "UP"
	price_change_down = "DONW"
)

func gatherFacePrice(face *Model.CodeFace) (priceList []*Model.CodePrice, errStr Model.Err) {
	client := &http.Client{}
	var codeStr = Tools.ConvertCodeStr(face.Code)
	url := fmt.Sprintf("https://vip.stock.finance.sina.com.cn/quotes_service/view/CN_TransListV2.php?num=150000&symbol=%s&rn=28950312", codeStr)

	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36")
	if err != nil {
		errStr = Model.Err(fmt.Sprintf("gatherFacePrice  http.NewRequest err:%v", err))
		return
	}

	resp, err := client.Do(request)
	if err != nil {
		errStr = Model.Err(fmt.Sprintf("gatherFacePrice  client.Do err:%v", err))
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		errStr = Model.Err(fmt.Sprintf("gatherFacePrice  io.ReadAll err:%v", err))
		return
	}

	list := strings.Split(string(data), "\n")
	if len(list) < 3 {
		return
	}

	priceList = make([]*Model.CodePrice, 0, len(list)-2)
	for i := len(list) - 1; i >= 0; i-- {
		v := list[i]
		tempIndex := strings.Index(v, "trade_item_list")
		if tempIndex > -1 && len(v) > 40 {
			tempIndex = strings.IndexRune(v, '(')

			valueStr := v[tempIndex+1:]
			valueStr = valueStr[:len(valueStr)-2]
			valueStr = strings.Replace(valueStr, "'", "", -1)
			valueStr = strings.Replace(valueStr, " ", "", -1)
			subValueStr := strings.Split(valueStr, ",")

			tempSeconds, _ := Tools.GetSecondsFromStr(subValueStr[0])
			tempPrice, _ := strconv.ParseFloat(subValueStr[2], 32)
			item := &Model.CodePrice{}
			item.Time = tempSeconds
			item.Volume, _ = strconv.Atoi(subValueStr[1])
			item.Price = int(tempPrice * 100)

			if subValueStr[3] == price_change_up {
				item.TradeType = 1
			} else if subValueStr[3] == price_change_down {
				item.TradeType = -1
			}

			priceList = append(priceList, item)
		}
	}

	if len(priceList) == 0 {
		errStr = Model.Err(fmt.Sprintf("gatherFacePrice code=%d priceList is empty", face.Code))
	}

	return
}
