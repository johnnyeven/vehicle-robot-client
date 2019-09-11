package client

func (c *RobotClient) DetectionObject(frame []byte) ([]DetectivedObject, error) {
	result := make([]DetectivedObject, 0)
	stat := c.sess.Call("/detection/object", frame, &result).Status()
	if !stat.OK() {
		return nil, stat.Cause()
	}

	return result, nil
}

func (c *RobotClient) CameraTransfer(frame []byte) error {
	stat := c.sess.Push("/camera/transfer", frame)
	if !stat.OK() {
		return stat.Cause()
	}
	return nil
}
