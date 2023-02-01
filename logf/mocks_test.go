package logf

type FormatterMock struct {
	formatFn func(msg string, fields Fields) string
}

func NewFormatterMock(formatFn func(msg string, fields Fields) string) *FormatterMock {
	return &FormatterMock{formatFn: formatFn}
}

func (m *FormatterMock) Format(msg string, fields Fields) string {
	return m.formatFn(msg, fields)
}
