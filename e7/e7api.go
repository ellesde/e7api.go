package e7

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultBaseURL = "https://api.epicsevendb.com/"
)

// A Client manages communication with the EpicSevenDB API.
type Client struct {
	// HTTP client used to communicate with API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// Reuse single struct instead of allocating one for each service on the heap.
	common service

	// Services used for talking to different parts of the EpicSevenDB API.
	Heroes *HeroesService
}

type service struct {
	client *Client
}

// NewClient returns a new EpicSevenDB API client.
func NewClient() *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		client:  &http.Client{},
		BaseURL: baseURL,
	}
	c.common.client = c
	c.Heroes = (*HeroesService)(&c.common)
	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.
func (c *Client) NewRequest(method, urlStr string) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// ErrNilContext is returned when a nil context is provided in Do.
var ErrNilContext error = errors.New("context must be non-nil")

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error occurred.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it
// is canceled or times out, ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}

	defer func() {
		// Ensure the response body is fully read and closed
		// before we reconnect, so that we reuse the same TCP connection.
		// Close the previous response's body but read at least some of
		// the body so if it's small the underlying TCP connection will
		// be re-used. No need to check for errors: if it fails, the
		// Transport won't reuse it anyway.
		const maxBodySlurpSize = 2 << 10
		if resp.ContentLength == -1 || resp.ContentLength <= maxBodySlurpSize {
			io.CopyN(ioutil.Discard, resp.Body, maxBodySlurpSize)
		}

		resp.Body.Close()
	}()

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return resp, err
}

// CheckResponse checks the API response for errors and returns them if
// present. A response is considered an error if it has a status code
// outside the 200 range or equal to 202 accepted.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{
		Response: r,
	}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	// Re-populate error response body
	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	switch {
	default:
		return errorResponse
	}
}

// An ErrorResponse reports one or more errors caused by an API request.
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response
	Message  string   `json:"error,omitempty"`
	Stack    string   `json:"stack,omitempty"`
	Metadata Metadata `json:"meta,omitempty"`
}

// Metadata is set of metadata that is returned in every EpicSevenDB API rsponse.
type Metadata struct {
	RequestDate string `json:"requestDate,omitempty"`
	APIVersion  string `json:"apiVersion,omitempty"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method,
		r.Response.Request.URL,
		r.Response.StatusCode,
		r.Message)
}
