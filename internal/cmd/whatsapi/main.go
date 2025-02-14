package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver

	"github.com/rafaelcoelhox/whatsapi/internal/adapter"
	"github.com/rafaelcoelhox/whatsapi/internal/config"
	whatisapiInterface "github.com/rafaelcoelhox/whatsapi/internal/interface"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	var client whatisapiInterface.WhatsInterface
	client, err = adapter.NewWhatsClient(cfg)
	if err != nil {
		log.Fatalf("Error creating WhatsClient: %v", err)
	}

	client.AddEventHandler(adapter.EventHandler)

	if err := client.Start(); err != nil {
		log.Fatal(err)
	}

	select {}
}
