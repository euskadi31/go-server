// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package server

import (
	"crypto/tls"
	"net/http"
	"testing"
	//"time"

	"github.com/stretchr/testify/assert"
)

func httpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

func TestServerNotConfigured(t *testing.T) {
	s := New(&Configuration{})

	err := s.Run()
	assert.EqualError(t, err, "http or https server is not configured")
}

func TestServerHTTP(t *testing.T) {
	s := New(&Configuration{
		HTTP: &HTTPConfiguration{
			Port: 12456,
		},
		Profiling:   true,
		Metrics:     true,
		HealthCheck: true,
	})

	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	go func() {
		err := s.Run()
		assert.NoError(t, err)
	}()

	defer func() {
		err := s.Shutdown()
		assert.NoError(t, err)
	}()

	resp, err := http.Get("http://localhost:12456/")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = http.Get("http://localhost:12456/health")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestServerHTTPS(t *testing.T) {
	s := New(&Configuration{
		HTTPS: &HTTPSConfiguration{
			Port:     12457,
			CertFile: "./testdata/server.crt",
			KeyFile:  "./testdata/server.key",
		},
		Profiling:   true,
		Metrics:     true,
		HealthCheck: true,
	})

	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	go func() {
		err := s.Run()
		assert.NoError(t, err)
	}()

	defer func() {
		err := s.Shutdown()
		assert.NoError(t, err)
	}()

	client := httpClient()

	resp, err := client.Get("https://localhost:12457/")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = client.Get("https://localhost:12457/health")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestServerHTTPAndHTTPS(t *testing.T) {
	s := New(&Configuration{
		HTTP: &HTTPConfiguration{
			Port: 12456,
		},
		HTTPS: &HTTPSConfiguration{
			Port:     12457,
			CertFile: "./testdata/server.crt",
			KeyFile:  "./testdata/server.key",
		},
		Profiling:   true,
		Metrics:     true,
		HealthCheck: true,
	})

	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	go func() {
		err := s.Run()
		assert.NoError(t, err)
	}()

	defer func() {
		err := s.Shutdown()
		assert.NoError(t, err)
	}()

	client := httpClient()

	resp, err := client.Get("http://localhost:12456/")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = client.Get("https://localhost:12457/")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = client.Get("https://localhost:12457/health")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
