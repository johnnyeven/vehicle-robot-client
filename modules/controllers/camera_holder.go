package controllers

import (
	"fmt"
	"gobot.io/x/gobot/drivers/gpio"
)

func CameraHolderController(servoHorizon *gpio.ServoDriver, servoVertical *gpio.ServoDriver) {
	var horizonAngleNumber, verticalAngleNumber uint8
	for {
		fmt.Println("input horizonAngle and verticalAngle")
		_, err := fmt.Scanln(&horizonAngleNumber, &verticalAngleNumber)
		if err != nil {
			fmt.Println("[CameraHolderController] fmt.Scanln err: ", err)
			continue
		}
		err = servoHorizon.Move(horizonAngleNumber)
		if err != nil {
			fmt.Printf("[CameraHolderController] servoHorizon.Move err: %v, angle: %d\n", err, horizonAngleNumber)
			continue
		}
		err = servoVertical.Move(verticalAngleNumber)
		if err != nil {
			fmt.Printf("[CameraHolderController] servoVertical.Move err: %v, angle: %d\n", err, verticalAngleNumber)
		}
	}
}
