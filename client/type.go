package client

import (
	"github.com/johnnyeven/service-vehicle-robot/constants/types"
)

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
	Direction types.MovingDirection `json:"direction"`
	Speed     float64               `json:"speed"`
}

type CameraHolderRequest struct {
	HorizonOffset  float64 `json:"horizonOffset"`
	VerticalOffset float64 `json:"verticalOffset"`
}
