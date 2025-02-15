package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
	"go.uber.org/zap"
	
	logger "github.com/rafaelhox/whatsapi/helper"
	"github.com/rafaelcoelhox/whatsapi/internal/adapter"
	"github.com/rafaelcoelhox/whatsapi/internal/config"
	whatisapiInterface "github.com/rafaelcoelhox/whatsapi/internal/interface"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}
	
	log := logger.InitLogger(cfg)
	defer log.Logger.Sync()
	


	/*
	var client .WhatsInterface
	client, err = adapter.NewWhatsClient(cfg)
	if err != nil {
		log.Logger.Error("Failed to create client", zap.Error(err))
	}
	client.AddEventHandler(adapter.EventHandler)
	if err := client.Start(); err != nil {
		log.Logger.Error("Failed to start client", zap.Error(err))
	}
	select {} */


}
