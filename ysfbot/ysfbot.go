package ysfbot

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	natsio "github.com/nats-io/nats"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/nats"
)

const (
	NatsOrderCreated  = "OrderCreated"
	OrderCreated      = "OrderCreated"
	UpdateOrderStatus = "UpdateOrder"
)

type OrderStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type Order struct {
	ID        string  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	Address   Address `json:"address"`
	Items     []Item  `json:"items"`
	Status    string  `json:"status"`
}

type Item struct {
	Quantitiy int    `json:"quantity"`
	Name      string `json:"name"`
	Price     string `json:"price"`
	Currency  string `json:"currency"`
	URL       string `json:"url"`
	Sku       string `json:"sku"`
}
type Address struct {
	RecipientName string `json:"recipient_name"`
	Line1         string `json:"line1"`
	City          string `json:"city"`
	State         string `json:"state"`
	PostalCode    string `json:"postal_code"`
	CountryCode   string `json:"country_code"`
}

type YSFRobot *gobot.Robot

func NewYSFRobot(url string, clientID int, user string, password string, orderUpdateUrl string) YSFRobot {
	natsAdaptor := nats.NewAdaptorWithAuth(url, clientID, user, password, natsio.Secure())

	y := gobot.NewRobot("ysfBot",
		[]gobot.Connection{natsAdaptor},
	)

	y.Work = func() {
		y.AddEvent(OrderCreated)
		y.AddEvent(UpdateOrderStatus)

		natsAdaptor.On(NatsOrderCreated, func(m nats.Message) {
			var o Order
			log.Println(string(m.Data[:]))
			err := json.Unmarshal(m.Data, &o)
			if err != nil {
				log.Println("error:", err)
			}

			y.Publish(OrderCreated, o)
		})

		y.On(UpdateOrderStatus, func(params interface{}) {
			json, err := json.Marshal(params.(OrderStatus))
			if err != nil {
				log.Println(err)
			}
			log.Println("Call update order service")
			_, err = http.Post(orderUpdateUrl, "application/json", bytes.NewBuffer(json))
			if err != nil {
				log.Println(err)
			}
		})
	}

	return y
}
