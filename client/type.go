package client

type ObjectDetectionBody struct {
	Image []byte
}

type DetectivedObject struct {
	Class       float32
	Box         []float32
	Probability float32
}
