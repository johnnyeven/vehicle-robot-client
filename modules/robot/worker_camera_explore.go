package robot

import (
	"fmt"
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/libtools/courier/enumeration"
	"github.com/johnnyeven/vehicle-robot-client/client"
	"github.com/johnnyeven/vehicle-robot-client/global"
	bus2 "github.com/mustafaturan/bus"
	"github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
)

const (
	cameraExploreWorkerID     = "camera-explore-worker"
	cameraCaptureTopic        = "camera.capture"
	cameraCaptureEventHandler = "camera-capture-handler"
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
	c.bus.RegisterTopic(cameraCaptureTopic)
	c.bus.RegisterHandler(cameraCaptureEventHandler, cameraCaptureTopic, func(e *bus2.Event) {
		cameraImage := gocv.NewMat()
		if !c.camera.Read(&cameraImage) {
			return
		}
		img, err := cameraImage.ToImage()
		if err != nil {
			fmt.Println("[CameraExploreWorker] cameraImage.ToImag err: ", err.Error())
			return
		}
		c.bus.Emit(cameraCaptureResultTopic, img, "")
	})
	c.bus.Emit(cameraCaptureTopic, nil, "")
}

func (c *CameraExploreWorker) Restart() error {
	return nil
}

func (c *CameraExploreWorker) Stop() error {
	c.bus.DeregisterHandler(cameraCaptureEventHandler)
	return c.camera.Close()
}
