package main

import (
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
	"go.uber.org/zap"

	logger "github.com/rafaelcoelhox/whatsapi/helper"
	"github.com/rafaelcoelhox/whatsapi/internal/adapter"
	"github.com/rafaelcoelhox/whatsapi/internal/config"
	whatisapiInterface "github.com/rafaelcoelhox/whatsapi/internal/interface"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config: ", err)
	}

	logger := logger.Init(cfg)
	defer logger.Logger.Sync()

	var client whatisapiInterface.WhatsInterface
	client, err = adapter.NewWhatsClient(cfg)
	if err != nil {
		logger.Logger.Error("Error creating new client", zap.Error(err))
	}

	client.AddEventHandler(adapter.EventHandler)

	if err := client.Start(); err != nil {
		log.Fatal(err)
	}

	select {}
}
