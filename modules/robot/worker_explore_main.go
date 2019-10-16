package robot

import (
	"bytes"
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/vehicle-robot-client/client"
	"github.com/johnnyeven/vehicle-robot-client/global"
	"github.com/johnnyeven/vehicle-robot-client/modules"
	bus2 "github.com/mustafaturan/bus"
	"github.com/sirupsen/logrus"
	"image"
	"image/jpeg"
)

const (
	exploreMainWorkerID             = "explore-main-worker"
	cameraCaptureResultTopic        = "camera.capture.result"
	cameraCaptureResultEventHandler = "camera-capture-result-handler"
)

type ExploreMainWorker struct {
	config *global.RobotConfiguration
	bus    *bus.MessageBus
	cli    *client.RobotClient
}

func NewExploreMainWorker(robot *Robot, config *global.RobotConfiguration, bus *bus.MessageBus, cli *client.RobotClient) *ExploreMainWorker {
	return &ExploreMainWorker{
		bus: bus,
		cli: cli,
	}
}

func (e *ExploreMainWorker) WorkerID() string {
	return exploreMainWorkerID
}

func (e *ExploreMainWorker) Start() {
	e.bus.RegisterTopic(cameraCaptureResultTopic)
	e.bus.RegisterHandler(cameraCaptureResultEventHandler, cameraCaptureResultTopic, func(evt *bus2.Event) {
		sourceImg, ok := evt.Data.(*image.RGBA)
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
		e.bus.Emit(cameraCaptureTopic, nil, "")
	})
}

func (e *ExploreMainWorker) Restart() error {
	return nil
}

func (e *ExploreMainWorker) Stop() error {
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
