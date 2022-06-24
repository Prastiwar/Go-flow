package logf

import (
	"encoding/json"
	"fmt"
	"strings"
)

type TextFormatter struct {
	leftFields Fields
}

func NewTextFormatter() *TextFormatter {
	return &TextFormatter{
		leftFields: map[string]interface{}{
			Level:   nil,
			LogTime: nil,
		},
	}
}

func NewTextFormatterWith(leftFieldNames ...string) *TextFormatter {
	leftFields := make(map[string]interface{}, len(leftFieldNames))
	for _, v := range leftFieldNames {
		leftFields[v] = nil
	}

	return &TextFormatter{
		leftFields: leftFields,
	}
}

func (f *TextFormatter) Format(msg string, fields Fields) string {
	leftBuilder := strings.Builder{}
	rightBuilder := strings.Builder{}
	rightFields := make(map[string]interface{}, len(fields))

	for k, v := range fields {
		_, isLeft := f.leftFields[k]
		if isLeft {
			leftBuilder.WriteRune('[')
			leftBuilder.WriteString(fmt.Sprintf("%v", v))
			leftBuilder.WriteString("] ")
		} else {
			rightFields[k] = v
		}
	}

	if len(rightFields) > 0 {
		rightBuilder.WriteRune(' ')
		jsonBytes, _ := json.Marshal(rightFields)
		rightBuilder.Write(jsonBytes)
	}
	rightBuilder.WriteRune('\n')

	return leftBuilder.String() + msg + rightBuilder.String()
}
