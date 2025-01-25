package ginServer

import (
	"context"
	"fmt"
	"sync"

	"SparkEven/govm/Src/Common/Def"
	"SparkEven/govm/Src/Common/Log"
	"SparkEven/govm/Src/Config"
	"SparkEven/govm/Src/Spark/MPark"
	"github.com/gin-gonic/gin"
)

const (
	con_fileName = "ginSerfer"
)

func init() {
	MPark.RegisterStart("httpServer", func(ctx context.Context) (errList []string) {
		errStr := startHServer(ctx)
		if errStr != "" {
			errList = append(errList, errStr)
		}

		return
	})
}

var (
	startOnce sync.Once
	rootMap   = make(map[string]func(c *gin.Context), 8)
	status    int32
)

func RegisterRoot(key string, f func(c *gin.Context)) {
	if status > 0 {
		Log.Error("can not register root after started")
		return
	}

	_, exist := rootMap[key]
	if exist {
		Log.Error(fmt.Sprintf("%s:k={%s} already exist", con_fileName, key))
		return
	}

	rootMap[key] = f
}

func startHServer(ctx context.Context) (errStr string) {

	startOnce.Do(func() {
		go func() {
			port, exists := Config.GetValue[int64](Def.Config_http_Port)
			if !exists {
				panic("config http port not exist")
			}

			engine := gin.Default()
			for k, v := range rootMap {
				str := "/express/h/" + k
				engine.GET(str, v)
			}

			err := engine.Run(fmt.Sprintf(":%d", port))
			if err != nil {
				errStr = fmt.Sprintf("%s fail to start http server err:%s", con_fileName, err.Error())
			}
		}()

	})

	return
}
