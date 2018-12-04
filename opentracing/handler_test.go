// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package opentracing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	tracer := mocktracer.New()

	middleware := alice.New(Handler(tracer, IgnoreNone)).ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		span := opentracing.SpanFromContext(r.Context())
		assert.NotNil(t, span)
	})

	middleware.ServeHTTP(w, req)
}

func TestHandlerOnRouter(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo/123", nil)

	tracer := mocktracer.New()

	r := mux.NewRouter()
	r.Use(Handler(tracer, IgnoreNone))
	r.HandleFunc("/foo/{id}", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET /foo/{id}", routeRequestName(r))
	}).Methods(http.MethodGet)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
}

func TestHandlerWithIgnore(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	tracer := mocktracer.New()

	middleware := alice.New(Handler(tracer, func(*http.Request) bool {
		return true
	})).ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		span := opentracing.SpanFromContext(r.Context())
		assert.Nil(t, span)
	})

	middleware.ServeHTTP(w, req)
}

func TestIgnoreNone(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)

	assert.False(t, IgnoreNone(req))
}

func TestMassageTemplate(t *testing.T) {
	out := massageTemplate("/articles/{category}/{id:[0-9]+}")

	assert.Equal(t, "/articles/{category}/{id}", out)

	out = massageTemplate("/articles")

	assert.Equal(t, "/articles", out)
}

func TestServerRequestName(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)

	assert.Equal(t, "GET /foo", serverRequestName(req))
}

func TestRouteRequestNameWithoutContext(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)

	assert.Equal(t, "GET /foo", routeRequestName(req))
}

func TestRouteRequestNameWithContext(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo/123", nil)

	r := mux.NewRouter()
	r.HandleFunc("/foo/{id}", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET /foo/{id}", routeRequestName(r))
	}).Methods(http.MethodGet)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
}

func TestStatusCodeResult(t *testing.T) {
	assert.Equal(t, "HTTP 2xx", statusCodeResult(http.StatusOK))
}