package loger

const (
	StdOut = 1 << iota
	File
)

type LogType int32

const (
	LogTypeEnum_Info LogType = 1 + iota

	LogTypeEnum_Warning

	LogTypeEnum_Err
)
