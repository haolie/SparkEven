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

var (
	startH int
	startM int
	startS int
)

func loadStartTime() string {
	startTimeStr, exists := Config.GetValue[string](Def.Config_Gather_StartTime)
	if !exists {
		return fmt.Sprintf("need %s Config", Def.Config_Gather_StartTime)
	}

	var success bool
	startH, startM, startS, success = Tools.ParseTimeStr(startTimeStr)
	if !success {
		return fmt.Sprintf("parse config %s err", Def.Config_Gather_StartTime)
	}

	return ""
}

type impl struct {
	startTime    time.Time
	nextCookTime time.Time
}

func NewImpl() *impl {
	return &impl{
		startTime: time.Now().Add(time.Minute),
	}
}

func getStartTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), startH, startM, startS, 0, time.Local)
}

func (impl *impl) StartCodeGather(ctx context.Context, stateDate time.Time) (err Model.Err) {

	go func() {

		now := time.Now()
		impl.startTime = getStartTime(now)
		impl.nextCookTime = impl.startTime
		if impl.startTime.Before(now) {
			Log.Info("------start start start -----")
			startGather(ctx, stateDate, impl)
			Log.Info("------end end end -----")
			impl.startTime = impl.startTime.AddDate(0, 0, 1)
			impl.nextCookTime = impl.startTime
		}

		//runTime := now.Add(-time.Duration(now.Minute()%30*60+now.Second()) * time.Second)
		for {
			Log.Info(fmt.Sprintf("wait gatherTime:%v ", impl.startTime))

			select {
			case <-ctx.Done():
				break
			case <-time.After(impl.startTime.Sub(time.Now())):
				Log.Info("------start start start -----")
				tx, _ := context.WithTimeout(ctx, time.Hour*7)
				startGather(tx, stateDate, impl)
				Log.Info("------end end end -----")
				impl.startTime = impl.startTime.AddDate(0, 0, 1)
				impl.nextCookTime = impl.startTime
				stateDate = impl.startTime
			}
		}
	}()

	return
}

func (impl *impl) FillCookie(cookie string) int {
	old := getCookStr()
	if len(old) == 0 {
		fillNewCook(cookie)
	}

	n := time.Now()
	if impl.nextCookTime.After(n) {
		return int(impl.nextCookTime.Sub(n).Seconds())
	} else {
		return 20
	}
}

func (impl impl) CodeGroup(date string) (err Model.Err) {
	savePath, exists := Config.GetValue[string](Def.Config_Group_SavePath)
	if !exists {
		return
	}

	group := CreateGroupper(date, savePath, 400)
	group.Start(context.Background())

	return
}

func onFinished(errMap map[int]Model.Err, total int, dateStr string, tryNum int64) {
	Log.Info(fmt.Sprintf("price gather finished date=%s total=%d  err=%d  tryNum=%v", dateStr, total, len(errMap), tryNum))
	if len(errMap) == 0 {
		return
	}

	errInfo := ""
	for code, errStr := range errMap {
		errInfo = fmt.Sprintf("%s %d:%s \r\n,", errInfo, code, errStr)
	}

	Log.Error(errInfo)
}
