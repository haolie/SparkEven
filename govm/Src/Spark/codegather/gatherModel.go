package codegather

type Col struct {
	name      string
	index     int
	id        string
	dateCheck bool
}

type Cl struct {
	Uintn         string `json:"uint"`
	Index_name    string `json:"index_name"`
	Key           string `json:"key"`
	Timestamp     string `json:"timestamp"`
	RealIndexName string `json:"realIndexName"`
}

type tokenObj struct {
	count   int
	token   string
	columns []Cl
}

const (
	// 数据列表名 涨跌
	Col_Ud = "ud"
	// 数据列表名 收盘价
	Col_LastPrice = "lastprice"
	// 数据列表名 开盘价
	Col_StartPrice = "startprice"
	// 数据列表名 最高价
	Col_Max = "max"
	// 数据列表名 最低价
	Col_Min = "min"
	// 数据列表名 换手率
	Col_TurnoverRate = "turnoverRate"
	// 数据列表名 成交额
	Col_Turnover = "turnover"
	// 数据列表名 成交量
	Col_Volume = "volume"
)
