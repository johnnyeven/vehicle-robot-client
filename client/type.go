package client

import "github.com/johnnyeven/vehicle-robot-client/constants"

type DetectivedObject struct {
	Class       float32
	Box         []float32
	Probability float32
}

type AuthRequestHeader struct {
	Token string
}

type CameraRequest struct {
	AuthRequestHeader
	Frame []byte
}

type AuthRequest struct {
	Key string
}

type PowerMovingRequest struct {
	Direction constants.MovingDirection `json:"direction"`
	Speed     float64                   `json:"speed"`
}
