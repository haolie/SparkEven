package Interface

import (
	"context"
)

type ICodeGather interface {
	StartCodeGather(ctx context.Context) string

	FillCookie(cookie string)

	StopCodeGather(ctx context.Context) string
}
