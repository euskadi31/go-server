// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	oaerrors "github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStatusCode(t *testing.T) {
	statusWant := http.StatusForbidden
	var statusGot, customStatusGot int

	type CustomResponseWriter struct {
		http.ResponseWriter
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, http.StatusText(statusWant), statusWant)
		statusGot = StatusCode(w)
		customStatusGot = StatusCode(CustomResponseWriter{w})
	}))
	defer ts.Close()

	if _, err := http.Get(ts.URL); err != nil {
		t.Fatal(err)
	}

	if statusWant != statusGot {
		t.Errorf("http.ResponseWriter: want %d, got %d", statusWant, statusGot)
	}

	if customStatusGot != statusGot {
		t.Errorf("CustomResponseWriter: want %d, got %d", customStatusGot, statusGot)
	}
}

func TestFailure(t *testing.T) {
	w := httptest.NewRecorder()

	Failure(w, http.StatusInternalServerError, ErrorMessage{
		Code:    1337,
		Message: "user_message",
	})

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.JSONEq(t, `{"error":{"code":1337,"message":"user_message"}}`, w.Body.String())
}

func TestFailureFailed(t *testing.T) {
	w := &mockResponseWriter{}

	w.On("Header").Return(http.Header{})
	w.On("WriteHeader", http.StatusInternalServerError).Return()
	w.On("Write", mock.Anything).Return(0, errors.New("fail"))

	Failure(w, http.StatusInternalServerError, ErrorMessage{
		Code:    1337,
		Message: "user_message",
	})
}

func TestFailureFromError(t *testing.T) {
	w := httptest.NewRecorder()

	FailureFromError(w, http.StatusInternalServerError, errors.New("dev_message"))

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.JSONEq(t, `{"error":{"code":500,"message":"dev_message"}}`, w.Body.String())
}

func TestNotFoundFailure(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/users", nil)
	w := httptest.NewRecorder()

	NotFoundFailure(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.JSONEq(t, `{"error":{"code":404,"message":"No route found for \"GET /users\""}}`, w.Body.String())
}

func TestMethodNotAllowedFailure(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/users", nil)
	w := httptest.NewRecorder()

	MethodNotAllowedFailure(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.JSONEq(t, `{"error":{"code":405,"message":"Method \"GET\" not allowed for \"/users\""}}`, w.Body.String())
}

func TestInternalServerFailure(t *testing.T) {
	w := httptest.NewRecorder()

	InternalServerFailure(w)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.JSONEq(t, `{"error":{"code":500,"message":"Internal Server Error"}}`, w.Body.String())
}

func TestServiceUnavailableFailure(t *testing.T) {
	w := httptest.NewRecorder()

	ServiceUnavailableFailure(w, 2*time.Minute)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.HeaderMap.Get("Content-Type"))
	assert.Equal(t, "nosniff", w.HeaderMap.Get("X-Content-Type-Options"))
	assert.Equal(t, "120", w.HeaderMap.Get("Retry-After"))
}

func TestFailureFromValidator(t *testing.T) {
	result := &validate.Result{
		Errors: []error{
			oaerrors.InvalidContentType("application/xml", []string{"application/json"}),
			errors.New("foo"),
		},
	}

	w := httptest.NewRecorder()

	FailureFromValidator(w, result)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.JSONEq(t, `{
		"errors":[
			{"name":"Content-Type","in":"header","value":"application/xml","message":"unsupported media type \"application/xml\", only [application/json] are allowed","values":["application/json"],"code":415},
			{"message":"foo"}
		]
	}`, w.Body.String())
}

func TestFailureFromValidatorFailed(t *testing.T) {
	result := &validate.Result{
		Errors: []error{
			oaerrors.InvalidContentType("application/xml", []string{"application/json"}),
		},
	}

	w := &mockResponseWriter{}

	w.On("Header").Return(http.Header{})
	w.On("WriteHeader", http.StatusBadRequest).Return()
	w.On("Write", mock.Anything).Return(0, errors.New("fail"))

	FailureFromValidator(w, result)
}
