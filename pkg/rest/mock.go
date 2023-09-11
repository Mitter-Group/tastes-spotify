package rest

import (
	"sync"
	"testing"
)

type MockClient struct {
	BaseURL string
	test    *testing.T
}

// MockCall is a struct that contains the information of a request that was made
type MockCall struct {
	URL          string
	ResponseBody string
	RequestBody  interface{}
}

// NewDefaultMockClient creates a new mock client with the default base URL
func NewDefaultMockClient(t *testing.T) *MockClient {
	mockCalls = make(map[string]MockResponse)
	requestsDone = make(map[string]int)
	requestStack = make([]MockCall, 0)
	return &MockClient{
		test: t,
	}
}

// NewMockClient creates a new mock client with the given base URL
func NewMockClient(t *testing.T, baseURL string) *MockClient {
	mockCalls = make(map[string]MockResponse)
	requestsDone = make(map[string]int)
	requestStack = make([]MockCall, 0)
	return &MockClient{
		test:    t,
		BaseURL: baseURL,
	}
}

var mockCalls map[string]MockResponse
var requestStack []MockCall
var requestsDone map[string]int
var mutex sync.Mutex

// MockResponse is a struct that contains the information of a mocked response
type MockResponse struct {
	StatusCode    int
	JSONBody      string
	InternalError error
}

func (c *MockClient) Get(url string) *Request {
	return &Request{
		Method:   "GET",
		URL:      url,
		isMocked: true,
		t:        c.test,
	}
}

func (c *MockClient) Post(url string, body interface{}) *Request {
	return &Request{
		Body:     body,
		Method:   "POST",
		URL:      url,
		isMocked: true,
		t:        c.test,
	}
}

func (c *MockClient) Put(url string, body interface{}) *Request {
	return &Request{
		Method:   "PUT",
		Body:     body,
		URL:      url,
		isMocked: true,
		t:        c.test,
	}
}

func (c *MockClient) Delete(url string) *Request {
	return &Request{
		Method:   "DELETE",
		URL:      url,
		isMocked: true,
		t:        c.test,
	}
}

func (c *MockClient) Request(method string, url string, body interface{}) *Request {
	return &Request{
		Method:   method,
		URL:      url,
		Body:     body,
		isMocked: true,
		t:        c.test,
	}
}

// SetMockCall is a function that can be used to mock the response of a request. DONT include the baseURL of the client
func (c *MockClient) SetMockCall(method string, url string, response MockResponse) {
	if mockCalls == nil {
		mockCalls = make(map[string]MockResponse)
	}
	mockCalls[method+url] = response
}

func (c *MockClient) CountRequestsDone(method string, url string) int {
	if res, ok := requestsDone[method+url]; ok {
		return res
	}
	return 0
}

func (c *MockClient) GetRequestStack() []MockCall {
	return requestStack
}

func (c *MockClient) ClearMockCalls() {
	mockCalls = make(map[string]MockResponse)

}

func handleMockRequest(r *Request, url, method string) *Response {
	if r.isMocked {
		mutex.Lock()
		defer mutex.Unlock()
		mock, ok := mockCalls[method+url]
		if !ok {
			r.t.Logf("No mock found for %s", method+url)
			r.t.Fail()
		}

		if mock.InternalError != nil {
			return &Response{
				StatusCode: 500,
				Error:      mockCalls[method+url].InternalError,
			}
		}
		requestsDone[method+url] = requestsDone[method+url] + 1
		requestStack = append(requestStack, MockCall{
			URL:          url,
			ResponseBody: mock.JSONBody,
			RequestBody:  r.Body,
		})

		bodyBytes := []byte(mock.JSONBody)
		return &Response{
			StatusCode:  mockCalls[method+url].StatusCode,
			Error:       nil,
			RawResponse: nil,
			BodyBytes:   bodyBytes,
		}
	}
	return nil
}
