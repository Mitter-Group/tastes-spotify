package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/chunnior/spotify/internal/util/log"
	"github.com/chunnior/spotify/pkg/tracing"
	"github.com/pkg/errors"
)

// Do perform the http request taking in consideration the fields of the request
// return the response
func (r *Request) Do() *Response {
	if r.isMocked {
		return handleMockRequest(r, r.URL, r.Method)
	}

	url := r.client.Config.BaseURL + r.URL

	if r.cached {
		log.Debug("Checking cache for url: ", url)
		_, cachedResponse := r.client.Cache.Get(r.ctx, r.Method+url)
		if cachedResponse != nil {
			return cachedResponse.(*Response)
		}
	}

	var req *http.Request
	var err error
	if r.Body != nil {
		var b []byte
		b, err = json.Marshal(r.Body)
		if err != nil {
			return &Response{
				StatusCode: 800,
				Error:      err,
			}
		}
		req, err = http.NewRequest(r.Method, url, bytes.NewBuffer(b))
	} else {
		req, err = http.NewRequest(r.Method, url, nil)
	}

	if r.newRelicTrace {
		segment := tracing.StartExternalSegmentFromRequest(r.ctx, req)
		defer segment.End()
	}

	if r.ctx != nil {
		req = req.WithContext(r.ctx)
	}

	if err != nil {
		return nil
	}

	if r.Headers != nil {
		for k, v := range r.Headers {
			headersConcat := strings.Join(v, ",")
			req.Header.Set(k, headersConcat)
		}
	}

	if r.AuthorizationToken != nil {
		req.Header.Set("Authorization", "Bearer "+*r.AuthorizationToken)
	}

	start := time.Now()
	res, err := r.client.Do(req)
	elapsed := time.Since(start).Milliseconds()

	if err != nil {
		statusCode := 500
		if res != nil {
			statusCode = res.StatusCode
		}
		return &Response{
			StatusCode: statusCode,
			Error:      errors.Wrap(err, "error on url: "+url),
		}
	}

	bodyBytes := []byte{}
	if res.Body != nil {
		bodyBytes, err = io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			return &Response{
				StatusCode: 800,
				Error:      errors.Wrap(err, "error reading response body"),
			}
		}
	}

	response := &Response{
		StatusCode:  res.StatusCode,
		Headers:     res.Header,
		RawResponse: res,
		//Body:        res.Body,
		BodyBytes: bodyBytes,
		Error:     nil,
		Duration:  elapsed,
	}

	if r.cached {
		log.DebugfWithContext(r.ctx, "Caching response for url: %s with ttl: %s", url, r.cacheTTL)
		r.client.Cache.SaveWithTTL(r.ctx, r.Method+url, response, r.cacheTTL)
	}

	//create a response object
	return response
}
