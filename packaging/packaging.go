package packaging

import (
	"log"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

const (
	Beer     = "Beer"
	Packaged = "Packaged"
	Beer1    = "ale"
	Beer2    = "stark"
	Beer3    = "weizen"
)

func NewPackagingBot(port string) *gobot.Robot {
	firmataAdaptor := firmata.NewAdaptor(port)
	button1 := gpio.NewButtonDriver(firmataAdaptor, "7")
	button2 := gpio.NewButtonDriver(firmataAdaptor, "5")
	button3 := gpio.NewButtonDriver(firmataAdaptor, "3")
	led1 := gpio.NewLedDriver(firmataAdaptor, "6")
	led2 := gpio.NewLedDriver(firmataAdaptor, "4")
	led3 := gpio.NewLedDriver(firmataAdaptor, "2")

	b := gobot.NewRobot("packaging",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{button1, button2, button3, led1, led2, led3},
	)

	b.Work = func() {
		// ale, stark, weizen
		selected := ""

		b.AddEvent(Beer)
		b.AddEvent(Packaged)
		b.On(Beer, func(s interface{}) {
			led1.Off()
			led2.Off()
			led3.Off()

			selected = s.(string)
			log.Println(selected)
			switch selected {
			case Beer1:
				led1.On()
			case Beer2:
				led2.On()
			case Beer3:
				led3.On()
			}
		})

		button1.On(gpio.ButtonRelease, func(s interface{}) {
			if selected == Beer1 {
				b.Publish(Packaged, "Picked")
				led1.Off()
				led2.Off()
				led3.Off()
			}
		})

		button2.On(gpio.ButtonRelease, func(s interface{}) {
			if selected == Beer2 {
				b.Publish(Packaged, "Picked")
				led1.Off()
				led2.Off()
				led3.Off()
			}
		})

		button3.On(gpio.ButtonRelease, func(s interface{}) {
			if selected == Beer3 {
				b.Publish(Packaged, "Picked")
				led1.Off()
				led2.Off()
				led3.Off()
			}
		})
	}

	return b
}
