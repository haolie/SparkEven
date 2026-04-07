package Tools

import (
	"context"
	"fmt"
	"testing"
	"time"
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
	}, nil)

	if i != 10 {
		t.Errorf("loopCtx fail A")
	}

	ctx, _ = context.WithCancel(context.Background())
	LoopCtx(ctx, func() bool {
		i++
		return i < 100
	}, nil)

	if i != 100 {
		fmt.Println(i)
		t.Errorf("loopCtx fail B")
	}
}

func TestHunUpAfter(t *testing.T) {
	t1 := time.Now()
	HunUpAfter(context.Background(), 10, "TestHunUpAfter", func() {
		time.Sleep(time.Second * time.Duration(12))
	})

	t2 := time.Now()
	s := t2.Sub(t1).Seconds()
	fmt.Println(s)
}
