package codegather

import (
	"context"
	"fmt"
	"sync"
	"time"

	"SparkEven/govm/Src/Common/Def"
	"SparkEven/govm/Src/Common/Log"
	"SparkEven/govm/Src/Model"
)

const (
	cookExpireSeconds = 60
	cookUseSeconds    = 50
)

var (
	cookStr string

	newCookStr string

	cookExpireTime time.Time

	cookUseTimes int

	cookLocker sync.RWMutex
)

func waitCookStr(ctx context.Context) (cookieStr string, er Model.Err) {
	for {
		select {
		case <-ctx.Done():
			er = "ctx cancel"
			return
		default:
			cookieStr = getCookStr()
			if cookieStr != "" {
				return
			}

			Log.Warn(fmt.Sprintf("waitCookStr cookStr:%s  newStr:%s  exprieTime:%v", cookStr, newCookStr, cookExpireTime))
			time.Sleep(time.Second)
		}
	}
}

func getCookStr() string {

	getFun := func() (cStr string, existNew bool) {
		cookLocker.RLock()
		defer cookLocker.RUnlock()
		if cookStr != "" && cookExpireTime.Before(time.Now()) {
			cStr = cookStr
			return
		}

		if newCookStr == "" {
			cStr = ""
		} else {
			existNew = true
		}

		return
	}

	result, existNew := getFun()
	if result != "" || !existNew {
		return result
	}

	cookLocker.Lock()
	defer cookLocker.Unlock()

	if cookStr == "" {
		cookStr = newCookStr
		cookExpireTime = time.Now().Add(cookExpireSeconds * time.Second)
	}

	return cookStr

}

func resetCook(isClear bool) {
	cookLocker.Lock()
	defer cookLocker.Unlock()

	if isClear {
		cookStr = ""
		newCookStr = ""
		return
	}

	if newCookStr != "" {
		cookStr = newCookStr
		cookExpireTime = time.Now().Add(cookExpireSeconds * time.Second)
	} else {
		cookStr = ""
		cookExpireTime = Def.MinDate
	}
}

func fillNewCook(setStr string) {
	cookLocker.Lock()
	defer cookLocker.Unlock()

	if cookStr == "" {
		cookStr = setStr
		newCookStr = ""
		cookExpireTime = time.Now().Add(cookExpireSeconds * time.Second)
	} else {
		newCookStr = cookStr
	}
}

func needNewCook(cook string) bool {
	cookLocker.RLock()
	defer cookLocker.RUnlock()

	return cookExpireTime.Sub(time.Now()) < time.Duration(cookUseSeconds)*time.Second
}
