package Interface

import (
	"time"

	"SparkEven/govm/Src/Model"
)

type IDbSupport interface {
	SearchCodeFace(date time.Time, code int) (list []*Model.CodeFace, err error)
	SaveFaceList([]*Model.CodeFace) error
}
