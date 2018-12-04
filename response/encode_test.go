// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package response

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/euskadi31/go-server/response/encoder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEncode(t *testing.T) {
	Register(encoder.JSON())

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	req.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()

	Encode(w, req, http.StatusOK, true)

	body, err := ioutil.ReadAll(w.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "true\n", string(body))
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	encoders = make(map[string]encoder.Encoder)
}

func TestEncodeWithBadMimeType(t *testing.T) {
	Register(encoder.JSON())

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	req.Header.Set("Accept", "application/xml")

	w := httptest.NewRecorder()

	Encode(w, req, http.StatusOK, true)

	body, err := ioutil.ReadAll(w.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "true\n", string(body))

	encoders = make(map[string]encoder.Encoder)
}

func TestEncodeWithError(t *testing.T) {
	provider := &encoder.MockEncoder{}

	provider.On("MimeType").Return("application/json")
	provider.On("Encode", mock.Anything, true).Return(errors.New("bad"))

	Register(provider)

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	req.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()

	Encode(w, req, http.StatusOK, true)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error":{"code":500,"message":"bad"}}`, w.Body.String())
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	encoders = make(map[string]encoder.Encoder)
}

func TestEncodeWithWriteError(t *testing.T) {
	Register(encoder.JSON())

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	req.Header.Set("Accept", "application/json")

	w := &mockResponseWriter{}

	w.On("Header").Return(http.Header{})
	w.On("WriteHeader", http.StatusOK).Return()
	w.On("Write", mock.Anything).Return(0, errors.New("fail"))

	Encode(w, req, http.StatusOK, true)

	encoders = make(map[string]encoder.Encoder)

	w.AssertExpectations(t)
}

func BenchmarkEncode(b *testing.B) {
	Register(encoder.JSON())

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	req.Header.Set("Accept", "application/xml")

	for n := 0; n < b.N; n++ {
		w := httptest.NewRecorder()

		Encode(w, req, http.StatusOK, true)
	}
}
