// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package encoder

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	Register(JSONEncoder())

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	req.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()

	Encode(w, req, http.StatusOK, true)

	body, err := ioutil.ReadAll(w.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "true\n", string(body))

	encoders = make(map[string]Encoder)
}

func TestEncodeWithBadMimeType(t *testing.T) {
	Register(JSONEncoder())

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	req.Header.Set("Accept", "application/xml")

	w := httptest.NewRecorder()

	Encode(w, req, http.StatusOK, true)

	body, err := ioutil.ReadAll(w.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "true\n", string(body))

	encoders = make(map[string]Encoder)
}

func TestEncodeWithError(t *testing.T) {
	provider := &MockEncoder{}

	provider.On("MimeType").Return("application/json")
	provider.On("Encode", mock.Anything, true).Return(errors.New("bad"))

	Register(provider)

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	req.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()

	Encode(w, req, http.StatusOK, true)

	body, err := ioutil.ReadAll(w.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "", string(body))

	encoders = make(map[string]Encoder)
}
