package bclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type (
	Client struct {
		Addr    string
		APIAddr string
	}

	MetricPush struct {
		UserID string `json:"userID"`
		Key    string `json:"key"`
		Value  int    `json:"value"`
	}

	LogPush struct {
		UserID string      `json:"userID"`
		Key    string      `json:"key"`
		Value  interface{} `json:"value"`
	}
)

func NewClient(addr string) Client {
	return Client{
		Addr:    addr,
		APIAddr: fmt.Sprintf("http://%s/", addr),
	}
}

func (c *Client) PushMetric(userID string, key string, value int) error {
	msg := MetricPush{
		UserID: userID,
		Key:    key,
		Value:  value,
	}
	container, err := NewAPIContainer(msg)
	if err != nil {
		return err
	}

	response, code := c.SendAPIQuery(container)
	if code != 200 {
		return errors.New(fmt.Sprintf(`{error: "%s", code: %d`, response, code))
	}
	return nil
}

func (c *Client) PushLog(userID string, key string, value interface{}) error {
	msg := LogPush{
		UserID: userID,
		Key:    key,
		Value:  value,
	}

	container, err := NewAPIContainer(msg)
	if err != nil {
		return err
	}

	response, code := c.SendAPIQuery(container)
	if code != 200 {
		return errors.New(fmt.Sprintf(`{error: "%s", code: %d`, response, code))
	}
	return nil
}

func (c *Client) SendAPIQuery(msg interface{}) (response string, code int) {
	query, err := json.Marshal(msg)
	if err != nil {
		return err.Error(), 0
	}

	form := url.Values{}
	form.Add("query", string(query[:]))

	resp, err := http.PostForm(c.APIAddr, form)
	if err != nil {
		return err.Error(), code
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body[:]), resp.StatusCode
}
