package controllers

import (
	"bytes"
	"fmt"
	"github.com/johnnyeven/service-vehicle-robot/constants/types"
	"github.com/johnnyeven/vehicle-robot-client/client"
	"github.com/johnnyeven/vehicle-robot-client/global"
	"github.com/johnnyeven/vehicle-robot-client/modules"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
)

func ObjectDetectiveController(config global.RobotConfiguration, camera *gocv.VideoCapture, cli *client.RobotClient) {
	for {
		cameraImage := gocv.NewMat()
		if !camera.Read(&cameraImage) {
			break
		}

		sourceImg, err := cameraImage.ToImage()
		if err != nil {
			fmt.Println("cameraImage.ToImag err: ", err.Error())
			return
		}

		b := sourceImg.Bounds()
		img := image.NewRGBA(b)
		draw.Draw(img, b, sourceImg, b.Min, draw.Src)

		buf := bytes.NewBuffer([]byte{})
		err = jpeg.Encode(buf, sourceImg, &jpeg.Options{Quality: 75})
		if err != nil {
			fmt.Println("jpeg.Encode err: ", err.Error())
			return
		}

		if config.CameraMode == types.CAMERA_MODE__OBJECT_DETECTIVE {
			resp, err := cli.DetectionObject(buf.Bytes())
			if err != nil {
				fmt.Println("cli.DetectionObject request err: ", err)
				return
			}

			for _, detectived := range resp {
				x1 := float32(img.Bounds().Max.X) * detectived.Box[1]
				x2 := float32(img.Bounds().Max.X) * detectived.Box[3]
				y1 := float32(img.Bounds().Max.Y) * detectived.Box[0]
				y2 := float32(img.Bounds().Max.Y) * detectived.Box[2]

				modules.Rect(img, int(x1), int(y1), int(x2), int(y2), 4, color.White)
			}

			//targetImage, err := modules.ConvertImageToMat(img)
			//if err != nil {
			//	fmt.Println("modules.ConvertImageToMat err: ", err.Error())
			//	return
			//}
		} else if config.CameraMode == types.CAMERA_MODE__NORMAL {
			err = cli.CameraTransfer(buf.Bytes())
			if err != nil {
				fmt.Println("cli.CameraTransfer push err: ", err)
				return
			}
		}
	}
}
