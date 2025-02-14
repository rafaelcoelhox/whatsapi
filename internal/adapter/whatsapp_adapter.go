package adapter

import (
	"context"
	"fmt"
	"os"

	"github.com/mdp/qrterminal/v3"
	"github.com/rafaelcoelhox/whatsapi/internal/config"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type WhatsClient struct {
	client *whatsmeow.Client
}

func NewWhatsClient(cfg *config.Config) (*WhatsClient, error) {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New(cfg.Cache.Storage, cfg.Cache.File, dbLog)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage container: %w", err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	return &WhatsClient{client: client}, nil
}

func (c *WhatsClient) Start() error {
	if c.client.Store.ID == nil {
		channel, err := c.QRChannel()
		if err != nil {
			return fmt.Errorf("failed to get QR channel: %w", err)
		}

		if err := c.Connect(); err != nil {
			return err
		}

		c.QRCode(channel)
		return nil
	}

	fmt.Println("Already authenticated")
	return c.Connect()
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
			config := qrterminal.Config{
				Level:     qrterminal.M,
				Writer:    os.Stdout,
				BlackChar: qrterminal.BLACK,
				WhiteChar: qrterminal.WHITE,
				QuietZone: 2,
			}
			qrterminal.GenerateWithConfig(evt.Code, config)
			return evt.Code
		}
	}
	return ""
}

func (c *WhatsClient) Connect() error {
	if err := c.client.Connect(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	return nil
}

func EventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received a message!", v.Message.GetConversation())
	}
}
