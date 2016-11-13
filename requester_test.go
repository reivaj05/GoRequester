package requester

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RequesterTestSuite struct {
	suite.Suite
	timesCalled  int
	getServer    *httptest.Server
	postServer   *httptest.Server
	processName  string
	getResponse  string
	postResponse string
	assert       *assert.Assertions
}

func (suite *RequesterTestSuite) SetupSuite() {
	suite.processName = "requester.test"
	suite.assert = assert.New(suite.T())

	suite.generateResponses()

	suite.getServer = httptest.NewServer(http.HandlerFunc(suite.getHandler))
	suite.postServer = httptest.NewServer(http.HandlerFunc(suite.postHandler))
}

func (suite *RequesterTestSuite) generateResponses() {
	suite.getResponse = `{"get": "Get response"}`
	suite.postResponse = `{"post": "Post response"}`
}

func (suite *RequesterTestSuite) getHandler(
	w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, suite.getResponse)
}

func (suite *RequesterTestSuite) postHandler(
	w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, suite.postResponse)
}

func (suite *RequesterTestSuite) TearDownSuite() {
	suite.getServer.Close()
	suite.postServer.Close()
}

func (suite *RequesterTestSuite) TestNew() {
	requesterObj := New()
	suite.assert.NotNil(requesterObj)
}

func (suite *RequesterTestSuite) TestMakeRequest() {
	suite.getTest()
	suite.postTest()
}

func (suite *RequesterTestSuite) getTest() {
	config := &RequestConfig{
		URL:    suite.getServer.URL,
		Method: "GET",
		Values: url.Values{},
		Headers: map[string]string{
			"headerTest1": "headerTestValue2,",
			"headerTest2": "headerTestValue2,",
		},
	}
	suite.successTest(config, suite.getResponse)
	suite.failTest(config)
}

func (suite *RequesterTestSuite) postTest() {
	config := &RequestConfig{
		URL:    suite.postServer.URL,
		Method: "POST",
		Values: nil,
	}
	suite.successTest(config, suite.postResponse)
	suite.failTest(config)
}

func (suite *RequesterTestSuite) successTest(
	config *RequestConfig, response string) {
	requesterObj := New()
	responseRetrieved, status, err := requesterObj.MakeRequest(config)
	suite.assert.Equal(response, responseRetrieved)
	suite.assert.Equal(status, http.StatusOK)
	suite.assert.Nil(err)
}

func (suite *RequesterTestSuite) failTest(config *RequestConfig) {
	config.URL = "//"
	requesterObj := New()
	response, status, err := requesterObj.MakeRequest(config)
	suite.assert.Equal(response, "")
	suite.assert.Equal(status, -1)
	suite.assert.NotNil(err)
}

func TestRequester(test *testing.T) {
	suite.Run(test, new(RequesterTestSuite))
}
