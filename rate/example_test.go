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

	store := memory.NewLimiterStore(context.Background(), sw, time.Hour)

	l := store.Limit("{id}")

	t := l.Take()
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
