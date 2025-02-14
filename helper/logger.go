package logger

import (
	"os"

	"github.com/rafaelcoelhox/whatsapi/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	Logger *zap.Logger
}

func Init(cfg *config.Config) Logger {

	var outputPaths []string
	encoding := "json"

	if cfg.AppEnv.Env == "development" {
		outputPaths = []string{"stdout"}
		encoding = "console"
	} else {
		outputPaths = []string{
			"stdout",
			cfg.AppEnv.LogFile,
		}
	}

	LogConfig := zap.Config{
		OutputPaths: outputPaths,
		Encoding:    encoding,
		Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			MessageKey:   "Message",
			LevelKey:     "Level",
			NameKey:      "Name",
			CallerKey:    "Caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
			EncodeLevel:  zapcore.CapitalLevelEncoder,
		},
	}

	if cfg.AppEnv.Env == "production" {
		if err := os.MkdirAll("logs", 07555); err != nil {
			panic(err)
		}
	}

	logger, _ := LogConfig.Build()

	return Logger{
		Logger: logger,
	}

}
