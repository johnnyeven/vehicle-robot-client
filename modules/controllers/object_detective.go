package controllers

import (
	"bytes"
	"fmt"
	"github.com/gofrs/flock"
	"github.com/johnnyeven/vehicle-robot-client/client"
	"github.com/johnnyeven/vehicle-robot-client/modules"
	"gobot.io/x/gobot/platforms/opencv"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
)

func ObjectDetectiveController(window *opencv.WindowDriver, camera *opencv.CameraDriver, cli *client.RobotClient) {
	locker := flock.New("/dev/lock/camera.lock")
	err := camera.On(opencv.Frame, func(data interface{}) {
		locked, err := locker.TryLock()
		if locked {
			defer locker.Unlock()
		} else {
			return
		}
		cameraImage := data.(gocv.Mat)

		sourceImg, err := cameraImage.ToImage()
		if err != nil {
			fmt.Println("cameraImage.ToImag err: ", err.Error())
			return
		}

		b := sourceImg.Bounds()
		fmt.Println(b.Dx(), b.Dy())
		img := image.NewRGBA(b)
		draw.Draw(img, b, sourceImg, b.Min, draw.Src)

		buf := bytes.NewBuffer([]byte{})
		err = jpeg.Encode(buf, sourceImg, &jpeg.Options{Quality: 75})
		if err != nil {
			fmt.Println("jpeg.Encode err: ", err.Error())
			return
		}

		request := client.ObjectDetectionBody{
			Image: buf.Bytes(),
		}
		resp, err := cli.DetectionObject(request)
		if err != nil {
			fmt.Println("request err: ", err)
			return
		}

		for _, detectived := range resp {
			x1 := float32(img.Bounds().Max.X) * detectived.Box[1]
			x2 := float32(img.Bounds().Max.X) * detectived.Box[3]
			y1 := float32(img.Bounds().Max.Y) * detectived.Box[0]
			y2 := float32(img.Bounds().Max.Y) * detectived.Box[2]

			modules.Rect(img, int(x1), int(y1), int(x2), int(y2), 4, color.White)
		}

		targetImage, err := modules.ConvertImageToMat(img)
		if err != nil {
			fmt.Println("modules.ConvertImageToMat err: ", err.Error())
			return
		}

		window.ShowImage(targetImage)
		window.WaitKey(1)
	})
	if err != nil {
		fmt.Println("camera.On err: ", err.Error())
	}
}
