package logger

import (
	"fmt"

	"github.com/klimenkokayot/vk-internship/libs/logger/adapters/logrus"
	"github.com/klimenkokayot/vk-internship/libs/logger/adapters/zap"
	"github.com/klimenkokayot/vk-internship/libs/logger/domain"
)

const (
	AdapterZap    = "zap"
	AdapterLogrus = "logrus"
)

var (
	ErrUnknownAdapter = fmt.Errorf("логгер не поддерживается")
)

type (
	Level  = domain.Level
	Field  = domain.Field
	Logger = domain.Logger
)

// Реэкспорт констант уровня
const (
	LevelDebug = domain.LevelDebug
	LevelInfo  = domain.LevelInfo
	LevelWarn  = domain.LevelWarn
	LevelError = domain.LevelError
	LevelFatal = domain.LevelFatal
)

type Config struct {
	Adapter string
	Level   Level
}

func NewAdapter(config *Config) (Logger, error) {
	switch config.Adapter {
	case AdapterZap:
		return zap.NewAdapter(config.Level)
	case AdapterLogrus:
		return logrus.NewAdapter(config.Level)
	default:
		return nil, ErrUnknownAdapter
	}
}
