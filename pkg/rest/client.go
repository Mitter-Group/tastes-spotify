package rest

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/chunnior/spotify/pkg/cache"
	"github.com/gojek/heimdall/v7/httpclient"
)

type Client interface {
	Get(url string) *Request
	Put(url string, body interface{}) *Request
	Post(url string, body interface{}) *Request
	Request(method string, url string, body interface{}) *Request
}

type RestClient struct {
	*httpclient.Client
	ReqInterceptors []RequestInterceptor  //Does not work yet
	ResInterceptors []ResponseInterceptor //Does not work yet
	Config          Config
	Cache           cache.Spec
}

type Config struct {
	BaseURL         string
	TimeoutInMillis int
	DefaultHeaders  http.Header
	Retries         int
}

func NewDefaultRestClient() *RestClient {
	customTransport := http.DefaultTransport
	//nolint:gosec // need insecure TLS option for testing while noto resolves the certificate issue
	customTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := http.Client{
		Transport: customTransport,
	}

	return &RestClient{
		Cache:  cache.NewMemoryCache("restclient", 30, 3600, false),
		Client: httpclient.NewClient(httpclient.WithHTTPClient(&client)),
	}
}

func NewCustomRestClient(cfg Config) *RestClient {
	customTransport := http.DefaultTransport
	//nolint:gosec // need insecure TLS option for testing while noto resolves the certificate issue
	customTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := http.Client{
		Timeout:   time.Duration(cfg.TimeoutInMillis) * time.Millisecond,
		Transport: customTransport,
	}

	return &RestClient{
		Client: httpclient.NewClient(httpclient.WithHTTPClient(&client)),
		Config: cfg,
	}
}

// RequestInterceptor is a function that can be used to modify the request or the response
type RequestInterceptor func(request *Request)

func (c *RestClient) WithRequestInterceptors(i ...RequestInterceptor) *RestClient {
	c.ReqInterceptors = append(c.ReqInterceptors, i...)
	return c
}

// ResponseInterceptor is a function that can be used to modify the response and executes after the response is received
type ResponseInterceptor func(response *Response)

func (c *RestClient) WithResponseInterceptors(i ...ResponseInterceptor) *RestClient {
	c.ResInterceptors = append(c.ResInterceptors, i...)
	return c
}

func (c *RestClient) Get(url string) *Request {
	return &Request{
		Method: http.MethodGet,
		URL:    url,
		client: c,
	}
}

func (c *RestClient) Put(url string, body interface{}) *Request {
	return &Request{
		Method: http.MethodPut,
		URL:    url,
		Body:   body,
		client: c,
	}
}

func (c *RestClient) Post(url string, body interface{}) *Request {
	return &Request{
		Method: http.MethodPost,
		URL:    url,
		Body:   body,
		client: c,
	}
}

func (c *RestClient) Request(method string, url string, body interface{}) *Request {
	return &Request{
		Method: method,
		URL:    url,
		Body:   body,
		client: c,
	}
}

// MapTo unmarshalls r.Body into bindTo. A pointer to the struct MUST be passed in.
func (r *Response) MapTo(bindTo interface{}) error {
	if r.BodyBytes == nil {
		return errors.New("response body is nil")
	}
	err := json.Unmarshal(r.BodyBytes, bindTo)
	if err != nil {
		return err
	}
	return nil
}
