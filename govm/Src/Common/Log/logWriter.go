package Log

import (
	"fmt"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"time"
)

type logWriter struct {
	logEndTime  int64
	locker      sync.Mutex
	fs          *os.File
	baseLogPath string
	typeName    string
}

func newLogWriter(logPath string, typeName string) *logWriter {
	l := &logWriter{
		baseLogPath: logPath,
		locker:      sync.Mutex{},
		logEndTime:  0,
		typeName:    typeName,
	}

	return l
}

func (writer *logWriter) Write(p []byte) (n int, err error) {
	err = writer.newWriter()
	if err != nil {
		return
	}

	n, err = writer.fs.Write(p)
	return
}

func (writer *logWriter) Close() error {
	writer.locker.Lock()
	defer writer.locker.Unlock()

	if writer.logEndTime == 0 {
		return nil
	}

	err := writer.fs.Close()
	if err != nil {
		return err
	}

	writer.logEndTime = 0
	writer.fs = nil

	return nil
}

func (writer *logWriter) newWriter() error {
	logEndTime := atomic.LoadInt64(&writer.logEndTime)
	if logEndTime < 0 {
		return fmt.Errorf("loger closed")
	}

	t := time.Now()
	tk := t.Unix()
	if tk < logEndTime {
		return nil
	}

	writer.locker.Lock()
	defer writer.locker.Unlock()

	if writer.logEndTime < 0 {
		return fmt.Errorf("loger closed")
	}

	if tk < writer.logEndTime {
		return nil
	}

	if writer.fs != nil {
		err := writer.fs.Close()
		if err != nil {
			return err
		}
	}

	y, m, d := t.Date()
	logPath := fmt.Sprintf("%d/%d", y, m)
	logPath = path.Join(writer.baseLogPath, logPath)

	err := os.MkdirAll(logPath, os.ModePerm)
	if err != nil {
		return err
	}

	logFile := path.Join(logPath, fmt.Sprintf("%d-%d-%d_%d_%s.txt", y, m, d, t.Hour(), writer.typeName))
	writer.fs, err = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	_, mi, s := t.Clock()
	s = mi*60 + s
	t = t.Add(-time.Second * time.Duration(s)).Add(time.Hour)
	writer.logEndTime = t.Unix()

	return nil
}
