package DBSupport

import (
	"fmt"

	"SparkEven/govm/Src/Common/Def"
	"SparkEven/govm/Src/Interface"
	"SparkEven/govm/Src/Model"
)

const (
	con_insert_face = "INSERT delayed INTO codeface(id,_date,no_id,_min,_max,_change,lastprice,startprice,volume,turnoverRate,turnover,face,dde,dde_b,dde_s,state,per) VALUES"

	Mysql_smUInt int = 65535
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

	faceItem.Code, _ = GetNoById(noId)
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

func (face *CodeFaceImpl) GetInsertPerStr() string {
	return con_insert_face
}

func (face *CodeFaceImpl) GetInsertValueStr(item Interface.IDbModel, noId int) string {
	codeface := item.(*Model.CodeFace)
	if codeface.Volume > Def.Max_Volume {
		codeface.Volume = Def.Max_Volume
	}

	turnoverRate := int(codeface.TurnoverRate * 100)
	if turnoverRate > Mysql_smUInt {
		turnoverRate = Mysql_smUInt
	}

	tempStr := fmt.Sprintf("(%d,'%s',%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d)",
		codeface.ID,
		codeface.Date,
		noId,
		int(codeface.MinValue*100),
		int(codeface.MaxValue*100),
		int(codeface.Change*100),
		int(codeface.LastPrice*100),
		int(codeface.StartPrice*100),
		codeface.Volume,
		turnoverRate,
		int(codeface.Turnover*100),
		codeface.Face,
		0,
		codeface.Dde_b,
		codeface.Dde_s,
		codeface.State,
		int(codeface.Percent*100))

	return tempStr
}
