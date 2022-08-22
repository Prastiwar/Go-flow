package logf

import (
	"encoding/json"
	"fmt"
	"strings"
)

// TextFormatter implements Formatter which outputs message as raw text formaat
// with control over side where fields are put.
type TextFormatter struct {
	leftFields Fields
}

// NewTextFormatter returns a new Text Formatter with Level and LogTime as fields which
// are outputed on the left side of message.
func NewTextFormatter() *TextFormatter {
	return &TextFormatter{
		leftFields: map[string]interface{}{
			Level:   nil,
			LogTime: nil,
		},
	}
}

// NewTextFormatterWith returns a new Text Formatter with field keys which are outputed on the left side of message.
func NewTextFormatterWith(leftFieldNames ...string) *TextFormatter {
	leftFields := make(map[string]interface{}, len(leftFieldNames))
	for _, v := range leftFieldNames {
		leftFields[v] = nil
	}

	return &TextFormatter{
		leftFields: leftFields,
	}
}

// Format returns a formatted string as message "[left fields] message [rest of the fields]".
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
