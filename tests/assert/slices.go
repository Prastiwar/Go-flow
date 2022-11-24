package assert

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

// ElementsMatch asserts both slices have the same amount and equal elements. Any extra items
// will be listed in error log.
func ElementsMatch[T any](t *testing.T, arrA, arrB []T, prefixes ...string) {
	t.Helper()
	extraA, extraB := diffLists(arrA, arrB)

	if len(extraA) == 0 && len(extraB) == 0 {
		return
	}

	msg := formatArrDiff(arrA, arrB, extraA, extraB)

	errorf(t, msg, prefixes...)
}

// diffLists diffs two arrays/slices and returns slices of elements that are only in A and only in B.
// If some element is present multiple times, each instance is counted separately (e.g. if something is 2x in A and
// 5x in B, it will be 0x in extraA and 3x in extraB). The order of items in both lists is ignored.
func diffLists[T any](arrA, arrB []T) (extraA, extraB []T) {
	aLen := len(arrA)
	bLen := len(arrB)

	visited := make([]bool, bLen)
	for i := 0; i < aLen; i++ {
		element := arrA[i]
		found := false
		for j := 0; j < bLen; j++ {
			if visited[j] {
				continue
			}
			if reflect.DeepEqual(arrB[j], element) {
				visited[j] = true
				found = true
				break
			}
		}
		if !found {
			extraA = append(extraA, element)
		}
	}

	for j := 0; j < bLen; j++ {
		if visited[j] {
			continue
		}
		extraB = append(extraB, arrB[j])
	}

	return
}

func formatArrDiff[T any](arrA, arrB, extraA, extraB []T) string {
	var msg bytes.Buffer

	msg.WriteString("elements differ")
	if len(extraA) > 0 {
		msg.WriteString("\nextra elements in list A:\n")
		for _, v := range extraA {
			msg.WriteString(fmt.Sprintf("%#v\n", v))
		}
	}
	if len(extraB) > 0 {
		msg.WriteString("\nextra elements in list B:\n")
		for _, v := range extraB {
			msg.WriteString(fmt.Sprintf("%#v\n", v))
		}
	}

	msg.WriteString("\narray A:\n")
	for _, v := range arrA {
		msg.WriteString(fmt.Sprintf("%#v\n", v))
	}

	msg.WriteString("\narray B:\n")
	for _, v := range arrB {
		msg.WriteString(fmt.Sprintf("%#v\n", v))
	}

	return msg.String()
}
