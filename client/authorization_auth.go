package client

func (c *RobotClient) AuthorizationAuth(key []byte) (token []byte, err error) {
	stat := c.sess.Call("/authorization/auth", key, &token).Status()
	if !stat.OK() {
		return nil, stat.Cause()
	}
	return
}
