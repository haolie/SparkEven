package common

// id MEDIUMINT UNSIGNED NOT NULL,
// _date DATE NOT NULL,
//  no_id SMALLINT UNSIGNED NOT NULL,
// _min MEDIUMINT UNSIGNED,
// _max MEDIUMINT UNSIGNED,
// _change MEDIUMINT,
// lastprice MEDIUMINT UNSIGNED,
// startprice MEDIUMINT UNSIGNED,
// dde BIGINT,
// dde_b BIGINT UNSIGNED,
// dde_s BIGINT UNSIGNED,
// face TINYINT UNSIGNED,
// volume BIGINT UNSIGNED,
// turnoverRate SMALLINT UNSIGNED,
// turnover BIGINT UNSIGNED,
// state TINYINT,
// per smallint,

type CodeFace struct {
	ID           int
	Code         int
	Date         string
	MaxValue     float32
	MinValue     float32
	Change       float32
	StartPrice   float32
	LastPrice    float32
	YestPrice    float32
	Dde_b        int
	Dde_s        int
	Face         int
	Volume       int
	Turnover     float32
	TurnoverRate float32
	State        int
	Percent      float32
}

type FaceEx struct {
	CodeFace
	DownTimes int
	FileState int //0 未开始处理  1  已完成处理  2 正在下载  3 已下载  4 数据读取中 5 数据已读取等待保存
}
