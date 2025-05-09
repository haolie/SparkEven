package Def

import (
	"time"

	"SparkEven/govm/Src/Common/Tools"
)

const (
	// 配置 http端口
	Config_http_Port = "http_port"
	Config_FileName  = "Config"

	// 数据库配置
	Config_Db_Ip   = "db_ip"
	Config_Db_Port = "db_port"
	Config_Db_Name = "db_name"
	Config_Db_User = "db_user"
	Config_Db_Pw   = "db_pw"

	Config_Log_Path         = "log_path"
	Config_Gather_StartTime = "gather_starttime"
	Config_Gather_FailTimes = "gather_failtimes"
	Config_Gather_Open      = "gather_open"
	Config_Group_SavePath   = "group_savepath"

	Model_Key_CodeFace = "CodeFace"

	// http 接口
	Http_Key_Ctl       = "ctl"
	Http_key_FillCook  = "fillcook"
	Http_key_CodePrice = "codePrice"
	Http_key_CodeGroup = "codeGroup"

	// 最大成交量
	Max_Volume = 16777215
)

var (
	MinDate   time.Time
	MaxDate   time.Time
	StartTick int
)

func init() {

	StartTick, _ = Tools.GetSecondsFromStr("09:00:00")

	t := time.Now()
	MinDate = time.Date(2000, 1, 1, 0, 0, 0, 0, t.Location())
	MaxDate = time.Date(2077, 1, 1, 0, 0, 0, 0, t.Location())
}
