package boombot

import (
	"log"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func NewBoomBot(name string, device string) *gobot.Robot {
	firmataAdaptor := firmata.NewTCPAdaptor(device)
	servo := gpio.NewServoDriver(firmataAdaptor, "13")
	servo.SetName("firework")

	work := func() {
		log.Println("Started")
		servo.Move(150)
	}

	return gobot.NewRobot(name,
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{servo},
		work,
	)
}
