// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package opentracing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/justinas/alice"
	"github.com/opentracing/opentracing-go/mocktracer"
)

func TestWrapRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	tracer := mocktracer.New()

	middleware := alice.New(Handler(tracer, IgnoreNone)).ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		subReq := httptest.NewRequest("GET", "http://example.com/foo", nil)
		subReq = WrapRequest(r.Context(), subReq)
	})

	middleware.ServeHTTP(w, req)
}
