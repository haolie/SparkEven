package dbFillImpl

import (
	"SparkEven/govm/Src/Interface"
	"SparkEven/govm/Src/Model"
)

type CodeFaceImpl struct{}

func (face *CodeFaceImpl) GetSelectStr() string {
	return ""
}

func (face *CodeFaceImpl) FillCodeFace(cb Interface.FillCodeCallBack) (model Interface.IDbModel, err error) {
	faceItem := &Model.CodeFace{}
	var noId, minValue, maxValue, change, startPrice, lastPrice, turnover, turnoverRate, percent int
	inList := make([]interface{}, 0, 16)
	inList = append(inList, &faceItem.ID,
		&faceItem.Date,
		&noId,
		&minValue,
		&maxValue,
		&change,
		&lastPrice,
		&startPrice,
		&faceItem.Dde_b,
		&faceItem.Dde_s,
		&faceItem.Face,
		&faceItem.Volume,
		&turnover,
		&turnoverRate,
		&faceItem.State,
		&percent)

	err = cb(inList)
	if err != nil {
		return
	}

	//face.Code = this.GetNoById(noId)
	faceItem.MinValue = float32(minValue) / 100
	faceItem.MaxValue = float32(maxValue) / 100
	faceItem.Change = float32(change) / 100
	faceItem.LastPrice = float32(lastPrice) / 100
	faceItem.StartPrice = float32(startPrice) / 100
	faceItem.Turnover = float32(turnover) / 100
	faceItem.TurnoverRate = float32(turnoverRate) / 100
	faceItem.Percent = float32(percent) / 100
	faceItem.YestPrice = faceItem.LastPrice - faceItem.Change

	model = faceItem
	return
}
