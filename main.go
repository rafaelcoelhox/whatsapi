package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
	"github.com/mdp/qrterminal/v3"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

const (
	CachePATH  = "cache"
	ServerPORT = ":8080"
)

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received a message!", v.Message.GetConversation())
	}
}

func main() {
	client := CreateWhatsClient()
	client.AddEventHandler(eventHandler)

	if client.client.Store.ID == nil {
		channel := client.QRChannel()
		client.Connect()
		client.QRCode(channel)
	}

}

type WhatsClient struct {
	client *whatsmeow.Client
}

func CreateWhatsClient() *WhatsClient {

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	return &WhatsClient{client: client}

}

func (c *WhatsClient) Disconnect() {
	c.client.Disconnect()
}

func (c *WhatsClient) AddEventHandler(handler func(interface{})) {
	c.client.AddEventHandler(handler)
}

func (c *WhatsClient) QRChannel() <-chan whatsmeow.QRChannelItem {
	qrChan, _ := c.client.GetQRChannel(context.Background())
	return qrChan
}

func (c *WhatsClient) QRCode(qrChan <-chan whatsmeow.QRChannelItem) string {

	for evt := range qrChan {
		if evt.Event == "code" {
			config := qrterminal.Config{
				Level:     qrterminal.L,
				Writer:    os.Stdout,
				BlackChar: qrterminal.BLACK,
				WhiteChar: qrterminal.WHITE,
				QuietZone: 1,
			}
			qrterminal.GenerateWithConfig(evt.Code, config)
			return evt.Code
		}

	}
	return ""
}

func (c *WhatsClient) Connect() {
	err := c.client.Connect()
	if err != nil {
		panic(err)
	}
}
