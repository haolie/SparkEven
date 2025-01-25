package impl

import (
	"fmt"

	"SparkEven/govm/Src/Common/Def"
	"SparkEven/govm/Src/Common/Log"
	"SparkEven/govm/Src/Spark/HServer/ginServer"

	"SparkEven/govm/Src/Spark/HServer/HttpTools"
	"github.com/gin-gonic/gin"
)

func init() {
	ginServer.RegisterRoot(Def.Http_Key_Ctl, codeGatherCtl)
}

var tempStatus = 0

func codeGatherCtl(ctx *gin.Context) {

	status, exists := HttpTools.GetQueryInt(ctx, "start")
	if !exists {
		ctx.JSON(200, HttpTools.CreateErrHSResponse("need start"))
	}

	if status > 0 {
		Log.Info("http ctl start code gather")
	} else {
		Log.Info("http ctl stop code gather")
	}

	if tempStatus != status {
		tempStatus = status
		Log.Debug(fmt.Sprintf("ctl status Changed :%d", tempStatus))
	}

	ctx.JSON(200, HttpTools.CreateSuccessHSResponse(fmt.Sprintf("ctl status Changed :%d", tempStatus)))
}
