package deliverybot

import (
	natsio "github.com/nats-io/nats"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/nats"
)

const (
	BeerBot = "beerbot"
)

func NewDeliveryBot(url string, clientID int, user string, password string) *gobot.Robot {
	natsAdaptor := nats.NewAdaptorWithAuth(url, clientID, user, password, natsio.Secure())

	b := gobot.NewRobot("DeliveryBot",
		[]gobot.Connection{natsAdaptor},
	)

	b.Work = func() {
		b.AddEvent(BeerBot)

		natsAdaptor.On("beerbot", func(m nats.Message) {
			b.Publish(BeerBot, string(m.Data[:]))
		})
	}

	return b
}
