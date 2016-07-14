package metric

import (
	"../../../ubm"
	"../service"
)

type (
	CompareStatement struct {
		Metric   string `json:"metric"`
		Operator string `json:"operator"`
		Value    int    `json:"value"`
	}
	MetricFindUsers []CompareStatement

	UserMetric struct {
		UserID  string         `json:"userID"`
		Metrics map[string]int `json:"metrics"`
	}
	MetricFindUsersResponse []UserMetric
)

func (msg *MetricFindUsers) Receive() interface{} {

	var cmp ubm.CompareStatement
	var compareStatements []ubm.CompareStatement

	for _, statement := range *msg {
		cmp.Metric = statement.Metric
		cmp.Operator = statement.Operator
		cmp.Value = statement.Value
		compareStatements = append(compareStatements, cmp)
	}

	usersMetrics, err := ubm.Metrics.FindUsers(compareStatements)

	if err != nil {
		return service.NewError(err.Error())
	}

	var tmpMetric UserMetric
	var response MetricFindUsersResponse
	for _, item := range usersMetrics {
		tmpMetric.UserID = item.UserID
		tmpMetric.Metrics = item.Metrics
		response = append(response, tmpMetric)
	}

	return response

}
