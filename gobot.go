package main

import (
	"bytes"
	"fmt"
	"github.com/johnnyeven/libtools/courier/client"
	"github.com/johnnyeven/vehicle-robot-client/client_vehicle_robot"
	"github.com/johnnyeven/vehicle-robot-client/modules"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/opencv"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
)

func main() {
	window := gocv.NewWindow("Detective")
	camera := opencv.NewCameraDriver(0)

	cli := &client_vehicle_robot.ClientVehicleRobot{
		Client: client.Client{
			Host: "www.profzone.net",
			Port: 50999,
			Mode: "grpc",
		},
	}
	cli.MarshalDefaults(cli)

	work := func() {
		err := camera.On(opencv.Frame, func(data interface{}) {
			cameraImage := data.(gocv.Mat)
			defer cameraImage.Close()

			sourceImg, err := cameraImage.ToImage()
			if err != nil {
				fmt.Println(err.Error())
			}

			b := sourceImg.Bounds()
			img := image.NewRGBA(b)
			draw.Draw(img, b, sourceImg, b.Min, draw.Src)

			buf := bytes.NewBuffer([]byte{})
			err = jpeg.Encode(buf, sourceImg, &jpeg.Options{Quality: 100})
			if err != nil {
				fmt.Println(err.Error())
			}

			request := client_vehicle_robot.ObjectDetectionRequest{
				Body: client_vehicle_robot.ObjectDetectionBody{
					Image: buf.Bytes(),
				},
			}
			resp, err := cli.ObjectDetection(request)
			if err != nil {
				fmt.Println("request err: ", err)
			}

			for _, detectived := range resp.Body {
				x1 := float32(img.Bounds().Max.X) * detectived.Box[1]
				x2 := float32(img.Bounds().Max.X) * detectived.Box[3]
				y1 := float32(img.Bounds().Max.Y) * detectived.Box[0]
				y2 := float32(img.Bounds().Max.Y) * detectived.Box[2]

				modules.Rect(img, int(x1), int(y1), int(x2), int(y2), 4, color.White)
			}

			targetImage, err := modules.ConvertImageToMat(img)
			if err != nil {
				fmt.Println(err.Error())
			}
			defer targetImage.Close()

			window.IMShow(targetImage)
		})
		if err != nil {

		}
	}

	robot := gobot.NewRobot("cameraBot",
		[]gobot.Device{camera},
		work,
	)

	robot.Start()
}
