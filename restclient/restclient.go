package restclient

import (
	"net/http"
	"encoding/json"
	"bytes"
	//"fmt"
)

type HTTPClient interface {
	Do(req * http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

// Post sends a post request to the URL with the body
func Post(url string, body map[string]interface{}, headers http.Header) (*http.Response, error) {
	jsonBytes, err := json.Marshal(body)
	//fmt.Print(string(jsonBytes))
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	request.Header = headers
	return Client.Do(request)
}

func Get(url string, headers http.Header) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header = headers
	return Client.Do(request)
}