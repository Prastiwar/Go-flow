package logf

import (
	"io"
	"strings"
)

type Fields map[string]interface{}

func MergeFields(source Fields, fields Fields) Fields {
	if fields == nil && source == nil {
		return make(Fields)
	}

	if fields == nil {
		return source
	}

	if source == nil {
		return fields
	}

	for k, v := range fields {
		source[k] = v
	}

	return source
}

type Formatter interface {
	Format(msg string, fields Fields) string
}

type formatterWriter struct {
	formatter Formatter
	writer    io.Writer
	fields    Fields
}

func (f formatterWriter) Write(p []byte) (n int, err error) {
	s := f.formatter.Format(strings.TrimSuffix(string(p), "\n"), f.fields)
	buf := []byte(s + "\n")
	return f.writer.Write(buf)
}

func DefaultFormatter() Formatter {
	return &TextFormatter{}
}
