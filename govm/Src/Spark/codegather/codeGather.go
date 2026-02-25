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

func init() {
	MPark.CodeGather = NewImpl()
	MPark.RegisterLoad(mName, onLoad)
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
	Log.Info(fmt.Sprintf("start gather startDate:%v", startDate))
	n := time.Now()
	n = time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.Local)
	endDate := n.AddDate(0, 0, 1)
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

	curDateStr := getCurCodeDate()
	for {
		if endDate.Before(date) {
			break
		}

		imp.nextCookTime = imp.startTime
		dateStr := date.Format("2006-01-02")
		faceList, errStr := gatherDateFace(ctx, conn, dateStr)
		if errStr.Exists() {
			tempTimes += 1
			Log.Error(fmt.Sprintf("gatherDateFace faild,date=%v times=%d errStr:%v", date, tempTimes, errStr))
			if tempTimes <= maxErrTimes {
				Tools.Wait(ctx, 5, "gatherDateFaceErr")
				continue
			}
		}

		imp.nextCookTime = imp.startTime.AddDate(0, 0, 1)
		if len(faceList) > 0 && dateStr == curDateStr {
			errStr = startPriceGather(ctx, conn, faceList, dateStr)
			if errStr.Exists() {
				Log.Error(string(errStr))
			}
		}

		tempTimes = 0
		date = date.AddDate(0, 0, 1)
	}

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

func startPriceGather(ctx context.Context, conn *sql.DB, faceList []*Model.CodeFace, dateStr string) (errStr Model.Err) {
	count := len(faceList)
	tryTimes, exists := Config.GetValue[int64](Def.Config_Gather_FailTimes)
	if !exists {
		panic(fmt.Errorf("need %s Config", Def.Config_Gather_FailTimes))
	}

	// 如有采集失败 重复尝试次数
	var failNum int64 = 0
	for ; failNum < tryTimes; failNum++ {
		errMap := make(map[int]Model.Err, 8)
		for i, face := range faceList {
			if face.State == 1 {
				continue
			}

			priceList, faceErr := gatherFacePrice(face)
			if faceErr.Exists() {
				errMap[face.Code] = faceErr
				continue
			}

			if len(priceList) == 0 {
				continue
			}

			faceErr = MPark.DbSupport.SaveFacePrices(conn, face, priceList)
			if faceErr.Exists() {
				errMap[face.Code] = faceErr
				continue
			}

			if i%100 == 0 {
				Log.Info(fmt.Sprintf("gatherFacePrice success code:%d   %d/%d", face.Code, i, count))
			} else {
				Log.Debug(fmt.Sprintf("gatherFacePrice success code:%d   %d/%d", face.Code, i, count))
			}

			Tools.Wait(ctx, con_gatherWait, "startPriceGatherAfter")
		}

		onFinished(errMap, count, dateStr, failNum+1)
		if len(errMap) == 0 {
			break
		}
	}

	return
}
