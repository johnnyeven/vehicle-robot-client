package robot

import (
	"bytes"
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/vehicle-robot-client/client"
	"github.com/johnnyeven/vehicle-robot-client/global"
	"github.com/johnnyeven/vehicle-robot-client/modules"
	"github.com/sirupsen/logrus"
	"image"
	"image/jpeg"
)

const (
	exploreMainWorkerID             = "explore-main-worker"
	cameraCaptureResultTopic        = "camera.Capture.result"
	cameraCaptureResultEventHandler = "camera-Capture-result-handler"
)

type ExploreMainWorker struct {
	config *global.RobotConfiguration
	bus    *bus.MessageBus
	cli    *client.RobotClient
	quit   chan struct{}

	cameraWorker       *CameraExploreWorker
	cameraHolderWorker *CameraHolderWorker
	attitudeWorker     *AttitudeGY85Worker
	distanceWorker     *DistanceHCSR04Worker
	powerWorker        *PowerWorker
}

func NewExploreMainWorker(robot *Robot, config *global.RobotConfiguration, bus *bus.MessageBus, cli *client.RobotClient) *ExploreMainWorker {
	cameraWorker, ok := robot.GetWorker(cameraExploreWorkerID).(*CameraExploreWorker)
	if !ok {
		logrus.Panicf("[ExploreMainWorker] robot.GetWorker(cameraExploreWorkerID).(*CameraExploreWorker) error")
	}
	cameraHolderWorker, ok := robot.GetWorker(cameraHolderWorkerID).(*CameraHolderWorker)
	if !ok {
		logrus.Panicf("[ExploreMainWorker] robot.GetWorker(cameraHolderWorkerID).(*CameraHolderWorker) error")
	}
	attitudeWorker, ok := robot.GetWorker(attitudeGY85WorkerID).(*AttitudeGY85Worker)
	if !ok {
		logrus.Panicf("[ExploreMainWorker] robot.GetWorker(attitudeGY85WorkerID).(*AttitudeGY85Worker) error")
	}
	distanceWorker, ok := robot.GetWorker(distanceHCSR04WorkerID).(*DistanceHCSR04Worker)
	if !ok {
		logrus.Panicf("[ExploreMainWorker] robot.GetWorker(distanceHCSR04WorkerID).(*DistanceHCSR04Worker) error")
	}
	powerWorker, ok := robot.GetWorker(powerWorkerID).(*PowerWorker)
	if !ok {
		logrus.Panicf("[ExploreMainWorker] robot.GetWorker(powerWorkerID).(*PowerWorker) error")
	}

	return &ExploreMainWorker{
		config: config,
		bus:    bus,
		cli:    cli,
		quit:   make(chan struct{}),

		cameraWorker:       cameraWorker,
		cameraHolderWorker: cameraHolderWorker,
		attitudeWorker:     attitudeWorker,
		distanceWorker:     distanceWorker,
		powerWorker:        powerWorker,
	}
}

func (e *ExploreMainWorker) WorkerID() string {
	return exploreMainWorkerID
}

func (e *ExploreMainWorker) Start() {
Run:
	for {
		select {
		case <-e.quit:
			break Run
		default:
			img, err := e.cameraWorker.Capture()
			if err != nil {
				return
			}
			sourceImg, ok := img.(*image.RGBA)
			if !ok {
				logrus.Error("[ExploreMainWorker] e.Data.(*image.RGBA) not *image.RGBA")
				return
			}
			detectived, err := e.objectDetective(sourceImg)
			if err != nil {
				logrus.Errorf("[ExploreMainWorker] ExploreMainWorker.objectDetective err: %v", err)
			}

			if e.config.ActivateCameraTransfer.True() {
				_ = e.transferScreen(sourceImg, detectived)
			}
		}
	}
}

func (e *ExploreMainWorker) Restart() error {
	return nil
}

func (e *ExploreMainWorker) Stop() error {
	e.quit <- struct{}{}
	return nil
}

func (e *ExploreMainWorker) objectDetective(sourceImg *image.RGBA) (resp []client.DetectivedObject, err error) {
	buf := bytes.NewBuffer([]byte{})
	err = jpeg.Encode(buf, sourceImg, &jpeg.Options{Quality: 75})
	if err != nil {
		logrus.Errorf("[ExploreMainWorker] jpeg.Encode err: ", err.Error())
		return
	}
	resp, err = e.cli.DetectionObject(buf.Bytes())
	return
}

func (e *ExploreMainWorker) transferScreen(sourceImg *image.RGBA, detectived []client.DetectivedObject) (err error) {
	for _, d := range detectived {
		x1 := float32(sourceImg.Bounds().Max.X) * d.Box[1]
		x2 := float32(sourceImg.Bounds().Max.X) * d.Box[3]
		y1 := float32(sourceImg.Bounds().Max.Y) * d.Box[0]
		y2 := float32(sourceImg.Bounds().Max.Y) * d.Box[2]

		modules.Rect(sourceImg, int(x1), int(y1), int(x2), int(y2), 4, modules.GetLabelColor(int(d.Class)))
		modules.DrawLabel(sourceImg, int(x1), int(y1), int(d.Class), d.Label)
	}

	buf := bytes.NewBuffer([]byte{})
	err = jpeg.Encode(buf, sourceImg, &jpeg.Options{Quality: 75})
	if err != nil {
		logrus.Errorf("[CameraExploreWorker] jpeg.Encode err: ", err.Error())
		return
	}

	err = e.cli.CameraTransfer(buf.Bytes())
	if err != nil {
		logrus.Errorf("[CameraExploreWorker] cli.CameraTransfer push err: ", err)
	}

	return
}
