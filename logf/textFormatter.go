package logf

type TextFormatter struct {
}

func NewTextFormatter() *TextFormatter {
	return &TextFormatter{}
}

func (f *TextFormatter) Format(msg string, fields Fields) string {
	// TODO: implement
	return msg
}
