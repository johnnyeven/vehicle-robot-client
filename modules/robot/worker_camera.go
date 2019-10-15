package robot

import (
	"bytes"
	"fmt"
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/libtools/courier/enumeration"
	"github.com/johnnyeven/service-vehicle-robot/constants/types"
	"github.com/johnnyeven/vehicle-robot-client/client"
	"github.com/johnnyeven/vehicle-robot-client/global"
	"github.com/johnnyeven/vehicle-robot-client/modules"
	"github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
	"image"
	"image/draw"
	"image/jpeg"
)

const (
	cameraWorkerID = "camera-worker"
)

type CameraWorker struct {
	camera *gocv.VideoCapture
	cli    *client.RobotClient
	bus    *bus.MessageBus

	cameraMode             types.CameraMode
	activateCameraTransfer enumeration.Bool
}

func NewCameraWorker(robot *Robot, bus *bus.MessageBus, cli *client.RobotClient, config *global.RobotConfiguration) *CameraWorker {
	camera, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		logrus.Panicf("[CameraWorker] gocv.VideoCaptureDevice err: %v", err)
	}
	camera.Set(gocv.VideoCaptureFrameWidth, float64(config.CameraCaptureWidth))
	camera.Set(gocv.VideoCaptureFrameHeight, float64(config.CameraCaptureHeight))

	return &CameraWorker{
		camera:                 camera,
		cli:                    cli,
		bus:                    bus,
		cameraMode:             config.CameraMode,
		activateCameraTransfer: config.ActivateCameraTransfer,
	}
}

func (c *CameraWorker) WorkerID() string {
	return cameraWorkerID
}

func (c *CameraWorker) Start() {
	for {
		cameraImage := gocv.NewMat()
		if !c.camera.Read(&cameraImage) {
			break
		}

		sourceImg, err := cameraImage.ToImage()
		if err != nil {
			fmt.Println("[CameraWorker] cameraImage.ToImag err: ", err.Error())
			return
		}

		b := sourceImg.Bounds()
		img := image.NewRGBA(b)
		draw.Draw(img, b, sourceImg, b.Min, draw.Src)

		if c.cameraMode == types.CAMERA_MODE__OBJECT_DETECTIVE {
			buf := bytes.NewBuffer([]byte{})
			err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 75})
			if err != nil {
				fmt.Println("[CameraWorker] jpeg.Encode err: ", err.Error())
				return
			}
			resp, err := c.cli.DetectionObject(buf.Bytes())
			if err != nil {
				fmt.Println("[CameraWorker] cli.DetectionObject request err: ", err)
				return
			}

			for _, detectived := range resp {
				x1 := float32(img.Bounds().Max.X) * detectived.Box[1]
				x2 := float32(img.Bounds().Max.X) * detectived.Box[3]
				y1 := float32(img.Bounds().Max.Y) * detectived.Box[0]
				y2 := float32(img.Bounds().Max.Y) * detectived.Box[2]

				modules.Rect(img, int(x1), int(y1), int(x2), int(y2), 4, modules.GetLabelColor(int(detectived.Class)))
				modules.DrawLabel(img, int(x1), int(y1), int(detectived.Class), detectived.Label)
			}
		}

		if c.activateCameraTransfer.True() {
			buf := bytes.NewBuffer([]byte{})
			err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 75})
			if err != nil {
				fmt.Println("[CameraWorker] jpeg.Encode err: ", err.Error())
				return
			}

			err = c.cli.CameraTransfer(buf.Bytes())
			if err != nil {
				fmt.Println("[CameraWorker] cli.CameraTransfer push err: ", err)
				return
			}
		}
	}
}

func (c *CameraWorker) Restart() error {
	panic("implement me")
}

func (c *CameraWorker) Stop() error {
	return c.camera.Close()
}
