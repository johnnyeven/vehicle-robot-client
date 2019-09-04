package main

import (
	"bytes"
	"fmt"
	"github.com/johnnyeven/libtools/courier/client"
	"github.com/johnnyeven/vehicle-robot-client/client_vehicle_robot"
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
			Host: "localhost",
			Port: 9900,
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

		Rect(img, int(x1), int(y1), int(x2), int(y2), 4, color.White)
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

// HLine draws a horizontal line
func HLine(img *image.RGBA, x1, y, x2 int, col color.Color) {
	for ; x1 <= x2; x1++ {
		img.Set(x1, y, col)
	}
}

// VLine draws a veritcal line
func VLine(img *image.RGBA, x, y1, y2 int, col color.Color) {
	for ; y1 <= y2; y1++ {
		img.Set(x, y1, col)
	}
}

// Rect draws a rectangle utilizing HLine() and VLine()
func Rect(img *image.RGBA, x1, y1, x2, y2, width int, col color.Color) {
	for i := 0; i < width; i++ {
		HLine(img, x1, y1+i, x2, col)
		HLine(img, x1, y2+i, x2, col)
		VLine(img, x1+i, y1, y2, col)
		VLine(img, x2+i, y1, y2, col)
	}
}
