package Interface

type IDbModelFill interface {
	FillCodeFace(cb FillCodeCallBack) (model IDbModel, err error)

	GetInsertPerStr() string

	GetInsertValueStr(item IDbModel, noId int) string
}

type FillCodeCallBack func([]interface{}) error
