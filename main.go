package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	logger "github.com/rafaelcoelhox/whatsapi/helper"
	"github.com/rafaelcoelhox/whatsapi/internal/adapter"
	"github.com/rafaelcoelhox/whatsapi/internal/config"
	whatisapiInterface "github.com/rafaelcoelhox/whatsapi/internal/interface"
	"go.uber.org/zap"
)

func Ping(c *gin.Context) {

}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	log := logger.InitLogger(cfg)
	defer log.Logger.Sync()

	var client whatisapiInterface.WhatsInterface
	client, err = adapter.NewWhatsClient(cfg)
	if err != nil {
		log.Logger.Error("Failed to create client", zap.Error(err))
	}
	client.AddEventHandler(adapter.EventHandler)
	if err := client.Start(); err != nil {
		log.Logger.Error("Failed to start client", zap.Error(err))
	}

	d := make(chan os.Signal, 1)
	signal.Notify(d, os.Interrupt, syscall.SIGTERM)
	<-d

	client.Disconnect()
}
