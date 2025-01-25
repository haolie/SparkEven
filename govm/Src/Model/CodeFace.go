package Model

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

func (codeFace *CodeFace) CreateInsertStr() string {
	return ""
}

func (codeFace *CodeFace) GetModelKey() string {
	return ""
}

type FaceEx struct {
	CodeFace
	DownTimes int
	FileState int //0 未开始处理  1  已完成处理  2 正在下载  3 已下载  4 数据读取中 5 数据已读取等待保存  8 数据错误
}
