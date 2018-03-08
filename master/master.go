package main

import (
	"gobot.io/x/gobot/api"
	"gobot.io/x/gobot/drivers/gpio"

	"github.com/joek/dkom/boombot"
	"github.com/joek/dkom/deliverybot"
	"github.com/joek/dkom/packaging"
	"github.com/joek/dkom/ysfbot"
	"gobot.io/x/gobot"
)

func main() {
	var order *ysfbot.Order

	master := gobot.NewMaster()
	a := api.NewAPI(master)
	a.Start()

	ysf := createNewYsf()
	pack := packaging.NewPackagingBot("/dev/tty.usbmodem14131")
	deliveryBot := newDeliveryBot()
	boom1 := boombot.NewBoomBot("boom1", "172.20.10.3:3030")
	boom2 := boombot.NewBoomBot("boom2", "172.20.10.13:3030")

	master.AddRobot(pack)
	master.AddRobot(deliveryBot)
	master.AddRobot(boom1)
	master.AddRobot(boom2)
	master.AddRobot(ysf)

	ysf.On(ysfbot.OrderCreated, func(m interface{}) {
		o := m.(ysfbot.Order)
		order = &o
		ysf.Publish(ysfbot.UpdateOrderStatus, ysfbot.OrderStatus{
			ID:     o.ID,
			Status: "In Preperation",
		})
		pack.Publish(packaging.Beer, order.Items[0].Sku)
	})

	pack.On(packaging.Packaged, func(interface{}) {
		ysf.Publish(ysfbot.UpdateOrderStatus, ysfbot.OrderStatus{
			ID:     order.ID,
			Status: "Packaged",
		})
	})

	deliveryBot.On(deliverybot.BeerBot, func(m interface{}) {
		if order != nil {
			switch m.(string) {
			case "loaded":
				ysf.Publish(ysfbot.UpdateOrderStatus, ysfbot.OrderStatus{
					ID:     order.ID,
					Status: "Shipping",
				})
			case "delivered":
				ysf.Publish(ysfbot.UpdateOrderStatus, ysfbot.OrderStatus{
					ID:     order.ID,
					Status: "Delivered",
				})
				boom1.Device("firework").(*gpio.ServoDriver).Move(0)
				boom2.Device("firework").(*gpio.ServoDriver).Move(0)
			}
		}
	})

	master.Start()
}

func createNewYsf() ysfbot.YSFRobot {
	natsURL := "127.0.0.1:32240"
	natsClientID := 1234
	natsUser := "natsuser"
	natsPassword := "xqwwVp9C&Hn6jXcux4r)vq"
	orderUpdateURL := "http://beershop.local/orderStatusUpdate"
	return ysfbot.NewYSFRobot(natsURL, natsClientID, natsUser, natsPassword, orderUpdateURL)
}

func newDeliveryBot() *gobot.Robot {
	natsURL := "172.20.10.8:32240"
	natsClientID := 1234
	natsUser := "natsuser"
	natsPassword := "xqwwVp9C&Hn6jXcux4r)vq"
	return deliverybot.NewDeliveryBot(natsURL, natsClientID, natsUser, natsPassword)
}
