package formatter

import "fmt"

type Formatter struct {
	prefix string
}

func NewFormatter(prefix string) *Formatter {
	return &Formatter{
		prefix,
	}
}

func (f *Formatter) FormatMessage(msg string) string {
	if f.prefix == "" {
		return msg
	}
	return fmt.Sprintf("[%s] %s", f.prefix, msg)
}
