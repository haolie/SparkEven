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
