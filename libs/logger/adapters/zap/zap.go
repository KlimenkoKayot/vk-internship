package zap

import (
	"fmt"
	"sync"

	"github.com/klimenkokayot/vk-internship/libs/logger/domain"
	"github.com/klimenkokayot/vk-internship/libs/logger/pkg/colorise"
	"github.com/klimenkokayot/vk-internship/libs/logger/pkg/formatter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapAdapter struct {
	logger    *zap.Logger
	fields    []zap.Field
	formatter *formatter.Formatter
	mu        sync.RWMutex
}

func (z *ZapAdapter) WithFields(fields ...domain.Field) domain.Logger {
	zapFields := toRouterFields(fields)

	z.mu.RLock()
	defer z.mu.RUnlock()

	return &ZapAdapter{
		logger:    z.logger,
		fields:    append(zapFields, z.fields...),
		formatter: z.formatter,
	}
}

func (z *ZapAdapter) WithLayer(name string) domain.Logger {
	z.mu.RLock()
	defer z.mu.RUnlock()

	return &ZapAdapter{
		logger:    z.logger,
		fields:    z.fields,
		formatter: formatter.NewFormatter(name),
	}
}

func (z *ZapAdapter) log(level zapcore.Level, msg string, fields []domain.Field, color colorise.Color) {
	formattedMsg := z.formatter.FormatMessage(msg)
	formattedMsg = colorise.ColorString(formattedMsg, color)

	zapFields := toRouterFields(fields)

	z.mu.RLock()
	allFields := append(zapFields, z.fields...)
	z.mu.RUnlock()

	switch level {
	case zap.DebugLevel:
		z.logger.Debug(formattedMsg, allFields...)
	case zap.InfoLevel:
		z.logger.Info(formattedMsg, allFields...)
	case zap.WarnLevel:
		z.logger.Warn(formattedMsg, allFields...)
	case zap.ErrorLevel:
		z.logger.Error(formattedMsg, allFields...)
	case zap.FatalLevel:
		z.logger.Fatal(formattedMsg, allFields...)
	}
}

// Методы-обёртки для удобства
func (z *ZapAdapter) Debug(msg string, fields ...domain.Field) {
	z.log(zap.DebugLevel, msg, fields, colorise.ColorReset)
}

func (z *ZapAdapter) Info(msg string, fields ...domain.Field) {
	z.log(zap.InfoLevel, msg, fields, colorise.ColorReset)
}

func (z *ZapAdapter) Warn(msg string, fields ...domain.Field) {
	z.log(zap.WarnLevel, msg, fields, colorise.ColorYellow)
}

func (z *ZapAdapter) Error(msg string, fields ...domain.Field) {
	z.log(zap.ErrorLevel, msg, fields, colorise.ColorRed)
}

func (z *ZapAdapter) Fatal(msg string, fields ...domain.Field) {
	z.log(zap.FatalLevel, msg, fields, colorise.ColorRed)
}

func (z *ZapAdapter) OK(msg string, fields ...domain.Field) {
	z.log(zap.InfoLevel, msg, fields, colorise.ColorGreen)
}

func NewAdapter(level domain.Level) (domain.Logger, error) {
	zapCfg := zap.NewProductionConfig()
	zapCfg.Encoding = "console"
	zapCfg.Level = toRouterLevel(level)

	zapLogger, err := zapCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrZapBuild, err.Error())
	}

	return &ZapAdapter{
		logger:    zapLogger,
		fields:    make([]zap.Field, 0),
		formatter: formatter.NewFormatter(""),
	}, nil
}

func toRouterLevel(level domain.Level) zap.AtomicLevel {
	switch level {
	case domain.LevelDebug:
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case domain.LevelInfo:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case domain.LevelWarn:
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case domain.LevelError:
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	case domain.LevelFatal:
		return zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	}
}

func toRouterFields(fields []domain.Field) []zap.Field {
	converted := []zap.Field{}
	for _, val := range fields {
		field := zap.Field{}
		switch val.Value.(type) {
		case string:
			field = zap.String(val.Key, val.Value.(string))
		case int:
			field = zap.Int(val.Key, val.Value.(int))
		default:
			field = zap.Any(val.Key, val.Value)
		}
		converted = append(converted, field)
	}
	return converted
}
