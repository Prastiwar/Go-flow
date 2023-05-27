package rate_test

import (
	"context"
	"time"

	"github.com/Prastiwar/Go-flow/rate"
	"github.com/Prastiwar/Go-flow/rate/memory"
	"github.com/Prastiwar/Go-flow/rate/slidingwindow"
)

func Example() {
	sw, err := slidingwindow.NewAlgorithm(40, time.Minute, 20)
	if err != nil {
		panic(err)
	}

	store, err := memory.NewLimiterStore(context.Background(), sw, time.Hour)
	if err != nil {
		panic(err)
	}

	l, err := store.Limit(context.Background(), "{id}")
	if err != nil {
		panic(err)
	}

	t, err := l.Take(context.Background())
	if err != nil {
		panic(err)
	}

	if err := t.Use(); err != nil {
		panic(err)
	}

	if err := rate.Wait(context.Background(), t.ResetsAt()); err != nil {
		panic(err)
	}

	for i := 0; i < 15; i++ {
		if err := rate.ConsumeAndWait(context.Background(), l); err != nil {
			panic(err)
		}
	}
}
