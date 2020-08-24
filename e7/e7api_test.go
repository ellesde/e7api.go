package e7

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const (
	// baseURLPath is a non-empty Client.BaseURL path to use during tests
	// to ensure relative URLs are used for all endpoints.
	baseURLPath = "/api-v2"
)

// setup sets up a test HTTP server along with a e7api.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provides mock responses for the API method being tested.
func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative. It only makes a difference
	// when there's a non-empty base URL path. So, use that.
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the EpicSevenDB client being tested and is configured to use test server.
	client = NewClient()
	url, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = url

	return client, mux, server.URL, server.Close
}

func TestNewClient(t *testing.T) {
	c := NewClient()

	if got, want := c.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}

	c2 := NewClient()
	if c.client == c2.client {
		t.Error("NewClient returned same http.Clients but they should differ")
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient()

	inURL, outURL := "/foo", defaultBaseURL+"foo"
	req, _ := c.NewRequest(http.MethodGet, inURL)

	// test that relative URL was expanded
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}
}

func TestNewRequest_badURL(t *testing.T) {
	c := NewClient()
	_, err := c.NewRequest(http.MethodGet, ":")
	testURLParseError(t, err)
}

func testURLParseError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	var e *url.Error
	if !errors.As(err, &e) || e.Op != "parse" {
		t.Errorf("expected URL parse error, got %+v", err)
	}
}

func TestNewRequest_badMethod(t *testing.T) {
	c := NewClient()
	_, err := c.NewRequest("bad method", ".")
	if err == nil {
		t.Errorf("expected error to be returned")
	}
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("request method: %v, want %v", got, want)
	}
}

func TestNewRequest_errorForNoTrailingSlash(t *testing.T) {
	tests := []struct {
		rawurl    string
		wantError bool
	}{
		{rawurl: "https://example.com/api/v3", wantError: true},
		{rawurl: "https://example.com/api/v3/", wantError: false},
	}
	c := NewClient()
	for _, test := range tests {
		u, err := url.Parse(test.rawurl)
		if err != nil {
			t.Fatalf("url.Parse returned unexpected error: %v.", err)
		}
		c.BaseURL = u
		if _, err := c.NewRequest(http.MethodGet, "test"); test.wantError && err == nil {
			t.Errorf("expected error to be returned.")
		} else if !test.wantError && err != nil {
			t.Fatalf("NewRequest returned unexpected error: %v.", err)
		}
	}
}

func TestDo(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest("GET", ".")
	body := new(foo)
	client.Do(context.Background(), req, body)

	want := &foo{"a"}
	if diff := cmp.Diff(body, want); diff != "" {
		t.Errorf("response body mismatch (-want +got):\n%s", diff)
	}
}

func TestDo_nilContext(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	req, _ := client.NewRequest("GET", ".")
	_, err := client.Do(nil, req, nil)

	if !errors.Is(err, ErrNilContext) {
		t.Errorf("expected context must be non-nil error")
	}
}

func TestDo_httpError(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest("GET", ".")
	resp, err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Fatal("expected HTTP 400 error, got no error.")
	}
	if resp.StatusCode != 400 {
		t.Errorf("expected HTTP 400 error, got %d status code.", resp.StatusCode)
	}
}

// Test handling of an error caused by the internal http client's Do()
// function.
func TestDo_redirectLoop(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, baseURLPath, http.StatusFound)
	})

	req, _ := client.NewRequest("GET", ".")
	_, err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Error("expected error to be returned.")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("expected a URL error; got %#v.", err)
	}
}

func TestDo_noContent(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	var body json.RawMessage

	req, _ := client.NewRequest("GET", ".")
	_, err := client.Do(context.Background(), req, &body)
	if err != nil {
		t.Fatalf("Do returned unexpected error: %v", err)
	}
}

func TestDo_contextTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}

	client, mux, _, teardown := setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	req, _ := client.NewRequest("GET", ".")
	body := new(foo)
	_, err := client.Do(ctx, req, body)
	if !errors.Is(err, context.DeadlineExceeded) || strings.Contains(err.Error(), "canceled") {
		t.Errorf("expected DeadlineExceeded error, got: %#v", err)
	}
}

func TestDo_contextCanceled(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping HTTP request in short mode")
	}

	client, mux, _, teardown := setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
	})

	ctx, cancel := context.WithCancel(context.Background())

	req, _ := client.NewRequest("GET", ".")
	body := new(foo)

	go func() {
		time.Sleep(5 * time.Millisecond)
		cancel()
	}()

	_, err := client.Do(ctx, req, body)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expect a Canceled error, got: %#v", err)
	}
}

func TestCheckResponse(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body: ioutil.NopCloser(strings.NewReader(
			`{
				"error": "Invalid request. Please read the API docs. Open an issue on Github if this keeps happening.",
  				"stack": "",
  				"meta": {
    				"requestDate": "Sat Aug 22 00:54:50 UTC 2020",
    				"apiVersion": "2.1.0"
  				}
			}`)),
	}
	err := CheckResponse(res).(*ErrorResponse)

	if err == nil {
		t.Errorf("expected error response.")
	}

	want := &ErrorResponse{
		Response: res,
		Message:  "Invalid request. Please read the API docs. Open an issue on Github if this keeps happening.",
		Stack:    "",
		Metadata: Metadata{
			RequestDate: "Sat Aug 22 00:54:50 UTC 2020",
			APIVersion:  "2.1.0",
		},
	}

	ignore := cmpopts.IgnoreUnexported(bytes.Buffer{}, http.Request{})
	if diff := cmp.Diff(err, want, ignore); diff != "" {
		t.Errorf("error mismatch (-want +got):\n%s", diff)
	}
}

// Ensure that we properly handle API errors that do not contain a response body
func TestCheckResponse_noBody(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	err := CheckResponse(res).(*ErrorResponse)

	if err == nil {
		t.Errorf("expected error response.")
	}

	want := &ErrorResponse{
		Response: res,
	}

	ignore := cmpopts.IgnoreUnexported(bytes.Buffer{}, http.Request{})
	if diff := cmp.Diff(err, want, ignore); diff != "" {
		t.Errorf("error mismatch (-want +got):\n%s", diff)
	}
}

func TestCheckResponse_unexpectedErrorStructure(t *testing.T) {
	httpBody := `{"error":"m", "errors": ["error 1"]}`
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader(httpBody)),
	}
	err := CheckResponse(res).(*ErrorResponse)

	if err == nil {
		t.Errorf("expected error response.")
	}

	want := &ErrorResponse{
		Response: res,
		Message:  "m",
	}

	ignore := cmpopts.IgnoreUnexported(bytes.Buffer{}, http.Request{})
	if diff := cmp.Diff(err, want, ignore); diff != "" {
		t.Errorf("error mismatch (-want +got):\n%s", diff)
	}

	data, err2 := ioutil.ReadAll(err.Response.Body)
	if err2 != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	got := string(data)
	if diff := cmp.Diff(got, httpBody); diff != "" {
		t.Errorf("ErrorResponse.Response.Body mismatch (-want +got):\n%s", diff)
	}
}

func TestErrorResponse_Error(t *testing.T) {
	res := &http.Response{Request: &http.Request{}}
	err := ErrorResponse{Message: "m", Response: res}
	if err.Error() == "" {
		t.Errorf("expected non-empty ErrorResponse.Error()")
	}
}
