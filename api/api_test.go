/*
 Copyright © 2021-2025 Dell Inc. or its subsidiaries. All Rights Reserved.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
      http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	types "github.com/dell/gopowermax/v2/types/v100"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type stubTypeWithMetaData struct{}

type MyReader struct {
	reader io.Reader
}

func (r *MyReader) Read(p []byte) (int, error) {
	return r.reader.Read(p)
}

func (r *MyReader) Close() error {
	return nil
}

func (s stubTypeWithMetaData) MetaData() http.Header {
	h := make(http.Header)
	h.Set("foo", "bar")
	return h
}

func Test_addMetaData(t *testing.T) {
	tests := []struct {
		name           string
		givenHeader    map[string]string
		expectedHeader map[string]string
		body           interface{}
	}{
		{"nil header is a noop", nil, nil, nil},
		{"nil body is a noop", nil, nil, nil},
		{"header is updated", make(map[string]string), map[string]string{"Foo": "bar"}, stubTypeWithMetaData{}},
		{"header is not updated", make(map[string]string), map[string]string{}, struct{}{}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			addMetaData(tt.givenHeader, tt.body)

			switch {
			case tt.givenHeader == nil:
				if tt.givenHeader != nil {
					t.Errorf("(%s): expected %s, actual %s", tt.body, tt.expectedHeader, tt.givenHeader)
				}
			case tt.body == nil:
				if len(tt.givenHeader) != 0 {
					t.Errorf("(%s): expected %s, actual %s", tt.body, tt.expectedHeader, tt.givenHeader)
				}
			default:
				if !reflect.DeepEqual(tt.expectedHeader, tt.givenHeader) {
					t.Errorf("(%s): expected %s, actual %s", tt.body, tt.expectedHeader, tt.givenHeader)
				}
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		host        string
		opts        ClientOptions
		debug       bool
		expectError bool
	}{
		{
			name:        "Valid host without options",
			host:        "http://example.com",
			opts:        ClientOptions{},
			debug:       false,
			expectError: false,
		},
		{
			name:        "Empty host",
			host:        "",
			opts:        ClientOptions{},
			debug:       false,
			expectError: true,
		},
		{
			name: "Valid host with timeout",
			host: "http://example.com",
			opts: ClientOptions{
				Timeout: 10 * time.Second,
			},
			debug:       false,
			expectError: false,
		},
		{
			name: "Valid host with insecure option",
			host: "http://example.com",
			opts: ClientOptions{
				Insecure: true,
			},
			debug:       false,
			expectError: false,
		},
		{
			name: "Valid host with dummy cert file",
			host: "http://example.com",
			opts: ClientOptions{
				CertFile: "../mock/cert.pem",
			},
			debug:       false,
			expectError: true,
		},
		{
			name: "Host with showHTTP option",
			host: "http://example.com",
			opts: ClientOptions{
				ShowHTTP: true,
			},
			debug:       false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New(tt.host, tt.opts, tt.debug)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected nil error, but got: %v", err)
				}
				if c == nil {
					t.Error("Expected non-nil client, but got nil")
				}
			}
		})
	}
}

func (m *MockClient) GetHTTPClient() *http.Client {
	return m.http
}

func (m *MockClient) DoWithHeaders(
	ctx context.Context,
	method, path string,
	headers map[string]string,
	body, resp interface{},
) error {
	args := m.Called(ctx, method, path, headers, body, resp)
	return args.Error(0)
}

func TestGetHTTPClient(t *testing.T) {
	mockClient := &MockClient{http: &http.Client{}}
	client := &client{http: mockClient.GetHTTPClient()}

	assert.Equal(t, mockClient.GetHTTPClient(), client.GetHTTPClient())
}

func TestBeginsWithSlash(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"/path", true},
		{"path/", false},
		{"/", true},
	}

	for _, test := range tests {
		result := beginsWithSlash(test.input)
		if result != test.expected {
			t.Errorf("beginsWithSlash(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestEndsWithSlash(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"/path/", true},
		{"path/", true},
		{"/", true},
		{"path", false},
	}

	for _, test := range tests {
		result := endsWithSlash(test.input)
		if result != test.expected {
			t.Errorf("endsWithSlash(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}

type MockClient struct {
	*client
	http *http.Client
	mock.Mock
}

func (m *MockClient) DoAndGetResponseBody(
	ctx context.Context,
	method, uri string,
	headers map[string]string,
	body interface{},
) (*http.Response, error) {
	args := m.Called(ctx, method, uri, headers, body)
	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}
	return resp.(*http.Response), args.Error(1)
}

func TestDoWithHeaders(t *testing.T) {
	mockHTTPClient := new(MockHTTPClient)
	httpClient := &http.Client{
		Transport: &MockTransport{mockHTTPClient: mockHTTPClient},
	}
	mockClient := &client{
		http:     httpClient,
		host:     "https://example.com",
		token:    "mockToken",
		showHTTP: false,
	}

	tests := []struct {
		name          string
		method        string
		uri           string
		headers       map[string]string
		body          interface{}
		resp          interface{}
		mockResponse  *http.Response
		mockDoError   error
		expectedError string
	}{
		{
			name:   "Successful GET request with response",
			method: http.MethodGet,
			uri:    "/test",
			headers: map[string]string{
				"Custom-Header": "value",
			},
			body: nil,
			resp: &map[string]interface{}{},
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"success": true}`)),
			},
			mockDoError:   nil,
			expectedError: "",
		},
		{
			name:          "Nil response",
			method:        http.MethodGet,
			uri:           "/test",
			headers:       nil,
			body:          nil,
			resp:          nil,
			mockResponse:  nil,
			mockDoError:   nil,
			expectedError: "",
		},
		{
			name:   "Failed to decode response body",
			method: http.MethodGet,
			uri:    "/test",
			headers: map[string]string{
				"Custom-Header": "value",
			},
			body: nil,
			resp: &map[string]interface{}{},
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"invalid_json":`)),
			},
			mockDoError:   nil,
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHTTPClient.On("Do", mock.Anything).Return(tt.mockResponse, tt.mockDoError)

			err := mockClient.DoWithHeaders(
				context.TODO(),
				tt.method,
				tt.uri,
				tt.headers,
				tt.body,
				tt.resp,
			)

			if tt.expectedError != "" {
				if err == nil || err.Error() != tt.expectedError {
					t.Errorf("Expected error: %v, but got: %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error, but got: %v", err)
				}
			}

			mockHTTPClient.AssertExpectations(t)
		})
	}
}

// MockHTTPClient is a mock http client
type MockHTTPClient struct {
	mock.Mock
	http.Client
}

// Do sends an HTTP request to the API
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}
	return resp.(*http.Response), args.Error(1)
}

// MockTransport is a mock http transport
type MockTransport struct {
	mockHTTPClient *MockHTTPClient
	Response       *http.Response
}

// RoundTrip sends an HTTP request to the API
func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.mockHTTPClient.Do(req)
}

func TestDoAndGetResponseBody(t *testing.T) {
	mockHTTPClient := new(MockHTTPClient)
	httpClient := &http.Client{
		Transport: &MockTransport{mockHTTPClient: mockHTTPClient},
	}

	tests := []struct {
		name          string
		method        string
		uri           string
		headers       map[string]string
		body          interface{}
		mockResponse  *http.Response
		mockError     error
		expectedError string
		c             *client
	}{
		{
			name:   "Successful GET request",
			method: http.MethodGet,
			uri:    "/test",
			headers: map[string]string{
				"Custom-Header": "value",
			},
			body: nil,
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"success": true}`)),
			},
			c: &client{
				http:     httpClient,
				host:     "https://example.com",
				token:    "mockToken",
				showHTTP: false,
			},
			mockError:     nil,
			expectedError: "",
		},
		{
			name:   "Failed POST request with invalid JSON body",
			method: http.MethodPost,
			uri:    "/test",
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			body:         make(chan int), // invalid JSON body
			mockResponse: nil,
			c: &client{
				http:     httpClient,
				host:     "https://example.com",
				token:    "mockToken",
				showHTTP: false,
			},
			mockError:     errors.New("unsupported type error"),
			expectedError: "json: unsupported type: chan int",
		},
		{
			name:   "Handle io.ReadCloser body",
			method: http.MethodPost,
			uri:    "/test",
			headers: map[string]string{
				"Content-Type": "application/octet-stream",
			},
			body: io.NopCloser(bytes.NewBufferString(`binary content`)),
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"success": true}`)),
			},
			c: &client{
				http:     httpClient,
				host:     "https://example.com",
				token:    "mockToken",
				showHTTP: false,
			},
			mockError:     nil,
			expectedError: "",
		},
		{
			name:   "POST request with JSON body and Content-Type header set",
			method: http.MethodPost,
			uri:    "/test",
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			body: map[string]string{"key": "value"},
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"success": true}`)),
			},
			c: &client{
				http:     httpClient,
				host:     "https://example.com",
				token:    "mockToken",
				showHTTP: false,
			},
			mockError:     nil,
			expectedError: "",
		},
		{
			name:    "POST request with JSON body without Content-Type header set",
			method:  http.MethodPost,
			uri:     "/test",
			headers: map[string]string{"Custom-Header": "application/json"},
			body: &MyReader{
				reader: strings.NewReader("Success"),
			},
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"success": true}`)),
			},
			c: &client{
				http:     httpClient,
				host:     "https://example.com",
				token:    "mockToken",
				showHTTP: false,
			},
			mockError:     nil,
			expectedError: "",
		},
		{
			name:   "Get request with path not starting with /",
			method: http.MethodGet,
			uri:    "test",
			headers: map[string]string{
				"content-type": "application/json",
			},
			body: nil,
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"success": true}`)),
			},
			c: &client{
				http:     httpClient,
				host:     "https://example.com",
				token:    "mockToken",
				showHTTP: false,
			},
			mockError:     nil,
			expectedError: "",
		},
		{
			name:   "Successful GET request having client with showHTTP set to true",
			method: http.MethodGet,
			uri:    "test",
			headers: map[string]string{
				"content-type": "application/json",
			},
			body: nil,
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"success": true}`)),
			},
			c: &client{
				http:     httpClient,
				host:     "https://example.com",
				token:    "mockToken",
				showHTTP: true,
			},
			mockError:     nil,
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHTTPClient.On("Do", mock.Anything).Return(tt.mockResponse, tt.mockError)

			res, err := tt.c.DoAndGetResponseBody(
				context.Background(),
				tt.method,
				tt.uri,
				tt.headers,
				tt.body,
			)

			if tt.expectedError != "" {
				if err == nil || err.Error() != tt.expectedError {
					t.Errorf("Expected error: %v, but got: %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error, but got: %v", err)
				}
				if !reflect.DeepEqual(res, tt.mockResponse) {
					t.Errorf("Expected response: %v, but got: %v", tt.mockResponse, res)
				}
			}

			mockHTTPClient.AssertExpectations(t)
		})
	}
}

func TestGetToken(t *testing.T) {
	// Test case: token is not set
	c := &client{}
	if c.GetToken() != "" {
		t.Errorf("GetToken() = %v, want %v", c.GetToken(), "")
	}

	// Test case: token is set
	c.token = "testToken"
	if c.GetToken() != "testToken" {
		t.Errorf("GetToken() = %v, want %v", c.GetToken(), "testToken")
	}
}

func TestSetToken(t *testing.T) {
	c := &client{}

	c.SetToken("token1")
	if c.token != "token1" {
		t.Errorf("Expected token to be 'token1', got '%s'", c.token)
	}

	c.SetToken("token2")
	if c.token != "token2" {
		t.Errorf("Expected token to be 'token2', got '%s'", c.token)
	}
}

func TestParseJSONError(t *testing.T) {
	tests := []struct {
		name          string
		responseBody  string
		statusCode    int
		expectedError *types.Error
	}{
		{
			name:         "Valid JSON error response",
			responseBody: `{"Message": "error occurred"}`,
			statusCode:   http.StatusBadRequest,
			expectedError: &types.Error{
				HTTPStatusCode: http.StatusBadRequest,
				Message:        "error occurred",
			},
		},
		{
			name:         "Invalid JSON error response",
			responseBody: `invalid json`,
			statusCode:   http.StatusInternalServerError,
			expectedError: &types.Error{
				HTTPStatusCode: http.StatusInternalServerError,
				Message:        http.StatusText(http.StatusInternalServerError),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &http.Response{
				StatusCode: tt.statusCode,
				Body:       io.NopCloser(bytes.NewBufferString(tt.responseBody)),
			}

			c := &client{}
			err := c.ParseJSONError(r)

			assert.Error(t, err)
			if e, ok := err.(*types.Error); ok {
				assert.Equal(t, tt.expectedError.HTTPStatusCode, e.HTTPStatusCode)
				assert.Equal(t, tt.expectedError.Message, e.Message)
			}
		})
	}
}

func TestDoLog(t *testing.T) {
	type fields struct {
		debug bool
	}
	type args struct {
		l   func(args ...interface{})
		msg string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "debug is true",
			fields: fields{
				debug: true,
			},
			args: args{
				l:   log.Println,
				msg: "test message",
			},
		},
		{
			name: "debug is false",
			fields: fields{
				debug: false,
			},
			args: args{
				l:   log.Println,
				msg: "test message",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			c := &client{
				debug: tt.fields.debug,
			}
			c.doLog(tt.args.l, tt.args.msg)
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		headers      map[string]string
		resp         interface{}
		expectedErr  error
		expectedBody string
	}{
		{
			name: "Successful Get Request",
			path: "/api/test",
			headers: map[string]string{
				"content-type": "application/json",
			},
			resp:         nil,
			expectedErr:  nil,
			expectedBody: `{"message":"Success"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				for key, value := range tt.headers {
					w.Header().Add(key, value)
				}

				if tt.expectedErr != nil {
					w.WriteHeader(http.StatusBadRequest)
					errData, _ := json.Marshal(tt.expectedErr)
					_, err := w.Write(errData)
					if err != nil {
						return
					}
				} else {
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte(tt.expectedBody))
					if err != nil {
						return
					}
				}
			}))
			defer ts.Close()

			c, err := New(ts.URL, ClientOptions{Timeout: 10 * time.Second}, true)
			if err != nil {
				t.Fatal(err)
			}

			err = c.Get(context.Background(), tt.path, tt.headers, &tt.resp)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			body, err := json.Marshal(tt.resp)
			if err != nil {
				t.Fatal(err)
			}
			if string(body) != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, string(body))
			}
		})
	}
}

func TestPost(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		headers      map[string]string
		resp         interface{}
		body         interface{}
		expectedErr  error
		expectedBody string
	}{
		{
			name: "Successful Post Request",
			path: "/api/test",
			headers: map[string]string{
				"content-type": "application/json",
			},
			resp:         nil,
			expectedErr:  nil,
			expectedBody: `{"message":"Success"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				for key, value := range tt.headers {
					w.Header().Add(key, value)
				}

				if tt.expectedErr != nil {
					w.WriteHeader(http.StatusBadRequest)
					errData, _ := json.Marshal(tt.expectedErr)
					_, err := w.Write(errData)
					if err != nil {
						t.Fatalf("error writing error response: %v", err)
					}
				} else {
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte(tt.expectedBody))
					if err != nil {
						t.Fatalf("error writing response: %v", err)
					}
				}
			}))
			defer ts.Close()

			c, err := New(ts.URL, ClientOptions{Timeout: 10 * time.Second}, true)
			if err != nil {
				t.Fatal(err)
			}

			err = c.Post(context.Background(), tt.path, tt.headers, tt.body, &tt.resp)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			body, err := json.Marshal(tt.resp)
			if err != nil {
				t.Fatal(err)
			}
			if string(body) != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, string(body))
			}
		})
	}
}

func TestPut(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		headers      map[string]string
		resp         interface{}
		body         interface{}
		expectedErr  error
		expectedBody string
	}{
		{
			name: "Successful Put Request",
			path: "/api/test",
			headers: map[string]string{
				"content-type": "application/json",
			},
			resp:         nil,
			body:         nil,
			expectedErr:  nil,
			expectedBody: `{"message":"Success"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				for key, value := range tt.headers {
					w.Header().Add(key, value)
				}

				if tt.expectedErr != nil {
					w.WriteHeader(http.StatusBadRequest)
					errData, _ := json.Marshal(tt.expectedErr)
					_, err := w.Write(errData)
					if err != nil {
						t.Fatalf("error writing error response: %v", err)
					}
				} else {
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte(tt.expectedBody))
					if err != nil {
						t.Fatalf("error writing response: %v", err)
					}
				}
			}))
			defer ts.Close()

			c, err := New(ts.URL, ClientOptions{Timeout: 10 * time.Second}, true)
			if err != nil {
				t.Fatal(err)
			}

			err = c.Put(context.Background(), tt.path, tt.headers, tt.body, &tt.resp)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			body, err := json.Marshal(tt.resp)
			if err != nil {
				t.Fatal(err)
			}
			if string(body) != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, string(body))
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		headers      map[string]string
		resp         interface{}
		expectedErr  error
		expectedBody string
	}{
		{
			name: "Successful Get Request",
			path: "/api/test",
			headers: map[string]string{
				"content-type": "application/json",
			},
			resp:         nil,
			expectedErr:  nil,
			expectedBody: `{"message":"Success"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				for key, value := range tt.headers {
					w.Header().Add(key, value)
				}

				if tt.expectedErr != nil {
					w.WriteHeader(http.StatusBadRequest)
					errData, _ := json.Marshal(tt.expectedErr)
					_, err := w.Write(errData)
					if err != nil {
						return
					}
				} else {
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte(tt.expectedBody))
					if err != nil {
						return
					}
				}
			}))
			defer ts.Close()

			c, err := New(ts.URL, ClientOptions{Timeout: 10 * time.Second}, true)
			if err != nil {
				t.Fatal(err)
			}

			err = c.Delete(context.Background(), tt.path, tt.headers, &tt.resp)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			body, err := json.Marshal(tt.resp)
			if err != nil {
				t.Fatal(err)
			}
			if string(body) != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, string(body))
			}
		})
	}
}

func TestDoMethod(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		headers      map[string]string
		resp         interface{}
		expectedErr  error
		expectedBody string
	}{
		{
			name: "Successful Get Request",
			path: "/api/test",
			headers: map[string]string{
				"content-type": "application/json",
			},
			resp:         nil,
			expectedErr:  nil,
			expectedBody: `{"message":"Success"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				for key, value := range tt.headers {
					w.Header().Add(key, value)
				}

				if tt.expectedErr != nil {
					w.WriteHeader(http.StatusBadRequest)
					errData, _ := json.Marshal(tt.expectedErr)
					_, err := w.Write(errData)
					if err != nil {
						return
					}
				} else {
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte(tt.expectedBody))
					if err != nil {
						return
					}
				}
			}))
			defer ts.Close()

			c, err := New(ts.URL, ClientOptions{Timeout: 10 * time.Second}, true)
			if err != nil {
				t.Fatal(err)
			}

			err = c.Do(context.Background(), http.MethodGet, tt.path, tt.headers, &tt.resp)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			body, err := json.Marshal(tt.resp)
			if err != nil {
				t.Fatal(err)
			}
			if string(body) != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, string(body))
			}
		})
	}
}
