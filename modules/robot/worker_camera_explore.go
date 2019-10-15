package robot

import (
	"bytes"
	"fmt"
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/libtools/courier/enumeration"
	"github.com/johnnyeven/vehicle-robot-client/client"
	"github.com/johnnyeven/vehicle-robot-client/global"
	"github.com/johnnyeven/vehicle-robot-client/modules"
	"github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
	"image"
	"image/jpeg"
)

const (
	cameraExploreWorkerID = "camera-explore-worker"
)

type CameraExploreWorker struct {
	camera *gocv.VideoCapture
	cli    *client.RobotClient
	bus    *bus.MessageBus

	activateCameraTransfer enumeration.Bool
}

func NewCameraExploreWorker(robot *Robot, bus *bus.MessageBus, cli *client.RobotClient, config *global.RobotConfiguration) *CameraExploreWorker {
	camera, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		logrus.Panicf("[CameraExploreWorker] gocv.VideoCaptureDevice err: %v", err)
	}
	camera.Set(gocv.VideoCaptureFrameWidth, float64(config.CameraCaptureWidth))
	camera.Set(gocv.VideoCaptureFrameHeight, float64(config.CameraCaptureHeight))

	return &CameraExploreWorker{
		camera: camera,
		cli:    cli,
		bus:    bus,

		activateCameraTransfer: config.ActivateCameraTransfer,
	}
}

func (c *CameraExploreWorker) WorkerID() string {
	return cameraExploreWorkerID
}

func (c *CameraExploreWorker) Start() {
	for {
		cameraImage := gocv.NewMat()
		if !c.camera.Read(&cameraImage) {
			break
		}

		img, err := cameraImage.ToImage()
		if err != nil {
			fmt.Println("[CameraExploreWorker] cameraImage.ToImag err: ", err.Error())
			return
		}
		sourceImg, ok := img.(*image.RGBA)
		if !ok {
			logrus.Error("[CameraExploreWorker] img.(*image.RGBA) not *image.RGBA")
			return
		}

		buf := bytes.NewBuffer([]byte{})
		err = jpeg.Encode(buf, sourceImg, &jpeg.Options{Quality: 75})
		if err != nil {
			fmt.Println("[CameraExploreWorker] jpeg.Encode err: ", err.Error())
			return
		}
		resp, err := c.cli.DetectionObject(buf.Bytes())
		if err != nil {
			fmt.Println("[CameraExploreWorker] cli.DetectionObject request err: ", err)
			return
		}

		if c.activateCameraTransfer.True() {
			for _, detectived := range resp {
				x1 := float32(sourceImg.Bounds().Max.X) * detectived.Box[1]
				x2 := float32(sourceImg.Bounds().Max.X) * detectived.Box[3]
				y1 := float32(sourceImg.Bounds().Max.Y) * detectived.Box[0]
				y2 := float32(sourceImg.Bounds().Max.Y) * detectived.Box[2]

				modules.Rect(sourceImg, int(x1), int(y1), int(x2), int(y2), 4, modules.GetLabelColor(int(detectived.Class)))
				modules.DrawLabel(sourceImg, int(x1), int(y1), int(detectived.Class), detectived.Label)
			}

			buf := bytes.NewBuffer([]byte{})
			err = jpeg.Encode(buf, sourceImg, &jpeg.Options{Quality: 75})
			if err != nil {
				fmt.Println("[CameraExploreWorker] jpeg.Encode err: ", err.Error())
				return
			}

			err = c.cli.CameraTransfer(buf.Bytes())
			if err != nil {
				fmt.Println("[CameraExploreWorker] cli.CameraTransfer push err: ", err)
				return
			}
		}
	}
}

func (c *CameraExploreWorker) Restart() error {
	panic("implement me")
}

func (c *CameraExploreWorker) Stop() error {
	return c.camera.Close()
}
