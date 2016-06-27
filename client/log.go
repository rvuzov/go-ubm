package bclient

import "./api/log"

func (c *Client) PushLog(userID string, key string, value interface{}) (err error) {
	requestMessage := log.LogPush{
		UserID: userID,
		Key:    key,
		Value:  value,
	}

	err = requestMessage.Receive(c.APIAddr)
	return
}
