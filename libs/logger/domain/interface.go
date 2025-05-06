package domain

import (
	"fmt"
)

type Field struct {
	Key   string
	Value interface{}
}

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	WithFields(fields ...Field) Logger
	WithLayer(name string) Logger
	OK(msg string, fields ...Field)
}

func String(msg, value string) Field {
	return Field{msg, value}
}

func Int(msg string, value int) Field {
	return Field{msg, value}
}

func Error(value error) Field {
	return Field{Value: value}
}
