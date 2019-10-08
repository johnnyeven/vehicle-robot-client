package client

func (c *RobotClient) AuthorizationAuth(key string) (token []byte, err error) {
	if c.sess == nil {
		return
	}

	request := AuthRequest{
		Key: key,
	}
	stat := c.sess.Call("/authorization/auth", request, &token).Status()
	if !stat.OK() {
		err = stat.Cause()
	}
	return
}
