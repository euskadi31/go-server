// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type DemoController struct {
}

func (c DemoController) Register(r *Router) {
	r.GET("/", c.GETIndexHandle)
}

func (c DemoController) GETIndexHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

func TestAddHealthcheck(t *testing.T) {
	config := &Configuration{}

	server := NewRouter(config)

	server.AddHealthCheck("mysql", func(ctx context.Context) bool {
		return true
	})

	assert.Equal(t, 1, len(server.healthchecks))
}

func TestCors(t *testing.T) {
	config := &Configuration{}

	server := NewRouter(config)
	server.EnableCors()
	server.AddController(DemoController{})

	ts := httptest.NewServer(server)
	defer ts.Close()

	client := &http.Client{}

	req, err := http.NewRequest("GET", ts.URL+"/", nil)
	assert.NoError(t, err)

	req.Header.Set("Origin", "http://test.com")

	res, err := client.Do(req)
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "*", res.Header.Get("Access-Control-Allow-Origin"))
}

func TestMetricsEndpoint(t *testing.T) {
	config := &Configuration{}

	server := NewRouter(config)
	server.EnableMetrics()

	ts := httptest.NewServer(server)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/metrics")
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)
}

func TestHealthEndpoint(t *testing.T) {
	config := &Configuration{}

	server := NewRouter(config)
	server.EnableHealthCheck()
	err := server.AddHealthCheck("mysql", func(ctx context.Context) bool {
		return true
	})

	assert.NoError(t, err)

	err = server.AddHealthCheck("mysql", func(ctx context.Context) bool {
		return true
	})
	assert.Error(t, err)

	ts := httptest.NewServer(server)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/health")
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)
}

func TestHealthEndpointFail(t *testing.T) {
	config := &Configuration{}

	server := NewRouter(config)
	server.EnableHealthCheck()
	err := server.AddHealthCheck("mysql", func(ctx context.Context) bool {
		return false
	})
	assert.NoError(t, err)

	ts := httptest.NewServer(server)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/health")
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 503, res.StatusCode)
}

func TestAddController(t *testing.T) {
	config := &Configuration{}

	server := NewRouter(config)
	server.AddController(DemoController{})
	ts := httptest.NewServer(server)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)
}

func TestPOST(t *testing.T) {
	config := &Configuration{}

	server := NewRouter(config)

	server.POST("/users", func(w http.ResponseWriter, r *http.Request) {
		_, err := ParamsFromContext(r.Context())
		assert.NoError(t, err)
	})

	req := httptest.NewRequest("POST", "http://example.com/users", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)
}

func TestPUT(t *testing.T) {
	config := &Configuration{}

	server := NewRouter(config)

	server.PUT("/users", func(w http.ResponseWriter, r *http.Request) {
		_, err := ParamsFromContext(r.Context())
		assert.NoError(t, err)
	})

	req := httptest.NewRequest("PUT", "http://example.com/users", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)
}

func TestHEAD(t *testing.T) {
	config := &Configuration{}

	server := NewRouter(config)

	server.HEAD("/users", func(w http.ResponseWriter, r *http.Request) {
		_, err := ParamsFromContext(r.Context())
		assert.NoError(t, err)
	})

	req := httptest.NewRequest("HEAD", "http://example.com/users", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)
}

func TestOPTIONS(t *testing.T) {
	config := &Configuration{}

	server := NewRouter(config)

	server.OPTIONS("/users", func(w http.ResponseWriter, r *http.Request) {
		_, err := ParamsFromContext(r.Context())
		assert.NoError(t, err)
	})

	req := httptest.NewRequest("OPTIONS", "http://example.com/users", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)
}

func TestPATCH(t *testing.T) {
	config := &Configuration{}

	server := NewRouter(config)

	server.PATCH("/users", func(w http.ResponseWriter, r *http.Request) {
		_, err := ParamsFromContext(r.Context())
		assert.NoError(t, err)
	})

	req := httptest.NewRequest("PATCH", "http://example.com/users", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)
}

func TestDELETE(t *testing.T) {
	config := &Configuration{}

	server := NewRouter(config)

	server.DELETE("/users", func(w http.ResponseWriter, r *http.Request) {
		_, err := ParamsFromContext(r.Context())
		assert.NoError(t, err)

		JSON(w, http.StatusOK, 1)
	})

	req := httptest.NewRequest("DELETE", "http://example.com/users", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)
}

func TestGET(t *testing.T) {
	config := &Configuration{}

	server := NewRouter(config)

	server.GET("/users/:id", func(w http.ResponseWriter, r *http.Request) {
		params, err := ParamsFromContext(r.Context())
		assert.NoError(t, err)

		assert.Equal(t, "345", params.ByName("id"))

		JSON(w, http.StatusOK, 1)
	})

	req := httptest.NewRequest("GET", "http://example.com/users/345", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)
}

func TestHandle(t *testing.T) {
	config := &Configuration{}

	server := NewRouter(config)

	server.Handle("GET", "/users/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, err := ParamsFromContext(r.Context())
		assert.NoError(t, err)

		assert.Equal(t, "123", params.ByName("id"))
	}))

	req := httptest.NewRequest("GET", "http://example.com/users/123", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)
}

/*
func TestHandleWithMiddleware(t *testing.T) {

	middleware := alice.New(tagMiddleware("m2"))

	subMiddleware := middleware.Append(tagMiddleware("m3"))

	config := &Configuration{}

	server := NewRouter(config)
	server.Use(tagMiddleware("m1"))

	server.Handle("GET", "/users/:id", subMiddleware.ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		params, err := ParamsFromContext(r.Context())
		assert.NoError(t, err)

		assert.Equal(t, "123", params.ByName("id"))
	}))

	req := httptest.NewRequest("GET", "http://example.com/users/123", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, "m1m2m3", w.Body.String())
}
*/

func TestRun(t *testing.T) {
	config := &Configuration{
		Host: "127.0.0.1",
		Port: 57474,
	}

	server := NewRouter(config)
	server.GET("/", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, http.StatusOK, 1)
	})

	go server.ListenAndServe()

	time.Sleep(time.Millisecond * 100)

	res, err := http.Get("http://127.0.0.1:57474/")
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)
}
