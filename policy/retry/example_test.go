package retry_test

import (
	"errors"
	"fmt"
	"time"

	"github.com/Prastiwar/Go-flow/policy/retry"
)

func Example() {
	var errPersistentError = errors.New("persistent-error")

	p := retry.NewPolicy(
		retry.WithCount(2),
		retry.WithWaitTimes(time.Second, 2*time.Second),
	)

	startedTime := time.Now()
	err := p.Execute(func() error {
		fmt.Println("executed after: " + time.Since(startedTime).Truncate(time.Second).String())
		startedTime = time.Now()

		return errPersistentError
	})

	fmt.Println(err == errPersistentError)

	// Output:
	// executed after: 0s
	// executed after: 1s
	// executed after: 2s
	// true
}
