/*
 How to run
 Pass serial port to use as the first param:

        go run examples/firmata_servo.go /dev/ttyACM0
*/

package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor("/dev/cu.usbmodem14201")
	servoHorizon := gpio.NewServoDriver(firmataAdaptor, "8")
	servoVertical := gpio.NewServoDriver(firmataAdaptor, "9")

	work := func() {
		var horizonAngleNumber, verticalAngleNumber uint8
		for {
			fmt.Println("input horizonAngle and verticalAngle")
			fmt.Scanln(&horizonAngleNumber, &verticalAngleNumber)
			servoHorizon.Move(horizonAngleNumber)
			servoVertical.Move(verticalAngleNumber)
		}
	}

	robot := gobot.NewRobot("servoBot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{servoHorizon, servoVertical},
		work,
	)

	robot.Start()
}
