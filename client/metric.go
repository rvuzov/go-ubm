package bclient

import (
	"./api/metric"
)

type (
	Cmp struct {
		Metric   string
		Operator string
		Value    int
	}

	UserMetric struct {
		UserID  string
		Metrics map[string]int
	}
)

func (c *Client) MetricGet(userID string, metrics []string) (result map[string]int, err error) {
	requestMessage := metric.MetricGet{
		UserID:  userID,
		Metrics: metrics,
	}

	result, err = requestMessage.Receive((*c).APIAddr)
	return
}

func (c *Client) MetricPush(userID string, key string, value int) (err error) {
	requestMessage := metric.MetricPush{
		UserID: userID,
		Key:    key,
		Value:  value,
	}

	err = requestMessage.Receive((*c).APIAddr)
	return
}

func (c *Client) MetricFindUsers(cmps []Cmp) (answer []UserMetric, err error) {
	var requestMessage metric.MetricFindUsers
	var tmpCmp metric.Cmp
	for _, cmp := range cmps {
		tmpCmp.Metric = cmp.Metric
		tmpCmp.Operator = cmp.Operator
		tmpCmp.Value = cmp.Value
		requestMessage = append(requestMessage, tmpCmp)
	}

	results, err := requestMessage.Receive((*c).APIAddr)
	if err != nil {
		return
	}

	var tmpMetric UserMetric
	for _, result := range results {
		tmpMetric.UserID = result.UserID
		tmpMetric.Metrics = result.Metrics
		answer = append(answer, tmpMetric)
	}

	return
}
