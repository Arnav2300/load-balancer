package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerConfig struct {
	Env      string
	LogLevel string
	LogFile  string
}

func NewLogger(cfg LoggerConfig) (*zap.Logger, error) {
	if cfg.Env == "" {
		cfg.Env = "development"
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}
	level := zapcore.InfoLevel
	if err := level.UnmarshalText([]byte(strings.ToLower(cfg.LogLevel))); err != nil {
		return nil, err
	}

	var encoderConfig zapcore.EncoderConfig
	if cfg.Env == "development" {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
	}
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	if cfg.Env == "development" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	var ws zapcore.WriteSyncer
	if cfg.LogFile != "" {
		file, err := os.OpenFile(cfg.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		ws = zapcore.AddSync(file)
	} else {
		ws = zapcore.AddSync(os.Stdout)
	}

	core := zapcore.NewCore(encoder, ws, level)

	opts := []zap.Option{}
	if cfg.Env == "development" {
		opts = append(opts, zap.AddCaller())
	}

	logger := zap.New(core, opts...)
	return logger, nil

}
