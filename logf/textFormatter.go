package logf

import (
	"encoding/json"
	"fmt"
	"strings"
)

// TextFormatter implements Formatter which outputs message as raw text formaat
// with control over side where fields are put.
type TextFormatter struct {
	leftFieldKeys []string
}

// NewTextFormatter returns a new Text Formatter with Level and LogTime as field keys which
// values will be outputted on the left side of message.
func NewTextFormatter() *TextFormatter {
	return &TextFormatter{
		leftFieldKeys: []string{Level, LogTime},
	}
}

// NewTextFormatterWith returns a new Text Formatter with field keys which are outputted on the left side of message.
func NewTextFormatterWith(leftFieldNames ...string) *TextFormatter {
	return &TextFormatter{
		leftFieldKeys: leftFieldNames,
	}
}

// Format returns a formatted string as message "[left fields] message [rest of the fields]".
func (f *TextFormatter) Format(msg string, fields Fields) string {
	leftBuilder := strings.Builder{}
	rightBuilder := strings.Builder{}

	for _, key := range f.leftFieldKeys {
		v, ok := fields[key]
		if !ok {
			continue
		}
		leftBuilder.WriteRune('[')
		leftBuilder.WriteString(fmt.Sprint(v))
		leftBuilder.WriteString("] ")
		delete(fields, key)
	}

	if len(fields) > 0 {
		rightBuilder.WriteRune(' ')
		jsonBytes, _ := json.Marshal(fields)
		rightBuilder.Write(jsonBytes)
	}

	return leftBuilder.String() + msg + rightBuilder.String()
}
