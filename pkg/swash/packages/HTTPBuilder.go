package packages

import (
	"io"
	"net/http"
	"time"
)

type HTTPBuilder struct {
	c *http.Client
	r *http.Request

	Ok bool
}

// Make will produce a brand new
func Make(method, url string) *HTTPBuilder {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {

		/* error handles by settings Ok to false */
		return &HTTPBuilder{
			Ok: false,
		}
	}

	return &HTTPBuilder{
		r: req,
		c: &http.Client{
			/* default license timeout */
			Timeout: 1 * time.Second,
		},
	}
}

// SetTimeout will modify the clients timeout period on making http requests
func (httpBuilder *HTTPBuilder) SetTimeout(timeout int) {
	httpBuilder.c.Timeout = time.Duration(timeout)
}

// AddHeader will append a header into the request handler
func (httpBuilder *HTTPBuilder) AddHeader(headerName, headerValue string) {
	httpBuilder.r.Header.Add(headerName, headerValue)
}

// Do will perform the request and return the response
func (httpBuilder *HTTPBuilder) Do() string {
	clone, err := httpBuilder.c.Do(httpBuilder.r)
	if err != nil {
		return err.Error()
	}

	content, err := io.ReadAll(clone.Body)
	if err != nil {
		return err.Error()
	}

	return string(content)
}

// DoWithJson shall make the request and then parse the response into json form
func (httpBuilder *HTTPBuilder) DoWithJson() any {
	content := httpBuilder.Do()
	if len(content) == 0 {
		return make(map[string]any)
	}

	return JsonDecode(content)
}
