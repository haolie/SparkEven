package impl

import (
	"SparkEven/govm/Src/Common/Def"
	"SparkEven/govm/Src/Spark/HServer/HttpTools"
	"SparkEven/govm/Src/Spark/HServer/ginServer"
	"SparkEven/govm/Src/Spark/MPark"
	"github.com/gin-gonic/gin"
)

func init() {
	ginServer.RegisterRoot(Def.Http_key_CodeGroup, false, codeGroup)
}

func codeGroup(ctx *gin.Context) {

	date, exists := ctx.GetQuery("date")
	if !exists || date == "" {
		ctx.JSON(200, HttpTools.CreateErrHSResponse("need date params"))
		return
	}

	MPark.CodeGather.CodeGroup(date)

	ctx.JSON(200, HttpTools.CreateSuccessHSResponse("success"))
}
