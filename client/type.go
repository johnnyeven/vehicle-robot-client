package client

type DetectivedObject struct {
	Class       float32
	Box         []float32
	Probability float32
}

type AuthRequest struct {
	Token string
}

type CameraRequest struct {
	AuthRequest
	Frame []byte
}
