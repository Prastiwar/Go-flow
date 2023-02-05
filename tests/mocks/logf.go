package mocks

import (
	"io"

	"github.com/Prastiwar/Go-flow/logf"
)

var (
	_ logf.Formatter = FormatterMock{}
	_ logf.Logger    = LoggerMock{}
)

type FormatterMock struct {
	OnFormat func(msg string, fields logf.Fields) string
}

func (m FormatterMock) Format(msg string, fields logf.Fields) string {
	return m.OnFormat(msg, fields)
}

type LoggerMock struct {
	OnDebug     func(v interface{})
	OnDebugf    func(format string, args ...any)
	OnError     func(v interface{})
	OnErrorf    func(format string, args ...any)
	OnFormatter func() logf.Formatter
	OnInfo      func(v interface{})
	OnInfof     func(format string, args ...any)
	OnOutput    func() io.Writer
	OnScope     func() logf.Fields
}

func (m LoggerMock) Debug(v interface{}) {
	m.OnDebug(v)
}

func (m LoggerMock) Debugf(format string, args ...any) {
	m.OnDebugf(format, args...)
}

func (m LoggerMock) Error(v interface{}) {
	m.OnError(v)
}

func (m LoggerMock) Errorf(format string, args ...any) {
	m.OnErrorf(format, args...)
}

func (m LoggerMock) Formatter() logf.Formatter {
	return m.OnFormatter()
}

func (m LoggerMock) Info(v interface{}) {
	m.OnInfo(v)
}

func (m LoggerMock) Infof(format string, args ...any) {
	m.OnInfof(format, args...)
}

func (m LoggerMock) Output() io.Writer {
	return m.OnOutput()
}

func (m LoggerMock) Scope() logf.Fields {
	return m.OnScope()
}
