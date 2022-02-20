package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpClient struct {
	c        http.Client
	ApiToken string
	BaseUrl  string
}

func (c *HttpClient) Get(url string) (resp *http.Response, err error) {

	url = c.BaseUrl + url

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err

	}

	return c.Do(req)
}

func (c *HttpClient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {

	url = c.BaseUrl + url

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err

	}

	req.Header.Set("Content-Type", contentType)

	return c.Do(req)

}

func (c *HttpClient) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

func (c *HttpClient) PatchForm(url string, data url.Values) (resp *http.Response, err error) {
	return c.Patch(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

func (c *HttpClient) Patch(url, contentType string, body io.Reader) (resp *http.Response, err error) {

	url = c.BaseUrl + url

	req, err := http.NewRequest("PATCH", url, body)
	if err != nil {
		return nil, err

	}

	req.Header.Set("Content-Type", contentType)

	return c.Do(req)

}

func (c *HttpClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+c.ApiToken)
	return c.c.Do(req)
}

func PostTimelog(date time.Time, description string) PostTimelogResponse {

	client := Client()

	layout := "2006-01-02T15:04"

	resp, err := client.PostForm("timelog", url.Values{"description": {description}, "range": {date.Format(layout) + "/" + date.Format(layout)}})

	if err != nil {
		fmt.Println(err)
	}

	var res PostTimelogResponse

	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		fmt.Println(err)
		return PostTimelogResponse{}
	}
	return res
}

func PatchTimelog(id string, date time.Time, description string) PostTimelogResponse {
	client := Client()

	layout := "2006-01-02T15:04"

	resp, err := client.PatchForm("timelog/"+id, url.Values{"description": {description}, "range": {date.Format(layout) + "/" + date.Format(layout)}})

	if err != nil {
		fmt.Println(err)
	}

	var res PostTimelogResponse

	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		fmt.Println(err)
		return PostTimelogResponse{}
	}

	return res
}

func Client() HttpClient {
	token := viper.GetString("access_token")

	client := HttpClient{ApiToken: token, BaseUrl: "https://api.nanda.io/api/v1/"}
	return client
}
