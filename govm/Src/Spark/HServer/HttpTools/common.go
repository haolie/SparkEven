package HttpTools

import (
	"strconv"
	"time"

	"SparkEven/govm/Src/Common/Tools"
	"github.com/gin-gonic/gin"
)

const (
	Con_HS_Status_Fail    = 0
	Con_HS_Status_Success = 200
)

func baseDataMap() map[string]interface{} {
	dataMap := make(map[string]interface{}, 5)
	dataMap["status"] = Con_HS_Status_Success
	dataMap["content"] = ""
	dataMap["data"] = "11111"
	dataMap["time"] = time.Now().Unix()
	dataMap["timeStr"] = Tools.ToDateTimeStr(time.Now())

	return dataMap
}

func CreateErrHSResponse(errInfo string) map[string]interface{} {
	dataMap := baseDataMap()
	dataMap["status"] = Con_HS_Status_Fail
	dataMap["content"] = errInfo

	return dataMap
}

func CreateSuccessHSResponse(info interface{}) map[string]interface{} {

	dataMap := baseDataMap()
	dataMap["status"] = Con_HS_Status_Success
	dataMap["content"] = info

	return dataMap
}

func GetQueryInt32(c *gin.Context, key string) (q int32, exists bool) {
	qs, exists := c.GetQuery(key)
	if !exists {
		return
	}

	temp, err := strconv.Atoi(qs)
	if err != nil {
		exists = false
		return
	}

	return int32(temp), true
}

func GetQueryInt(c *gin.Context, key string) (q int, exists bool) {
	qs, exists := c.GetQuery(key)
	if !exists {
		return
	}

	temp, err := strconv.Atoi(qs)
	if err != nil {
		exists = false
		return
	}

	return temp, true
}

func GetQueryBool(c *gin.Context, key string) (q bool, exists bool) {
	qs, exists := c.GetQuery(key)
	if !exists {
		return
	}

	q = qs == "true"

	return
}
