package main

import (
	"bytes"
	"fmt"
	"github.com/johnnyeven/libtools/courier/client"
	"github.com/johnnyeven/vehicle-robot-client/client_vehicle_robot"
	"github.com/johnnyeven/vehicle-robot-client/modules"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
)

func main() {
	webcam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", 0)
		return
	}
	defer webcam.Close()

	webcam.Set(gocv.VideoCaptureFrameWidth, 640)
	webcam.Set(gocv.VideoCaptureFrameHeight, 480)

	cameraImage := gocv.NewMat()
	defer cameraImage.Close()

	cli := &client_vehicle_robot.ClientVehicleRobot{
		Client: client.Client{
			Host: "www.profzone.net",
			Port: 50999,
			Mode: "grpc",
		},
	}
	cli.MarshalDefaults(cli)

	if ok := webcam.Read(&cameraImage); !ok {
		fmt.Printf("Device closed: %v\n", 0)
		return
	}

	sourceImg, err := cameraImage.ToImage()
	if err != nil {
		fmt.Println(err.Error())
	}

	b := sourceImg.Bounds()
	img := image.NewRGBA(b)
	draw.Draw(img, b, sourceImg, b.Min, draw.Src)

	imageFile, err := os.Create("./test.jpg")
	if err != nil {
		return
	}

	err = jpeg.Encode(imageFile, sourceImg, &jpeg.Options{Quality: 75})
	if err != nil {
		fmt.Println(err.Error())
	}

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
		//addLabel(sourceImg, int(x1), int(y1), int(classes[curObj]), getLabel(curObj, probabilities, classes))
	}

	imageFile, err = os.Create("./test1.jpg")
	if err != nil {
		return
	}

	err = jpeg.Encode(imageFile, img, &jpeg.Options{Quality: 75})
	if err != nil {
		fmt.Println(err.Error())
	}
}
