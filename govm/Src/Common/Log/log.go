package Log

import (
	"log/slog"
	"os"
)

var (
	//sl = slog.New(slog.NewTextHandler())
	debugLogger *slog.Logger
	infoLogger  *slog.Logger
	wangLogger  *slog.Logger
	errorLogger *slog.Logger
)

type logLvl int

func (lvl logLvl) Level() slog.Level {
	return -10
}

func InitLog(basePath string) {

	debugLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		Level:       logLvl(0),
		ReplaceAttr: nil,
	}))

	infoLogger = slog.New(slog.NewTextHandler(newLogWriter(basePath, "info"), nil))
	wangLogger = slog.New(slog.NewTextHandler(newLogWriter(basePath, "wan"), nil))
	errorLogger = slog.New(slog.NewTextHandler(newLogWriter(basePath, "err"), nil))
}

func Debug(str string) {
	debugLogger.Debug(str)
}

func Info(str string) {
	infoLogger.Info(str)
	debugLogger.Info(str)
}

func Warn(str string) {
	wangLogger.Warn(str)
	debugLogger.Warn(str)
}

func Error(str string) {
	errorLogger.Error(str)
	debugLogger.Error(str)
}
