package whatisapiInterface

import "go.mau.fi/whatsmeow"

type WhatsInterface interface {
	Start() error
	Disconnect()
	AddEventHandler(func(interface{}))
	QRChannel() (<-chan whatsmeow.QRChannelItem, error)
	QRCode(<-chan whatsmeow.QRChannelItem) string
	Connect() error
}
