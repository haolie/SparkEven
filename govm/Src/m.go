package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"SparkEven/govm/Src/Common/Def"
	"SparkEven/govm/Src/Common/Log"
	"SparkEven/govm/Src/Config"
	"SparkEven/govm/Src/Spark/DBSupport"
	"SparkEven/govm/Src/Spark/MPark"

	_ "SparkEven/govm/Src/Spark"
)

func main() {
	Config.Load("")

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

	conn, errStr := DBSupport.CreateConn()
	if errStr != "" {
		Log.Error(errStr)
		return
	}

	errStr = DBSupport.InitSupport(conn)
	if errStr != "" {
		Log.Error(errStr)
		return
	}

	Log.Info("park start success")

	testNo, exists := DBSupport.GetNoById(0)
	if exists {
		Log.Debug(fmt.Sprintf("get id=%d no=%d", 0, testNo))
	}

	faceList, errStr := DBSupport.GetCodeFace(ctx, conn, 1600720, "")
	if errStr != "" {
		Log.Error(errStr)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	Log.Debug(strconv.Itoa(len(faceList)))
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
