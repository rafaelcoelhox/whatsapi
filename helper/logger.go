package logger

import (
	"go.uber.org/zap"
)

// Logger é uma instância global do logger
var Logger *zap.Logger

// Init inicializa o logger
func Init() {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

// Sync realiza a sincronização do logger
func Sync() {
	if err := Logger.Sync(); err != nil {
		panic(err)
	}
}
