package adapter

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite3 driver
	"github.com/mdp/qrterminal/v3"
	logger "github.com/rafaelcoelhox/whatsapi/helper"
	"github.com/rafaelcoelhox/whatsapi/internal/config"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	"go.uber.org/zap"
)

type WhatsClient struct {
	client *whatsmeow.Client
}

func NewWhatsClient(cfg *config.Config) (*WhatsClient, error) {
	log := logger.GetLogger()
	log.Logger.Info("Creating new client")
	container, err := sqlstore.New(cfg.Cache.Storage, cfg.Cache.File, nil)
	if err != nil {
		log.Logger.Fatal("Failed to create container", zap.Error(err))
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		log.Logger.Fatal("Failed to Get Devicer", zap.Error(err))
	}

	client := whatsmeow.NewClient(deviceStore, nil)
	if client == nil {
		log.Logger.Fatal("Failed to load store", zap.Error(err))
	}
	return &WhatsClient{client: client}, nil
}

func (c *WhatsClient) Start() error {
	log := logger.GetLogger()
	log.Logger.Info("Checling if store is nil")
	if c.client.Store.ID == nil {
		channel, err := c.QRChannel()
		if err != nil {
			log.Logger.Fatal("Failed to get QR channel", zap.Error(err))
		}
		if err := c.Connect(); err != nil {
			log.Logger.Fatal("Failed to connect", zap.Error(err))
		}
		log.Logger.Info("Waiting for QR code")
		c.QRCode(channel)
	} else {
		log.Logger.Info("Logged User")
		err := c.client.Connect()
		if err != nil {
			log.Logger.Fatal("Failed to connect", zap.Error(err))
		}
	}
	d := make(chan os.Signal, 1)
	signal.Notify(d, os.Interrupt, syscall.SIGTERM)
	<-d
	c.client.Disconnect()

	return nil
}

func (c *WhatsClient) Disconnect() {
	c.client.Disconnect()
}

func (c *WhatsClient) AddEventHandler(handler func(interface{})) {
	c.client.AddEventHandler(handler)
}

func (c *WhatsClient) QRChannel() (<-chan whatsmeow.QRChannelItem, error) {
	return c.client.GetQRChannel(context.Background())
}

func (c *WhatsClient) QRCode(qrChan <-chan whatsmeow.QRChannelItem) string {
	for evt := range qrChan {
		if evt.Event == "code" {
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			return evt.Code
		}
	}
	return ""
}

func (c *WhatsClient) Connect() error {
	log := logger.GetLogger()
	if err := c.client.Connect(); err != nil {
		log.Logger.Error("Failed to connect", zap.Error(err))
	}
	return nil
}

func EventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received a message!", v.Message.GetConversation())
	}
}
