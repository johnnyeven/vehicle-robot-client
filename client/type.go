package client

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
