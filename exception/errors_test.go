package exception_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/Prastiwar/Go-flow/exception"
	"github.com/Prastiwar/Go-flow/tests/assert"
)

func TestAggregate(t *testing.T) {
	tests := []struct {
		name string
		errs []error
		want exception.AggregatedError
	}{
		{
			name: "success-nested-flat",
			errs: []error{
				errors.New("1"),
				exception.Aggregate(errors.New("2"), errors.New("3"), errors.New("4")),
				exception.Aggregate(errors.New("5")),
				errors.New("6"),
			},
			want: exception.Aggregate(
				errors.New("1"),
				errors.New("2"),
				errors.New("3"),
				errors.New("4"),
				errors.New("5"),
				errors.New("6"),
			),
		},
		{
			name: "success-flat-flat",
			errs: []error{
				errors.New("1"),
				errors.New("2"),
				errors.New("3"),
			},
			want: exception.Aggregate(
				errors.New("1"),
				errors.New("2"),
				errors.New("3"),
			),
		},
		{
			name: "success-shuffled",
			errs: []error{
				exception.Aggregate(errors.New("2"), errors.New("3"), errors.New("4")),
				errors.New("1"),
				errors.New("6"),
				exception.Aggregate(errors.New("5")),
			},
			want: exception.Aggregate(
				errors.New("2"),
				errors.New("3"),
				errors.New("4"),
				errors.New("1"),
				errors.New("6"),
				errors.New("5"),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := exception.Aggregate(tt.errs...).Flat()
			assert.Equal(t, tt.want.Error(), got.Error())
		})
	}
}

func TestAggregatef(t *testing.T) {
	tests := []struct {
		name    string
		errors  []error
		wantErr error
	}{
		{
			name:    "nil",
			wantErr: nil,
		},
		{
			name: "errors",
			errors: []error{
				errors.New("test"),
				errors.New("test2"),
			},
			wantErr: errors.New("[\"test\", \"test2\"]"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := exception.Aggregatedf(tt.errors...)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestStackTrace(t *testing.T) {
	got := exception.StackTrace()

	lines := strings.Split(got, "\n")
	if len(lines) < 6 {
		t.Error("too few lines")
	}

	const expected = "exception_test.TestStackTrace"
	contains := false
	for _, v := range lines {
		if strings.Contains(v, expected) {
			contains = true
			break
		}
	}
	assert.Equal(t, true, contains, "stack trace does not contain current function path")
}
