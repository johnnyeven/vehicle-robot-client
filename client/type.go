package client

import "github.com/johnnyeven/vehicle-robot-client/constants"

type BroadcastRequest struct {
	Port uint16 `json:"port"`
}

type DetectivedObject struct {
	Class       float32   `json:"class"`
	Label       string    `json:"label"`
	Box         []float32 `json:"box"`
	Probability float32   `json:"probability"`
}

type CameraRequest struct {
	Frame []byte `json:"frame"`
}

type AuthRequest struct {
	Key string `json:"key"`
}

type PowerMovingRequest struct {
	Direction constants.MovingDirection `json:"direction"`
	Speed     float64                   `json:"speed"`
}

type CameraHolderRequest struct {
	Direction constants.HolderDirection `json:"direction"`
	Offset    float64                   `json:"offset"`
}
