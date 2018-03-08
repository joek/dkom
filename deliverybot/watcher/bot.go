package main

import (
	"log"

	natsio "github.com/nats-io/nats"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
	"gobot.io/x/gobot/platforms/nats"
)

// On Push -> Nats loaded
// on Release -> Nats delivered

func main() {
	natsURL := "172.20.10.8:32240"
	natsClientID := 1235
	natsUser := "natsuser"
	natsPassword := "xqwwVp9C&Hn6jXcux4r)vq"

	// _, currentfile, _, _ := runtime.Caller(0)
	// cascade := path.Join(path.Dir(currentfile), "haarcascade_frontalface_alt.xml")

	// camera := opencv.NewCameraDriver(0)
	// fbox := facebox.New("http://facebox.local/")

	natsAdaptor := nats.NewAdaptorWithAuth(natsURL, natsClientID, natsUser, natsPassword, natsio.Secure())

	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
	button := gpio.NewButtonDriver(firmataAdaptor, "2")

	work := func() {
		// mat := gocv.NewMat()
		// img.Store(mat)

		// camera.On(opencv.Frame, func(data interface{}) {
		// 	i := data.(gocv.Mat)
		// 	img.Store(i)
		// })

		button.On(gpio.ButtonPush, func(interface{}) {
			log.Println("loaded")
			data := []byte("loaded")
			natsAdaptor.Publish("beerbot", data)
		})
		button.On(gpio.ButtonRelease, func(interface{}) {
			log.Println("delivered")
			data := []byte("delivered")
			natsAdaptor.Publish("beerbot", data)
		})

		// gobot.Every(1*time.Second, func() {
		// 	i := img.Load().(gocv.Mat)
		// 	if i.Empty() {
		// 		return
		// 	}
		// 	faces := opencv.DetectObjects(cascade, i)
		// 	opencv.DrawRectangles(i, faces, 0, 255, 0, 5)

		// 	log.Println("Check Image")
		// 	buf, err := gocv.IMEncode(".jpg", i)
		// 	if err != nil {
		// 		log.Printf("unable to encode matrix: %v", err)
		// 		return
		// 	}

		// 	found, err := fbox.Check(bytes.NewReader(buf))
		// 	if err != nil {
		// 		log.Printf("unable to recognize face: %v", err)
		// 	}

		// 	var caption = "I don't know you"
		// 	if len(found) > 0 {
		// 		caption = fmt.Sprintf("I know you %s", found[0].Name)
		// 		natsAdaptor.Publish("beerbot", []byte(fmt.Sprintf("user:%s", found[0].Name)))
		// 	}

		// 	log.Println(caption)
		// })
	}

	robot := gobot.NewRobot("beerbot",
		[]gobot.Connection{natsAdaptor, firmataAdaptor},
		[]gobot.Device{button},
		work,
	)

	robot.Start()

}
