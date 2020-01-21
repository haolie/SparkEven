package nohelper

import (
	"fmt"
	"net/http"
	"sync"
	"io/ioutil"
	"strings"
	"time"
	"reflect"
	"strconv"
	"github.com/widuu/gojson"
)

var(
	newSession=""
	lock sync.Mutex
)



func StartListenSession(){
	mux:=http.NewServeMux()
	mux.HandleFunc("/",func (res http.ResponseWriter,r *http.Request)  {
		lock.Lock()
		defer lock.Unlock()
		newSession=strings.TrimLeft(r.URL.Path,"/") 
		//fmt.Println(newSession)
		fmt.Fprintf(res, "%v\n", time.Now().Format("2006-01-02 15:04:05"))
	})
	if err:=http.ListenAndServe(":9001",mux);err!=nil{
		fmt.Println(err)
	}
	fmt.Println("codeId  监听中……")
}

func GetToken(date string)(*gojson.Js,bool){

	for{
		if len(newSession)>0{
			break
		}
		time.Sleep(time.Second*2)
	}

	client:=&http.Client{}
	var url= "http://www.iwencai.com/stockpick/load-data?typed=0&preParams=&ts=1&f=1&qs=result_original&selfsectsn=&querytype=stock&searchfilter=&tid=stockpick&w=%E6%B6%A8%E8%B7%8C" +
	date +
	"+%E4%BB%B7%E6%A0%BC" +
	date +
	"+%E6%88%90%E4%BA%A4%E9%87%8F" +
	date +
	"+%E6%8D%A2%E6%89%8B" +
	date +
	"+&queryarea="
	fmt.Println(url)
	request,err:=http.NewRequest("GET",url,nil)
	request.Header.Add("Connection","keep-alive")
	request.Header.Add("Host","www.iwencai.com")
	request.Header.Add("User-Agent","Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.75 Safari/537.36")
	request.Header.Add("Cookie",newSession)
	resp,err:=client.Do(request)
	defer resp.Body.Close()
	if err!=nil{
		return nil,false
	}
	result,err:=ioutil.ReadAll(resp.Body)
	if err!=nil{
		return nil,false
	}
	c1 := gojson.Json(string(result)).Get("data").Get("result") 
	return c1,true
}

type Col struct{
	name string
	index int
	id string
}

type CodeBase struct{
	date string
	state int
	index int

}

func getNo(token string){

}

func createColMap() map[string]Col {
	m:=make(map[string]Col)
	m["ud"]=Col{"涨跌",-1,"ud"}
	m["lastprice"]=Col{"收盘价:不复权",-1,"lastprice"}
	m["startprice"]=Col{"开盘价:不复权",-1,"startprice"}
	m["max"]=Col{"最高价:不复权",-1,"max"}
	m["min"]=Col{"涨最低价:不复权跌",-1,"min"}
	m["turnoverRate"]=Col{"换手率",-1,"turnoverRate"}
	m["turnover"]=Col{"成交额",-1,"turnover"}
	m["volume"]=Col{"成交量",-1,"volume"}
	return m
}

func GetNocodesFromWeb(date string){
	tokenObj,_:=GetToken(date)
	// perCount,index,pageCount,allCodes:=1000,1,0,make([]string,3000)
	// codeCount,_:=strconv.Atoi( tokenObj.Get("code_count").Tostring())
	// pageCount=codeCount/perCount
	// if codeCount%pageCount>0{
	// 	pageCount++
	// }
	// ms:=createColMap()
	strconv.Atoi("1")
	//_,cls2:= tokenObj.Get("columnsIndexID").ToArray()

	fmt.Println(reflect.TypeOf(tokenObj.Getdata()["columnsIndexID"]))
	//tokenObj.Get("columnsIndexID").Type()
	//fmt.Println(tokenObj.Get("columnsIndexID").Type())
	// for i,v:=range tokenObj.Getdata()["columnsIndexID"]{
	// 	fmt.Println(i,v)
	// }
     gp :=tokenObj.Tostring()
//	var ays []interface {}=tokenObj.Get("columnsIndexID").data
	fmt.Println(gp)

	//fmt.Println(cls2)



}