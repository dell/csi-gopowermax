/*
 Copyright © 2025 Dell Inc. or its subsidiaries. All Rights Reserved.

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
	"fmt"
	"io"
	"net/http"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type ErrorReader struct{}

func (r *ErrorReader) Close() error {
	return fmt.Errorf("error closing the body")
}

func (r *ErrorReader) Read(_ []byte) (n int, err error) {
	return 0, fmt.Errorf("error reading the body")
}

func TestIsBinOctetBody(t *testing.T) {
	tests := []struct {
		name     string
		header   http.Header
		expected bool
	}{
		{"BinaryOctetStream", http.Header{HeaderKeyContentType: []string{headerValContentTypeBinaryOctetStream}}, true},
		{"NonBinaryOctetStream", http.Header{HeaderKeyContentType: []string{"text/plain"}}, false},
		{"EmptyHeader", http.Header{}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isBinOctetBody(test.header)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestLogRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	assert.NoError(t, err)

	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetLevel(log.DebugLevel)

	logRequest(context.Background(), req, func(lf func(args ...interface{}), msg string) {
		lf(msg)
	})

	assert.Contains(t, buf.String(), "POWERMAX HTTP REQUEST")
}

func TestLogResponse(t *testing.T) {
	res := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBufferString(`{"key":"value"}`)),
	}

	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetLevel(log.DebugLevel)

	logResponse(context.Background(), res, func(lf func(args ...interface{}), msg string) {
		lf(msg)
	})

	assert.Contains(t, buf.String(), "POWERMAX HTTP RESPONSE")
}

func TestWriteIndentedN(t *testing.T) {
	var buf bytes.Buffer
	err := WriteIndentedN(&buf, []byte("line1\nline2"), 2)
	assert.NoError(t, err)
	assert.Equal(t, "  line1\n  line2", buf.String())
}

func TestWriteIndented(t *testing.T) {
	var buf bytes.Buffer
	err := WriteIndented(&buf, []byte("line1\nline2"))
	assert.NoError(t, err)
	assert.Equal(t, "    line1\n    line2", buf.String())
}

func TestDrainBody(t *testing.T) {
	body := io.NopCloser(bytes.NewBufferString("test body"))
	r1, r2, err := drainBody(body)
	assert.NoError(t, err)

	buf1 := new(bytes.Buffer)
	buf1.ReadFrom(r1)
	assert.Equal(t, "test body", buf1.String())

	buf2 := new(bytes.Buffer)
	buf2.ReadFrom(r2)
	assert.Equal(t, "test body", buf2.String())
}

func TestDumpRequest(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		url         string
		body        io.Reader
		headers     map[string]string
		expectError bool
		expected    []string
	}{
		{
			name:     "GET request without body",
			method:   "GET",
			url:      "http://example.com",
			body:     bytes.NewBufferString(""),
			headers:  map[string]string{},
			expected: []string{"GET / HTTP/1.1", "Host: example.com"},
		},
		{
			name:     "POST request with body",
			method:   "POST",
			url:      "http://example.com",
			body:     bytes.NewBufferString("test body"),
			headers:  map[string]string{"Content-Type": "application/json"},
			expected: []string{"POST / HTTP/1.1", "Host: example.com", "Content-Type: application/json", "test body"},
		},
		{
			name:     "Request with Authorization header",
			method:   "GET",
			url:      "http://example.com",
			body:     bytes.NewBufferString(""),
			headers:  map[string]string{"Authorization": "Basic dXNlcjpwYXNz"},
			expected: []string{"GET / HTTP/1.1", "Host: example.com"},
		},
		{
			name:   "Request with invalid Authorization header",
			method: http.MethodGet,
			url:    "http://example.com",
			body:   bytes.NewBufferString(""),
			headers: map[string]string{
				"Authorization": "Basic invalid_base64_string",
			},
			expectError: true,
			expected:    []string{"GET / HTTP/1.1", "Host: example.com"},
		},
		{
			name:     "Request with empty Host",
			method:   http.MethodGet,
			body:     bytes.NewBufferString(""),
			headers:  map[string]string{"Authorization": "Basic dXNlcjpwYXNz"},
			expected: []string{"GET / HTTP/1.1", "Host: example.com"},
		},
		{
			name:        "Request with invalid body",
			method:      http.MethodGet,
			body:        &ErrorReader{},
			headers:     map[string]string{"Authorization": "Basic dXNlcjpwYXNz"},
			expectError: true,
			expected:    []string{"GET / HTTP/1.1", "Host: example.com"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.method, test.url, test.body)
			req.TransferEncoding = []string{"chunked"}
			req.Close = true
			assert.NoError(t, err)

			for key, value := range test.headers {
				req.Header.Set(key, value)
			}
			req.URL.Host = "example.com"

			var buf bytes.Buffer
			log.SetOutput(&buf)
			log.SetLevel(log.DebugLevel)

			dump, err := dumpRequest(req, true)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				for _, expected := range test.expected {
					assert.Contains(t, string(dump), expected)
				}
				if _, isAuth := test.headers["Authorization"]; isAuth {
					assert.Contains(t, buf.String(), "username: user , password: *****")
				}
			}
		})
	}
}
