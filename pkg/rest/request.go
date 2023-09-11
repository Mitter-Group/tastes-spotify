package rest

import (
	"context"
	"net/http"
	"testing"
	"time"
)

type Request struct {
	Method             string
	URL                string
	Body               interface{}
	Headers            http.Header
	TimeoutInMillis    int
	Retries            int
	AuthorizationToken *string
	cached             bool
	cacheTTL           time.Duration
	client             *RestClient
	ctx                context.Context
	newRelicTrace      bool
	isMocked           bool
	t                  *testing.T
}

type Response struct {
	StatusCode  int
	Headers     http.Header
	RawResponse *http.Response
	BodyBytes   []byte
	Duration    int64
	Error       error
}

func (r *Request) WithHeader(key, value string) *Request {
	if r.Headers == nil {
		r.Headers = make(http.Header)
	}
	r.Headers.Set(key, value)
	return r
}

func (r *Request) WithHeaders(headers http.Header) *Request {
	if r.Headers == nil {
		r.Headers = make(http.Header)
	}
	for k, v := range headers {
		r.Headers[k] = v
	}
	return r
}

func (r *Request) WithMapHeaders(headers map[string]string) *Request {
	for k, v := range headers {
		r.WithHeader(k, v)
	}
	return r
}

func (r *Request) WithNewRelicTrace(ctx context.Context) *Request {
	r.newRelicTrace = true
	r.ctx = ctx
	return r
}

func (r *Request) WithContext(ctx context.Context) *Request {
	r.ctx = ctx
	return r
}

func (r *Request) WithTimeout(timeoutInMillis int) *Request {
	r.TimeoutInMillis = timeoutInMillis
	return r
}

// WithRetries sets the number of retries for the request
func (r *Request) WithRetries(retries int) *Request {
	r.Retries = retries
	return r
}

func (r *Request) WithCache(ttl time.Duration) *Request {
	r.cached = true
	r.cacheTTL = ttl
	return r
}

func (r *Request) WithAuthorizationToken(token string) *Request {
	r.AuthorizationToken = &token
	return r
}
