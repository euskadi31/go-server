// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouterEnableMetrics(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/metrics", nil)
	w := httptest.NewRecorder()

	router := NewRouter()

	router.EnableMetrics()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestRouterEnableHealthCheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/health", nil)
	w := httptest.NewRecorder()

	router := NewRouter()

	router.EnableHealthCheck()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestRouterEnableCors(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	req.Header.Add("Origin", "http://localhost")

	w := httptest.NewRecorder()

	router := NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	router.EnableCors()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestRouterEnableProxy(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	req.Header.Add("X-Forwarded-For", "127.0.0.1")

	w := httptest.NewRecorder()

	router := NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "127.0.0.1", r.RemoteAddr)
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	router.EnableProxy()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestRouterEnableProfiling(t *testing.T) {
	router := NewRouter()

	router.EnableProfiling()

	for _, item := range []struct {
		status   int
		endpoint string
		method   string
	}{
		{
			status:   http.StatusMovedPermanently,
			endpoint: "/debug/pprof",
			method:   http.MethodGet,
		},
		{
			status:   http.StatusMovedPermanently,
			endpoint: "/debug/pprof/cmdline",
			method:   http.MethodGet,
		},
		{
			status:   http.StatusMovedPermanently,
			endpoint: "/debug/pprof/symbol",
			method:   http.MethodGet,
		},
		{
			status:   http.StatusMovedPermanently,
			endpoint: "/debug/pprof/profile",
			method:   http.MethodGet,
		},
		{
			status:   http.StatusMovedPermanently,
			endpoint: "/debug/pprof/heap",
			method:   http.MethodGet,
		},
		{
			status:   http.StatusMovedPermanently,
			endpoint: "/debug/pprof/goroutine",
			method:   http.MethodGet,
		},
		{
			status:   http.StatusMovedPermanently,
			endpoint: "/debug/pprof/block",
			method:   http.MethodGet,
		},
		{
			status:   http.StatusMovedPermanently,
			endpoint: "/debug/pprof/threadcreate",
			method:   http.MethodGet,
		},
	} {
		req := httptest.NewRequest(item.method, "http://example.com/"+item.endpoint, nil)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, item.status, w.Code)
	}
}
