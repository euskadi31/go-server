// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package locale

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/justinas/alice"
	"github.com/stretchr/testify/assert"
)

func TestLocaleHandlerWithoutAcceptLanguage(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	middleware := alice.New(Handler()).ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		locale := FromContext(r.Context())

		assert.Equal(t, "fr", locale.Language)
	})

	middleware.ServeHTTP(w, req)
}

func TestLocaleHandlerWithComplexAcceptLanguage(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	middleware := alice.New(Handler()).ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		locale := FromContext(r.Context())

		assert.Equal(t, "fr", locale.Language)
	})

	req.Header.Set("Accept-Language", "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7")

	middleware.ServeHTTP(w, req)
}

func TestLocaleHandlerWithoutSupportedLanguage(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	middleware := alice.New(Handler()).ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		locale := FromContext(r.Context())

		assert.Equal(t, "fr", locale.Language)
	})

	req.Header.Set("Accept-Language", "en-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7")

	middleware.ServeHTTP(w, req)
}

func TestLocaleHandlerWithConfig(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	middleware := alice.New(HandlerWithConfig([]string{"fr", "en", "es", "it"})).ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		locale := FromContext(r.Context())

		assert.Equal(t, "en", locale.Language)
		assert.Equal(t, "US", locale.Region)
	})

	req.Header.Set("Accept-Language", "en-FR,en;q=0.9,en-US;q=0.8,en;q=0.7")

	middleware.ServeHTTP(w, req)
}

func TestLocaleHandlerWithConfiAndAcceptLanguageMixed(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	middleware := alice.New(HandlerWithConfig([]string{"fr", "en", "es", "it"})).ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		locale := FromContext(r.Context())

		assert.Equal(t, "en", locale.Language)

		//@BUG: expected FR, bug with language lib ?
		assert.Equal(t, "US", locale.Region)
	})

	req.Header.Set("Accept-Language", "en-FR")

	middleware.ServeHTTP(w, req)
}

func TestLocaleHandlerWithConfigWithSimpleLocale(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	middleware := alice.New(HandlerWithConfig([]string{"fr", "en", "es", "it"})).ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		locale := FromContext(r.Context())

		assert.Equal(t, "en", locale.Language)
		assert.Equal(t, "US", locale.Region)
	})

	req.Header.Set("Accept-Language", "en")

	middleware.ServeHTTP(w, req)
}

func TestLocaleHandlerWithConfigWithSimpleLocaleNotSupported(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	middleware := alice.New(HandlerWithConfig([]string{"fr", "en", "es", "it"})).ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		locale := FromContext(r.Context())

		assert.Equal(t, "fr", locale.Language)
		assert.Equal(t, "FR", locale.Region)
	})

	req.Header.Set("Accept-Language", "de")

	middleware.ServeHTTP(w, req)
}

func TestFromContext(t *testing.T) {
	ctx := context.Background()

	locale := FromContext(ctx)

	assert.Equal(t, "fr", locale.Language)
	assert.Equal(t, "FR", locale.Region)
}

func BenchmarkHandler(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)

	middleware := alice.New(HandlerWithConfig([]string{"fr", "en", "es", "it"})).ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		FromContext(r.Context())
	})

	for n := 0; n < b.N; n++ {
		w := httptest.NewRecorder()

		req.Header.Set("Accept-Language", "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7")

		middleware.ServeHTTP(w, req)
	}
}
