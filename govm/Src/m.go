package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"SparkEven/govm/Src/Common/Def"
	"SparkEven/govm/Src/Common/Log"
	"SparkEven/govm/Src/Config"
	"SparkEven/govm/Src/Spark/MPark"

	_ "SparkEven/govm/Src/Spark"
)

func main() {

	err := Config.Load(path.Dir(os.Args[0]))
	if err != nil {
		panic("load config file error:" + err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	//errList, err := MPark.RunLoad(ctx)
	//if err != nil {
	//
	//}

	logPath, exists := Config.GetValue[string](Def.Config_Log_Path)
	if !exists {
		panic("log path not exists")
	}

	Log.InitLog(logPath)

	errList, err := MPark.RunLoad(ctx)
	if err != nil {
		panic(err)
	}

	if len(errList) > 0 {
		Log.Error("start server fail")
		for _, str := range errList {
			Log.Error(str)
		}
	}

	start(ctx)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go waitExit(c, cancel)

	waitEx := ctx.Done()
	<-waitEx
}

func waitExit(signalChan chan os.Signal, cancel context.CancelFunc) {
	for i := range signalChan {
		switch i {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			Log.Info(fmt.Sprintf("exit signal :%v", i))
			cancel()
		}
	}
}

func start(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			Log.Error(fmt.Sprintf("%v", r))
		}
	}()

	isGather, exists := Config.GetValue[bool](Def.Config_Gather_Open)
	if exists && isGather {
		MPark.CodeGather.StartCodeGather(ctx, time.Now())
	}
}

// SELECT tc._no,cf._date,tp.time,tp.price,tp.volume FROM `timeprice2025_1` tp JOIN `codeface` cf ON tp.face_id=cf.id JOIN `tbl_codes` tc ON cf.no_id=tc.id WHERE tc._no=1600699 AND cf._date='2025-05-07'
