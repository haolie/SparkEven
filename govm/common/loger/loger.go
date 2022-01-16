package loger

var (
	isInit = false
)

func Init() error {
	if isInit {
		return nil
	}

	return nil
}

func TypeLog(logType LogType, format string, a ...interface{}) {

}

func WriteLog(outPut int, format string, a ...interface{}) {

}
