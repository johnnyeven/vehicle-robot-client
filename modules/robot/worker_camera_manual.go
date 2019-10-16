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
	cameraManualWorkerID = "camera-manual-worker"
)

type CameraManualWorker struct {
	camera *gocv.VideoCapture
	cli    *client.RobotClient
	bus    *bus.MessageBus

	robotMode              types.RobotMode
	cameraMode             types.CameraMode
	activateCameraTransfer enumeration.Bool
}

func NewCameraManualWorker(robot *Robot, bus *bus.MessageBus, cli *client.RobotClient, config *global.RobotConfiguration) *CameraManualWorker {
	camera, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		logrus.Panicf("[CameraManualWorker] gocv.VideoCaptureDevice err: %v", err)
	}
	camera.Set(gocv.VideoCaptureFrameWidth, float64(config.CameraCaptureWidth))
	camera.Set(gocv.VideoCaptureFrameHeight, float64(config.CameraCaptureHeight))

	return &CameraManualWorker{
		camera:                 camera,
		cli:                    cli,
		bus:                    bus,
		robotMode:              config.RobotMode,
		cameraMode:             config.CameraMode,
		activateCameraTransfer: config.ActivateCameraTransfer,
	}
}

func (c *CameraManualWorker) WorkerID() string {
	return cameraManualWorkerID
}

func (c *CameraManualWorker) Start() {
	for {
		cameraImage := gocv.NewMat()
		if !c.camera.Read(&cameraImage) {
			break
		}

		sourceImg, err := cameraImage.ToImage()
		if err != nil {
			fmt.Println("[CameraManualWorker] cameraImage.ToImag err: ", err.Error())
			return
		}

		b := sourceImg.Bounds()
		img := image.NewRGBA(b)
		draw.Draw(img, b, sourceImg, b.Min, draw.Src)

		if c.cameraMode == types.CAMERA_MODE__OBJECT_DETECTIVE {
			buf := bytes.NewBuffer([]byte{})
			err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 75})
			if err != nil {
				fmt.Println("[CameraManualWorker] jpeg.Encode err: ", err.Error())
				return
			}
			resp, err := c.cli.DetectionObject(buf.Bytes())
			if err != nil {
				fmt.Println("[CameraManualWorker] cli.DetectionObject request err: ", err)
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
				fmt.Println("[CameraManualWorker] jpeg.Encode err: ", err.Error())
				return
			}

			err = c.cli.CameraTransfer(buf.Bytes())
			if err != nil {
				fmt.Println("[CameraManualWorker] cli.CameraTransfer push err: ", err)
				return
			}
		}
	}
}

func (c *CameraManualWorker) Restart() error {
	return nil
}

func (c *CameraManualWorker) Stop() error {
	return c.camera.Close()
}
