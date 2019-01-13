package requester

import (
	"bytes"
	"io/ioutil"

	"net/http"
	"net/url"
)

// Requester struct holds an http client to make calls
type Requester struct {
	client *http.Client
}

// RequestConfig Holds configuration values to make a request
type RequestConfig struct {
	URL     string
	Method  string
	Values  url.Values
	Body    []byte
	Headers map[string]string
}

type createRequest func(config *RequestConfig) (*http.Request, error)

var (
	requestsMaker = map[string]createRequest{
		"POST": createPostRequest,
		"PUT":  createPostRequest,
		"GET":  createGetRequest,
	}
)

// New creates a new requester object
func New() *Requester {
	return &Requester{
		client: &http.Client{},
	}
}

// MakeRequest Makes a HTTP request to the config passed as parameter
func (requester *Requester) MakeRequest(
	config *RequestConfig) (string, int, error) {

	response, err := requester._makeRequest(config)
	if err != nil {
		return "", -1, err
	}
	return parseResponse(response)
}

func (requester *Requester) _makeRequest(
	config *RequestConfig) (*http.Response, error) {

	request, _ := requestsMaker[config.Method](config)
	addHeaders(request, config.Headers)
	return requester.client.Do(request)
}

func createGetRequest(config *RequestConfig) (*http.Request, error) {
	return http.NewRequest(config.Method,
		config.URL+"?"+config.Values.Encode(), nil)
}

func createPostRequest(config *RequestConfig) (*http.Request, error) {
	return http.NewRequest(config.Method,
		config.URL+"?"+config.Values.Encode(), bytes.NewBuffer(config.Body))
}

func addHeaders(request *http.Request, headers map[string]string) {
	for header, value := range headers {
		request.Header.Add(header, value)
	}
}

func parseResponse(response *http.Response) (string, int, error) {
	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	return string(body), response.StatusCode, err
}
