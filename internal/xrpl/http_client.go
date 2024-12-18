package xrpl

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type HTTPClient struct {
	BaseURL string
	Client  *http.Client
}

// create a new HTTP client with the provided URL
func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

// send a POST request to the xrpl api with the given payload
func (hc *HTTPClient) Post(endpoint string, payload interface{}) ([]byte, error) {
	url := hc.BaseURL + endpoint

	reqBody, err := json.Marshal(payload) // convert payload to JSON
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody)) // create a new POST request
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := hc.Client.Do(req) // send the request
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body) // read the response body
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("non-200 status code received: " + resp.Status)
	}

	return respBody, nil
}
