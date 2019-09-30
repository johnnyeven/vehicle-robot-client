package client

func (c *RobotClient) CameraTransfer(frame []byte) error {
	request := CameraRequest{
		Frame: frame,
	}
	cmd := c.sess.Call("/camera/transfer", request, nil)
	if !cmd.Status().OK() {
		return cmd.Status().Cause()
	}
	return nil
}
