package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Response represents an HTTP response along with additional metadata.
type Response struct {
	*http.Response // Embedded http.Response for composition
	hasRead       bool
}

// Request represents an HTTP request with customizable parameters.
type Request struct {
	url     string        // URL of the request
	headers [][2]string   // Headers of the request
	method  string        // HTTP method of the request
	body    []byte        // Body of the request
	err     error         // Error associated with the request
}

// Fetch returns a new GET Request for the given URL.
func Fetch(url string) *Request {
	return NewRequest("GET", url)
}

// NewRequest creates a new Request with the specified method and URL.
func NewRequest(method, url string) *Request {
	return &Request{url: url, method: strings.ToUpper(method), headers: [][2]string{}, body: []byte{}, err: nil}
}

// Header sets a single header key-value pair for the request.
func (r *Request) Header(key, value string) *Request {
	r.headers = append(r.headers, [2]string{key, value})
	return r
}

// Headers sets multiple header key-value pairs for the request.
func (r *Request) Headers(headers [][2]string) *Request {
	r.headers = append(r.headers, headers...)
	return r
}

// Method sets the HTTP method for the request.
func (r *Request) Method(method string) *Request {
	r.method = strings.ToUpper(method)
	return r
}

// Body sets the body content for the request.
func (r *Request) Body(bytes []byte) *Request {
	r.body = bytes
	return r
}

// JSON sets the request body to a JSON representation of the given content.
func (r *Request) JSON(content any) *Request {
	r.Header("Content-Type", "application/json")

	bytes, err := json.Marshal(content)
	if err != nil {
		r.err = err
	} else {
		r.Body(bytes)
	}

	return r
}

// Execute sends the HTTP request and returns the response.
func (request *Request) Execute() (*Response, error) {
	if request.err != nil {
		return nil, request.err
	}

	client := http.Client{}

	method := request.method
	if method == "" {
		method = "GET"
	}
	httpRequest, err := http.NewRequest(method, request.url, bytes.NewBuffer(request.body))
	if err != nil {
		return nil, err
	}

	for _, header := range request.headers {
		httpRequest.Header.Set(header[0], header[1])
	}

	var response *http.Response
	response, err = client.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	return &Response{Response: response}, nil
}

// IntoBytes reads the response body into bytes.
func (response *Response) IntoBytes() ([]byte, error) {
	if response.hasRead {
		return nil, fmt.Errorf("cannot read body multiple times")
	}
	defer response.Body.Close()
	bytes, err := io.ReadAll(response.Body)
	response.hasRead = true
	return bytes, err
}

// IntoString reads the response body into a string.
func (response *Response) IntoString() (string, error) {
	bytes, err := response.IntoBytes()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJSON reads the response body and decodes it into the provided structure.
func (response *Response) FromJSON(item any) error {
	bytes, err := response.IntoBytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, item)
}
