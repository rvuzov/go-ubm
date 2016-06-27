package query

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func SendQuery(APIAddr string, msg interface{}) (responce string, code int) {
	container, err := NewAPIContainer(msg)
	if err != nil {
		return
	}

	query, err := json.Marshal(container)
	if err != nil {
		return err.Error(), 0
	}

	form := url.Values{}
	form.Add("query", string(query[:]))

	resp, err := http.PostForm(APIAddr, form)
	if err != nil {
		return err.Error(), code
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	sbody, err := strconv.Unquote(string(body[:]))
	if err != nil {
		return "", 0
	}
	return sbody, resp.StatusCode
}
