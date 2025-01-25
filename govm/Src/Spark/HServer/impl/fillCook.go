package impl

import (
	"SparkEven/govm/Src/Spark/HServer/HttpTools"
	"SparkEven/govm/Src/Spark/MPark"
	"github.com/gin-gonic/gin"
)

func fillCook(ctx *gin.Context) {

	cook := ctx.DefaultQuery("cook", "")
	if cook == "" {
		ctx.JSON(200, HttpTools.CreateErrHSResponse("need cook"))
	}

	MPark.CodeGather.FillCookie(cook)

	ctx.JSON(200, HttpTools.CreateSuccessHSResponse("ok"))
}
