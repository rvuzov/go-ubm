package bclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type (
	Client struct {
		Addr string
	}

	UMetric struct {
		UserID string `json:"user"`
		Key    string `json:"key"`
		Value  int    `json:"value"`
	}

	ULog struct {
		UserID string      `json:"user"`
		Key    string      `json:"key"`
		Value  interface{} `json:"value"`
	}
)

const (
	PushMetricMethod = "push.metric"
	PushLogMethod    = "push.log"
)

func New(addr string) Client {
	client := Client{
		Addr: addr,
	}
	return client
}

func (c *Client) PushMetric(userID string, key string, value int) error {
	params := url.Values{}
	params.Add("userID", userID)
	params.Add("key", key)
	params.Add("value", strconv.Itoa(value))

	uri := c.MakeURI(PushMetricMethod, params.Encode())

	resp, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Error: %d", resp.StatusCode))
	}

	return nil
}

func (c *Client) PushLog(userID string, key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	svalue := string(bytes[:])

	params := url.Values{}
	params.Add("userID", userID)
	params.Add("key", key)
	params.Add("value", svalue)

	uri := c.MakeURI(PushLogMethod, params.Encode())

	resp, err := http.Get(uri)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Error: %d", resp.StatusCode))
	}

	return nil
}

func (c *Client) MakeURI(method, params string) string {
	return fmt.Sprintf("http://%s/%s?%s", c.Addr, method, params)
}
