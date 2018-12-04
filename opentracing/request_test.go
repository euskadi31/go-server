// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package opentracing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/justinas/alice"
	"github.com/opentracing/opentracing-go/mocktracer"
)

func TestWrapRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	tracer := mocktracer.New()

	middleware := alice.New(Handler(tracer, IgnoreNone)).ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		subReq := httptest.NewRequest("GET", "http://example.com/foo", nil)
		subReq2 := WrapRequest(r.Context(), subReq)

		assert.Equal(t, subReq, subReq2)
	})

	middleware.ServeHTTP(w, req)
}

func TestWrapRequestWithPort(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	tracer := mocktracer.New()

	middleware := alice.New(Handler(tracer, IgnoreNone)).ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		subReq := httptest.NewRequest("GET", "http://example.com:8080/foo", nil)
		subReq2 := WrapRequest(r.Context(), subReq)

		assert.Equal(t, subReq, subReq2)
	})

	middleware.ServeHTTP(w, req)
}

func TestWrapRequestWithoutTracer(t *testing.T) {
	r := httptest.NewRequest("GET", "http://example.com/foo", nil)
	r2 := WrapRequest(r.Context(), r)

	assert.Equal(t, r, r2)
}
