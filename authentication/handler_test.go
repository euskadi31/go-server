// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package authentication

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/justinas/alice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler(t *testing.T) {
	provider := &MockProvider{}

	provider.On("Validate", mock.Anything).Return(true)

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	middleware := alice.New(Handler(&Configuration{
		Realm: "Test",
	}, provider)).ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandlerWithBadAuth(t *testing.T) {
	provider := &MockProvider{}

	provider.On("Validate", mock.Anything).Return(false)

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	middleware := alice.New(Handler(&Configuration{
		Realm: "Test",
	}, provider)).ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandlerOnHealthEndpoint(t *testing.T) {
	provider := &MockProvider{}

	provider.On("Validate", mock.Anything).Return(false)

	req := httptest.NewRequest(http.MethodGet, "http://example.com/health", nil)
	w := httptest.NewRecorder()

	middleware := alice.New(Handler(&Configuration{
		Realm: "Test",
	}, provider)).ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
