package bclient

import (
	"./api/metric"
)

func (c *Client) GetMetric(userID string, metrics []string) (result map[string]int, err error) {
	requestMessage := metric.MetricGet{
		UserID:  userID,
		Metrics: metrics,
	}

	result, err = requestMessage.Receive((*c).APIAddr)
	return
}

func (c *Client) PushMetric(userID string, key string, value int) (err error) {
	requestMessage := metric.MetricPush{
		UserID: userID,
		Key:    key,
		Value:  value,
	}

	err = requestMessage.Receive((*c).APIAddr)
	return
}
