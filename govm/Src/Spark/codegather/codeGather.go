package codegather

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"SparkEven/govm/Src/Common/Def"
	"SparkEven/govm/Src/Common/Log"
	"SparkEven/govm/Src/Common/Tools"
	"SparkEven/govm/Src/Config"
	"SparkEven/govm/Src/Model"
	"SparkEven/govm/Src/Spark/MPark"
)

const (
	mName = "codeGather"
)

// 采集最早日期
var minDate time.Time

func init() {
	MPark.CodeGather = NewImpl()
	MPark.RegisterLoad(mName, onLoad)

	minDate = time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
}

func onLoad(ctx context.Context) []string {
	errList := make([]string, 0)
	errStr := loadStartTime()
	if errStr != "" {
		errList = append(errList, errStr)
	}

	return errList
}

func startGather(ctx context.Context, startDate time.Time, imp *impl) (errStr Model.Err) {
	resetCook(true)

	n := time.Now()
	n = Tools.GetDate(n)
	date := startDate
	var tempTimes int64
	maxErrTimes, exists := Config.GetValue[int64](Def.Config_Gather_FailTimes)
	if !exists {
		panic(fmt.Sprintf("need %s Config", Def.Config_Gather_FailTimes))
	}

	conn, errStr := MPark.DbSupport.GetConn()
	if errStr.Exists() {
		Log.Error(string(errStr))
		return errStr
	}

	defer conn.Close()

	// 获取当日交易日期
	curDateStr := getCurCodeDate()
	Tools.LoopCtx(ctx, func() bool {
		if minDate.After(date) {
			return false
		}

		imp.nextCookTime = imp.startTime
		dateStr := date.Format("2006-01-02")
		Log.Info(fmt.Sprintf("start gather startDate:%v", dateStr))
		faceList, errStr := gatherDateFace(ctx, conn, dateStr)
		if errStr.Exists() {
			tempTimes += 1
			Log.Error(fmt.Sprintf("gatherDateFace faild,date=%v times=%d errStr:%v", date, tempTimes, errStr))
			if tempTimes <= maxErrTimes {
				Tools.Wait(ctx, 5, "gatherDateFaceErr")
				return true
			}
		}

		if len(faceList) > 0 && dateStr == curDateStr {
			go startPriceGather(ctx, faceList, dateStr)
		}

		return false
		//tempTimes = 0
		//date = date.AddDate(0, 0, 1)
		//return true
	}, func() {
		Log.Info(fmt.Sprintf("startGather ctx.Down"))
	})

	imp.nextCookTime = imp.startTime.AddDate(0, 0, 1)
	return errStr
}

func gatherDateFace(ctx context.Context, conn *sql.DB, dateStr string) (faceList []*Model.CodeFace, errStr Model.Err) {
	faceList, errStr = MPark.DbSupport.SearchCodeFace(conn, dateStr, -1)
	if errStr.Exists() {
		return
	}

	if len(faceList) > 0 {
		return
	}

	//ctx, _ = context.WithTimeout(ctx, 3*time.Minute)
	faceList, errStr = GetNocodesFromWeb(ctx, dateStr)
	if !errStr.Exists() {
		errStr = MPark.DbSupport.SaveFaceList(conn, faceList)
	}

	Tools.Wait(ctx, 2, "gatherDateFace")
	return faceList, errStr
}

func startPriceGather(ctx context.Context, faceList []*Model.CodeFace, dateStr string) {
	count := len(faceList)
	tryTimes, exists := Config.GetValue[int64](Def.Config_Gather_FailTimes)
	if !exists {
		panic(fmt.Errorf("need %s Config", Def.Config_Gather_FailTimes))
	}

	conn, errStr := MPark.DbSupport.GetConn()
	if errStr.Exists() {
		Log.Error(string(errStr))
	}

	defer conn.Close()

	// 如有采集失败 重复尝试次数
	var failNum int64 = 0
	var index int
	errMap := make(map[int]Model.Err, 8)
	Tools.LoopCtx(ctx, func() bool {
		// 超失败次数 || 采集结束
		if failNum >= tryTimes || index >= count {
			onFinished(errMap, count, dateStr, failNum+1)
			return false
		}

		curIndex := index
		index += 1
		face := faceList[curIndex]
		// 跳过已采集
		if face.State == 1 {
			return true
		}

		// 采集价格
		priceList, faceErr := gatherFacePrice(face)
		if faceErr.Exists() {
			errMap[face.Code] = faceErr
			return true
		}

		if len(priceList) == 0 {
			return true
		}

		// 保存数据库
		faceErr = MPark.DbSupport.SaveFacePrices(conn, face, priceList)
		if faceErr.Exists() {
			errMap[face.Code] = faceErr
			return true
		}

		if curIndex%100 == 0 {
			Log.Info(fmt.Sprintf("gatherFacePrice success code:%d   %d/%d", face.Code, curIndex, count))
		} else {
			Log.Debug(fmt.Sprintf("gatherFacePrice success code:%d   %d/%d", face.Code, curIndex, count))
		}

		Tools.Wait(ctx, con_gatherWait, "startPriceGatherAfter")
		return true
	}, func() {
		Log.Info(fmt.Sprintf("startPriceGather ctx.Down"))
	})

	return
}
