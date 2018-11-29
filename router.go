// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/euskadi31/go-server/metrics"
	"github.com/euskadi31/go-server/response"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Router struct
type Router struct {
	*mux.Router
	healthchecks map[string]HealthCheckHandler
}

// NewRouter constructor
func NewRouter() *Router {
	return &Router{
		Router:       mux.NewRouter(),
		healthchecks: make(map[string]HealthCheckHandler),
	}
}

// AddHealthCheck handler
func (r *Router) AddHealthCheck(name string, handle HealthCheckHandler) error {
	if _, ok := r.healthchecks[name]; ok {
		return fmt.Errorf("the %s healthcheck handler already exists", name)
	}

	r.healthchecks[name] = handle

	return nil
}

// EnableHealthCheck endpoint
func (r *Router) EnableHealthCheck() {
	r.AddRouteFunc("/health", r.healthHandler).Methods(http.MethodGet, "HEAD")
}

func (r *Router) healthHandler(w http.ResponseWriter, req *http.Request) {
	code := http.StatusOK

	resp := healthCheckProcessor(r.healthchecks)

	if !resp.Status {
		code = http.StatusServiceUnavailable
	}

	response.Encode(w, req, code, resp)
}

// EnableMetrics endpoint
func (r *Router) EnableMetrics() {
	r.Use(metrics.Handler())

	r.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)
}

// EnableCors for all endpoint
func (r *Router) EnableCors() {
	r.EnableCorsWithOptions(handlers.AllowedOrigins([]string{
		"*",
	}))
}

// EnableProxy for populating r.RemoteAddr and r.URL.Scheme based on the X-Forwarded-For,
// X-Real-IP, X-Forwarded-Proto and RFC7239 Forwarded headers when running
// a Go server behind a HTTP reverse proxy.
func (r *Router) EnableProxy() {
	r.Use(func() func(h http.Handler) http.Handler {
		return handlers.ProxyHeaders
	}())
}

// EnableProfiling with pprof
func (r *Router) EnableProfiling() {
	r.HandleFunc("/debug/pprof", pprof.Index).Methods(http.MethodGet)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline).Methods(http.MethodGet)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol).Methods(http.MethodGet)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol).Methods(http.MethodPost)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile).Methods(http.MethodGet)

	r.Handle("/debug/pprof/heap", pprof.Handler("heap")).Methods(http.MethodGet)
	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine")).Methods(http.MethodGet)
	r.Handle("/debug/pprof/block", pprof.Handler("block")).Methods(http.MethodGet)
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate")).Methods(http.MethodGet)
}

// EnableRecovery for all endpoint
func (r *Router) EnableRecovery() {
	r.Use(handlers.RecoveryHandler())
}

// EnableCorsWithOptions for all endpoint
func (r *Router) EnableCorsWithOptions(opts ...handlers.CORSOption) {
	r.Use(handlers.CORS(opts...))
}

// AddController to Router
func (r *Router) AddController(controller Controller) {
	controller.Mount(r)
}

// AddRoute to Router
// Deprecated: Use server.Handle() instead.
func (r *Router) AddRoute(path string, handler http.Handler) *mux.Route {
	return r.Handle(path, handler)
}

// AddRouteFunc to Router
// Deprecated: Use server.HandleFunc() instead.
func (r *Router) AddRouteFunc(path string, handler http.HandlerFunc) *mux.Route {
	return r.HandleFunc(path, handler)
}

// AddPrefixRoute to Router
// Deprecated: Use server.PathPrefix(prefix).Handler() instead.
func (r *Router) AddPrefixRoute(prefix string, handler http.Handler) *mux.Route {
	return r.PathPrefix(prefix).Handler(handler)
}

// AddPrefixRouteFunc to Router
// Deprecated: Use server.PathPrefix(prefix).HandlerFunc() instead.
func (r *Router) AddPrefixRouteFunc(prefix string, handler http.HandlerFunc) *mux.Route {
	return r.PathPrefix(prefix).HandlerFunc(handler)
}

// SetNotFound handler
func (r *Router) SetNotFound(handler http.Handler) {
	r.NotFoundHandler = handler
}

// SetNotFoundFunc handler
func (r *Router) SetNotFoundFunc(handler http.HandlerFunc) {
	r.NotFoundHandler = handler
}
