package Tools

import (
	"context"
	"fmt"
	"testing"
)

func TestLoopCtx(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	i := 0
	LoopCtx(ctx, func() bool {
		i++
		if i >= 10 {
			cancel()
		}

		return i < 100
	})

	if i != 10 {
		t.Errorf("loopCtx fail A")
	}

	ctx, _ = context.WithCancel(context.Background())
	LoopCtx(ctx, func() bool {
		i++
		return i < 100
	})

	if i != 100 {
		fmt.Println(i)
		t.Errorf("loopCtx fail B")
	}
}
