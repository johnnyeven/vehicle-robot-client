package client

func (c *RobotClient) CameraTransfer(frame []byte) error {
	request := CameraRequest{
		Frame: frame,
	}
	stat := c.sess.Push("/camera/transfer", request)
	if !stat.OK() {
		return stat.Cause()
	}
	return nil
}
