package client

func (c *RobotClient) DetectionObject(frame []byte) ([]DetectivedObject, error) {
	if c.sess == nil {
		return nil, nil
	}

	result := make([]DetectivedObject, 0)
	request := CameraRequest{
		Frame: frame,
	}
	stat := c.sess.Call("/detection/object", request, &result).Status()
	if !stat.OK() {
		return nil, stat.Cause()
	}

	return result, nil
}
