package log

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Options struct {
	Level       zapcore.Level
	Development bool
}

type Option func(*Options)

func WithLevel(level zapcore.Level) Option {
	return func(o *Options) {
		o.Level = level
	}
}

func WithDevelopment(development bool) Option {
	return func(o *Options) {
		o.Development = development
	}
}

var Logger logr.Logger

func Init(opts ...Option) logr.Logger {
	options := Options{
		Level:       zap.InfoLevel,
		Development: false,
	}

	for _, opt := range opts {
		opt(&options)
	}

	var zapConfig zap.Config
	if options.Development {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}

	zapConfig.Level.SetLevel(options.Level)

	zapLog, err := zapConfig.Build()
	if err != nil {
		fmt.Errorf("%v", err)
	}

	Logger = zapr.NewLogger(zapLog)
	return Logger
}
