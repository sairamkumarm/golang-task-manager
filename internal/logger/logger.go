package logger

import (
	"fmt"
	"golang-task-manager/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func BuildLogger(config *config.Config) (*zap.Logger, error) {
	level, err := parseLevel(config.LOG_LEVEL)
	if err != nil {
		return nil, err
	}

	atom:=zap.NewAtomicLevel()
	atom.SetLevel(level)

	zapConfig:= zap.Config{
		Level: atom,
		Development: false,
		Encoding: "json",
		EncoderConfig: zap.NewProductionEncoderConfig(),
		OutputPaths: []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	// fmt.Println(zapConfig)
	return zapConfig.Build()
}

func parseLevel(level string) (zapcore.Level, error) {
	switch level {
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("invalid log level: %s", level)
	}
}