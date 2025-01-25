package MPark

import (
	"context"
	"fmt"
)

const (
	runStep_error      = -1
	runStep_unStart    = 0
	runStep_loading    = 1
	runStep_starting   = 2
	runStep_startAfter = 3
	runStep_runing     = 4
)

var (
	status = runStep_error

	loadMap  = make(map[string]ParkCallBack, 4)
	startMap = make(map[string]ParkCallBack, 4)
	afterMap = make(map[string]ParkCallBack, 4)
)

type ParkCallBack func(ctx context.Context) []string
type runStep int32

func GetStatus() int32 {
	return int32(status)
}

func RegisterLoad(key string, cb ParkCallBack) {
	if status >= runStep_loading {
		panic("can not register a Load function after loading")
	}

	_, exists := loadMap[key]
	if exists {
		panic(fmt.Sprintf("can not register a Load function again,key:%s", key))
	}

	loadMap[key] = cb
}

func RegisterStart(key string, cb ParkCallBack) {
	if status >= runStep_starting {
		panic("can not register a start function after start")
	}

	_, exists := startMap[key]
	if exists {
		panic(fmt.Sprintf("can not register a start function again,key:%s", key))
	}

	startMap[key] = cb
}

func RegisterAfterStart(key string, cb ParkCallBack) {
	if status >= runStep_startAfter {
		panic("can not register a start function now")
	}

	_, exists := afterMap[key]
	if exists {
		panic(fmt.Sprintf("can not register a startAfter function again,key:%s", key))
	}

	afterMap[key] = cb
}

func RunLoad(ctx context.Context) (errList []string, err error) {
	status = runStep_loading
	for _, cb := range loadMap {
		tempList := cb(ctx)
		if tempList != nil && len(tempList) > 0 {
			errList = append(errList, tempList...)
		}
	}

	if len(errList) > 0 {
		status = runStep_error
		return errList, nil
	}

	status = runStep_starting
	for _, cb := range startMap {
		tempList := cb(ctx)
		if tempList != nil && len(tempList) > 0 {
			errList = append(errList, tempList...)
		}
	}

	if len(errList) > 0 {
		if len(errList) > 0 {
			status = runStep_error
			return errList, nil
		}
	}

	status = runStep_startAfter
	for _, cb := range afterMap {
		tempList := cb(ctx)
		if tempList != nil && len(tempList) > 0 {
			errList = append(errList, tempList...)
		}
	}

	if len(errList) > 0 {
		if len(errList) > 0 {
			status = runStep_error
			return errList, nil
		}
	}

	status = runStep_runing

	return
}
