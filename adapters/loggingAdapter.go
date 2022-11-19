package adapters

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/AlecSmith96/fourth-floor-exercise/entities"
)

func NewLogger(config entities.ConfigLogging) (*zap.Logger, error) {
	// get loglevel from config
	level := zap.NewAtomicLevel()
	err := level.UnmarshalText([]byte(config.LogLevel))
	if err != nil {
		return nil, fmt.Errorf("initializing logging: %v", err)
	}

	// encoderConfig sets the default fields for each log
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// global config for logger instance
	logger, err := zap.Config{
		Level:       level,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:          config.Encoding,
		EncoderConfig:     encoderConfig,
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableStacktrace: true,
	}.Build()

	if err != nil {
		return nil, err
	}

	logger.Named("logging").Info("Initialized logging adapter")
	return logger, nil
}