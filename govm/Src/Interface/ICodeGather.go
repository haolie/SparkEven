package Interface

import (
	"context"
	"time"

	"SparkEven/govm/Src/Model"
)

type ICodeGather interface {
	StartCodeGather(ctx context.Context, stateDate time.Time) (err Model.Err)

	FillCookie(cookie string)

	CodeGroup(date string) (err Model.Err)
}
