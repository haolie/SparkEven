package codegather

import (
	"context"
	"fmt"
	"time"

	"SparkEven/govm/Src/Common/Def"
	"SparkEven/govm/Src/Common/Log"
	"SparkEven/govm/Src/Common/Tools"
	"SparkEven/govm/Src/Config"
	"SparkEven/govm/Src/Model"
)

const (
	// 价格采集间隔时间（秒
	con_gatherWait = 1200 * time.Millisecond
)

type impl int32

func (impl impl) StartCodeGather(ctx context.Context, stateDate time.Time) (err Model.Err) {

	go func() {
		startTimeStr, exists := Config.GetValue[string](Def.Config_Gather_StartTime)
		if !exists {
			panic(fmt.Sprintf("need %s Config", Def.Config_Gather_StartTime))
		}

		h, m, s, success := Tools.ParseTimeStr(startTimeStr)
		if !success {
			panic(fmt.Sprintf("parse config %s err", Def.Config_Gather_StartTime))
		}

		now := time.Now()
		startTime := time.Date(now.Year(), now.Month(), now.Day(), h, m, s, 0, time.Local)
		if startTime.Before(now) {
			Log.Info("------start start start -----")
			startGather(ctx, stateDate)
			Log.Info("------end end end -----")
			startTime = startTime.AddDate(0, 0, 1)
		}

		//runTime := now.Add(-time.Duration(now.Minute()%30*60+now.Second()) * time.Second)
		for {
			Log.Info(fmt.Sprintf("wait gatherTime:%v ", startTime))

			select {
			case <-ctx.Done():
				break
			case <-time.After(startTime.Sub(time.Now())):
				Log.Info("------start start start -----")
				startGather(ctx, stateDate)
				Log.Info("------end end end -----")
				startTime = startTime.AddDate(0, 0, 1)
				stateDate = startTime
			}
		}
	}()

	return
}

func (impl impl) FillCookie(cookie string) {
	old := getCookStr()
	if len(old) == 0 {
		fillNewCook(cookie)
	}
}

func onFinished(errMap map[int]Model.Err, total int, dateStr string) {
	Log.Info(fmt.Sprintf("price gather finished date=%s total=%d  err=%d", dateStr, total, len(errMap)))
	if len(errMap) == 0 {
		return
	}

	errInfo := ""
	for code, errStr := range errMap {
		errInfo = fmt.Sprintf("%s %d:%s \r\n,", errInfo, code, errStr)
	}

	Log.Error(errInfo)
}
