package impl

import (
	"time"

	"SparkEven/govm/Src/Common/Def"
	"SparkEven/govm/Src/Common/Log"
	"SparkEven/govm/Src/Spark/HServer/HttpTools"
	"SparkEven/govm/Src/Spark/HServer/ginServer"
	"SparkEven/govm/Src/Spark/MPark"
	"github.com/gin-gonic/gin"
)

func init() {
	ginServer.RegisterRoot(Def.Http_key_CodePrice, false, codePrice)
}

func codePrice(ctx *gin.Context) {

	date, exists := ctx.GetQuery("date")
	if !exists || date == "" {
		ctx.JSON(200, HttpTools.CreateErrHSResponse("need date params"))
		return
	}

	startTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		ctx.JSON(200, HttpTools.CreateErrHSResponse("error date params"))
		return
	}

	code, exists := HttpTools.GetQueryInt(ctx, "code")
	if !exists || code == 0 {
		ctx.JSON(200, HttpTools.CreateErrHSResponse("need code params"))
	}

	conn, errStr := MPark.DbSupport.GetConn()
	if errStr.Exists() {
		Log.Error(string(errStr))
		ctx.JSON(200, HttpTools.CreateErrHSResponse(string(errStr)))
		return
	}

	//startTime = startTime.Add(time.Hour * 9)
	priceList, errStr := MPark.DbSupport.GetDateCodePrice(conn, date, code)
	if errStr.Exists() {
		Log.Error(string(errStr))
		ctx.JSON(200, HttpTools.CreateErrHSResponse(string(errStr)))
		return
	}

	list := make([]*timePrice, 0, len(priceList))
	for _, price := range priceList {
		item := &timePrice{
			Time:  startTime.Add(time.Second * time.Duration(price.Time)).Format("15:04:05"),
			Price: price.Price,
		}

		list = append(list, item)
	}

	ctx.JSON(200, HttpTools.CreateSuccessHSResponse(list))
}

type timePrice struct {
	Time  string
	Price int
}
