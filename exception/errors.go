package exception

import (
	"fmt"
	"strings"
)

func Aggregate(errors ...error) error {
	count := len(errors)
	if count == 0 {
		return nil
	}

	b := strings.Builder{}

	b.WriteString(errors[0].Error())
	for i := 1; i < count; i++ {
		b.WriteString(", ")
		b.WriteString(errors[i].Error())
	}

	return fmt.Errorf("[%v]", b.String())
}
