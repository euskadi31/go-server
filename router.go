// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"net/http"

	"github.com/euskadi31/go-server/metrics"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

// Router struct
type Router struct {
	*mux.Router
	middleware   alice.Chain
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
	r.AddRouteFunc("/health", r.healthHandler).Methods("GET", "HEAD")
}

func (r *Router) healthHandler(w http.ResponseWriter, req *http.Request) {
	code := http.StatusOK

	response := healthCheckProcessor(req.Context(), r.healthchecks)

	if !response.Status {
		code = http.StatusServiceUnavailable
	}

	JSON(w, code, response)
}

// EnableMetrics endpoint
func (r *Router) EnableMetrics() {
	r.Use(metrics.NewHandler())

	r.Handle("/metrics", promhttp.Handler()).Methods("GET")
}

// EnableCors for all endpoint
func (r *Router) EnableCors() {
	r.EnableCorsWithOptions(cors.Options{
		AllowedOrigins: []string{"*"},
	})
}

// EnableRecovery for all endpoint
func (r *Router) EnableRecovery() {
	r.Use(handlers.RecoveryHandler())
}

// EnableCorsWithOptions for all endpoint
func (r *Router) EnableCorsWithOptions(options cors.Options) {
	c := cors.New(options)

	r.Use(c.Handler)

	r.MethodNotAllowedHandler = c.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))
}

// Use middleware
func (r *Router) Use(middleware ...alice.Constructor) {
	r.middleware = r.middleware.Append(middleware...)
}

// AddController to Router
func (r *Router) AddController(controller Controller) {
	controller.Mount(r)
}

// AddRoute to Router
func (r *Router) AddRoute(path string, handler http.Handler) *mux.Route {
	return r.Handle(path, r.middleware.Then(handler))
}

// AddRouteFunc to Router
func (r *Router) AddRouteFunc(path string, handler http.HandlerFunc) *mux.Route {
	return r.Handle(path, r.middleware.ThenFunc(handler))
}

// AddPrefixRoute to Router
func (r *Router) AddPrefixRoute(prefix string, handler http.Handler) *mux.Route {
	return r.PathPrefix(prefix).Handler(r.middleware.Then(handler))
}

// AddPrefixRouteFunc to Router
func (r *Router) AddPrefixRouteFunc(prefix string, handler http.HandlerFunc) *mux.Route {
	return r.PathPrefix(prefix).Handler(r.middleware.ThenFunc(handler))
}

// SetNotFound handler
func (r *Router) SetNotFound(handler http.Handler) {
	r.NotFoundHandler = r.middleware.Then(handler)
}

// SetNotFoundFunc handler
func (r *Router) SetNotFoundFunc(handler http.HandlerFunc) {
	r.NotFoundHandler = r.middleware.ThenFunc(handler)
}
