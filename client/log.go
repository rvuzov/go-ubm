package bclient

import "./api/log"

func (c *Client) LogPush(userID string, key string, value interface{}) (err error) {
	requestMessage := log.LogPush{
		UserID: userID,
		Key:    key,
		Value:  value,
	}

	err = requestMessage.Receive(c.APIAddr)
	return
}
