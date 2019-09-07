package client

func (c *RobotClient) DetectionObject(req ObjectDetectionBody) ([]DetectivedObject, error) {
	result := make([]DetectivedObject, 0)
	stat := c.sess.Call("", req, &result).Status()
	if !stat.OK() {
		return nil, stat.Cause()
	}

	return result, nil
}
