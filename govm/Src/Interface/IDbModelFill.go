package Interface

type IDbModelFill interface {
	FillCodeFace(cb FillCodeCallBack) (model IDbModel, err error)
}

type FillCodeCallBack func([]interface{}) error
