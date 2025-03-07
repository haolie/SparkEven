package impl

import (
	"fmt"
	"io"

	"SparkEven/govm/Src/Common/Def"
	"SparkEven/govm/Src/Common/Log"
	"SparkEven/govm/Src/Spark/HServer/HttpTools"
	"SparkEven/govm/Src/Spark/HServer/ginServer"
	"SparkEven/govm/Src/Spark/MPark"
	"github.com/gin-gonic/gin"
)

func init() {
	ginServer.RegisterRoot(Def.Http_key_FillCook, true, fillCook)
}

func fillCook(ctx *gin.Context) {

	cookBye, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		Log.Error(fmt.Sprintf("read body err: %v", err))
	}

	cook := string(cookBye)
	if cook == "" {
		ctx.JSON(200, HttpTools.CreateErrHSResponse("need cook"))
	}

	MPark.CodeGather.FillCookie(cook)

	ctx.JSON(200, HttpTools.CreateSuccessHSResponse("ok"))
}
