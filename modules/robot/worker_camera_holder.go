package robot

import (
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/vehicle-robot-client/client"
	"github.com/johnnyeven/vehicle-robot-client/global"
	bus2 "github.com/mustafaturan/bus"
	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

const (
	cameraHolderWorkerID           = "camera-holder-worker"
	CameraHolderTopic              = "camera.holder"
	cameraHolderEventHandler       = "camera-holder-handler"
	MaxAngle                 uint8 = 180
	CentreAngle              uint8 = 90
)

type CameraHolderWorker struct {
	servoHorizon  *gpio.ServoDriver
	servoVertical *gpio.ServoDriver
	bus           *bus.MessageBus

	currentHorizonAngle  uint8
	currentVerticalAngle uint8
}

func NewCameraHolderWorker(robot *Robot, bus *bus.MessageBus, config *global.RobotConfiguration) *CameraHolderWorker {
	var firmataAdaptor *firmata.Adaptor
	var ok bool
	conn := robot.GetConnection(config.FirmataConnectionName)
	if conn == nil {
		firmataAdaptor = firmata.NewAdaptor(config.ArduinoDeviceID)
		firmataAdaptor.SetName(config.FirmataConnectionName)
		robot.AddConnection(firmataAdaptor)
	} else {
		if firmataAdaptor, ok = conn.(*firmata.Adaptor); !ok {
			logrus.Panicf("[CameraHolderWorker] 连接器已存在，但并不是 *firmata.Adaptor 类型")
		}
	}

	servoHorizon := gpio.NewServoDriver(firmataAdaptor, config.ServoHorizonPin)
	servoHorizon.SetName(config.ServoHorizonName)
	servoVertical := gpio.NewServoDriver(firmataAdaptor, config.ServoVerticalPin)
	servoVertical.SetName(config.ServoVerticalName)

	robot.AddDevice(servoHorizon, servoVertical)

	return &CameraHolderWorker{
		servoHorizon:         servoHorizon,
		servoVertical:        servoVertical,
		bus:                  bus,
		currentHorizonAngle:  CentreAngle,
		currentVerticalAngle: CentreAngle,
	}
}

func (c *CameraHolderWorker) WorkerID() string {
	return cameraHolderWorkerID
}

func (c *CameraHolderWorker) Start() {
	logrus.Infof("[CameraHolderWorker] Init servos to center angle: %d", CentreAngle)
	err := c.servoHorizon.Move(CentreAngle)
	if err != nil {
		logrus.Errorf("[CameraHolderWorker] horizon servo move failed with err: %v", err)
		return
	}
	err = c.servoVertical.Move(CentreAngle)
	if err != nil {
		logrus.Errorf("[CameraHolderWorker] vertical servo move failed with err: %v", err)
		return
	}

	c.bus.RegisterTopic(CameraHolderTopic)
	c.bus.RegisterHandler(cameraHolderEventHandler, CameraHolderTopic, func(e *bus2.Event) {
		var err error
		if evt, ok := e.Data.(*client.CameraHolderRequest); ok {
			c.currentHorizonAngle = servoAngleChange(c.currentHorizonAngle, evt.HorizonOffset)
			err = c.servoHorizon.Move(c.currentHorizonAngle)
			if err != nil {
				logrus.Errorf("[CameraHolderWorker] camera-holder-handler servoHorizon.Move err: %v, angle: %d, event: %+v", err, c.currentHorizonAngle, evt)
			}

			c.currentVerticalAngle = servoAngleChange(c.currentVerticalAngle, evt.VerticalOffset)
			err = c.servoVertical.Move(c.currentVerticalAngle)
			if err != nil {
				logrus.Errorf("[CameraHolderWorker] camera-holder-handler servoVertical.Move err: %v, angle: %d, event: %+v", err, c.currentVerticalAngle, evt)
			}
		} else {
			logrus.Errorf("[CameraHolderWorker] camera-holder-handler Data type err: %s", "not CameraHolderRequest struct")
		}
	})
}

func (c *CameraHolderWorker) Restart() error {
	return nil
}

func (c *CameraHolderWorker) Stop() error {
	c.bus.DeregisterHandler(cameraHolderEventHandler)
	c.servoVertical.Move(CentreAngle)
	c.servoHorizon.Move(CentreAngle)
	return nil
}

func servoAngleChange(current uint8, offset float64) uint8 {
	current = uint8(float64(current) + offset)
	if current < 0 {
		current = 0
	} else if current > MaxAngle {
		current = MaxAngle
	}
	return current
}
