// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package server

import (
	"io/ioutil"
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

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouterEnableHealthCheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/health", nil)
	w := httptest.NewRecorder()

	router := NewRouter()

	router.EnableHealthCheck()

	err := router.AddHealthCheck("redis", HealthCheckHandlerFunc(func() bool {
		return true
	}))
	assert.NoError(t, err)

	err = router.AddHealthCheck("redis", HealthCheckHandlerFunc(func() bool {
		return true
	}))
	assert.Error(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	b, err := ioutil.ReadAll(w.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(b), "redis")
}

func TestRouterHealthCheckFailed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/health", nil)
	w := httptest.NewRecorder()

	router := NewRouter()

	router.EnableHealthCheck()

	err := router.AddHealthCheck("redis", HealthCheckHandlerFunc(func() bool {
		return false
	}))
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	b, err := ioutil.ReadAll(w.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(b), "redis")
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

	assert.Equal(t, http.StatusOK, w.Code)
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

	assert.Equal(t, http.StatusOK, w.Code)
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

func TestRouterEnableRecovery(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	req.Header.Add("X-Forwarded-For", "127.0.0.1")

	w := httptest.NewRecorder()

	router := NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	}).Methods(http.MethodGet)

	router.EnableRecovery()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

type testController struct {
}

func (c testController) Mount(r *Router) {
	r.HandleFunc("/controller", c.handler).Methods(http.MethodGet)
}

func (c testController) handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func TestRouterAddController(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/controller", nil)

	w := httptest.NewRecorder()

	router := NewRouter()

	router.AddController(&testController{})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

type testRoute struct {
}

func (testRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func TestRouterAddRoute(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/route", nil)

	w := httptest.NewRecorder()

	router := NewRouter()

	router.AddRoute("/route", &testRoute{})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestRouterAddRouteFunc(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/route-func", nil)

	w := httptest.NewRecorder()

	router := NewRouter()

	router.AddRouteFunc("/route-func", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestRouterAddPrefixRoute(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/route/foo", nil)

	w := httptest.NewRecorder()

	router := NewRouter()

	router.AddPrefixRoute("/route/", &testRoute{})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestRouterAddPrefixRouteFunc(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/route-func/foo", nil)

	w := httptest.NewRecorder()

	router := NewRouter()

	router.AddPrefixRouteFunc("/route-func/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

type testNotFound struct {
}

func (testNotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	if _, err := w.Write([]byte("true")); err != nil {
		panic(err)
	}
}

func TestRouterSetNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/route-func/foo", nil)

	w := httptest.NewRecorder()

	router := NewRouter()
	router.SetNotFound(&testNotFound{})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	b, err := ioutil.ReadAll(w.Body)
	assert.NoError(t, err)
	assert.Equal(t, "true", string(b))
}

func TestRouterSetNotFoundFunc(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/route-func/foo", nil)

	w := httptest.NewRecorder()

	router := NewRouter()
	router.SetNotFoundFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("true")); err != nil {
			panic(err)
		}
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	b, err := ioutil.ReadAll(w.Body)
	assert.NoError(t, err)
	assert.Equal(t, "true", string(b))
}

func BenchmarkRouter(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/route-func", nil)

	router := NewRouter()

	router.AddRouteFunc("/route-func", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	for n := 0; n < b.N; n++ {
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
	}
}
