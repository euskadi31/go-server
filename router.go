// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/euskadi31/go-server/metrics"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
)

// Router struct
type Router struct {
	*httprouter.Router
	chain        alice.Chain
	config       *Configuration
	healthchecks map[string]HealthCheckHandler
}

// NewRouter with configuration
func NewRouter(config *Configuration) *Router {
	return &Router{
		Router: &httprouter.Router{
			RedirectTrailingSlash:  true,
			RedirectFixedPath:      true,
			HandleMethodNotAllowed: true,
			NotFound:               NotFoundFailure,
			MethodNotAllowed:       MethodNotAllowedFailure,
			PanicHandler:           InternalServerFailure,
		},
		chain:        alice.New(),
		config:       config,
		healthchecks: make(map[string]HealthCheckHandler),
	}
}

// Use middleware
func (r *Router) Use(constructors ...alice.Constructor) {
	r.chain = r.chain.Append(constructors...)
}

// AddController to Router
func (r *Router) AddController(c Controller) {
	c.Register(r)
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
	r.GET("/health", r.healthHandler)
}

// EnableMetrics endpoint
func (r *Router) EnableMetrics() {
	r.Use(metrics.NewHandler())

	r.Router.Handle("GET", "/metrics", stripParams(promhttp.Handler()))
}

// EnableCors for all endpoint
func (r *Router) EnableCors() {
	r.EnableCorsWithOptions(cors.Options{
		AllowedOrigins: []string{"*"},
	})
}

// EnableCorsWithOptions for all endpoint
func (r *Router) EnableCorsWithOptions(options cors.Options) {
	c := cors.New(options)

	r.Use(c.Handler)
}

func (r *Router) healthHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	code := http.StatusOK

	response := healthCheckProcessor(req.Context(), r.healthchecks)

	if !response.Status {
		code = http.StatusServiceUnavailable
	}

	JSON(w, code, response)
}

// ListenAndServe HTTP Server
func (r *Router) ListenAndServe() error {
	addr := fmt.Sprintf("%s:%d", r.config.Host, r.config.Port)

	s := &http.Server{
		Addr:           addr,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		IdleTimeout:    30 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        r,
	}

	log.Info().Msgf("Listening and serving on %s", addr)

	return s.ListenAndServe()
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (r *Router) GET(path string, handle http.HandlerFunc) {
	r.Handle("GET", path, http.HandlerFunc(handle))
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (r *Router) HEAD(path string, handle http.HandlerFunc) {
	r.Handle("HEAD", path, http.HandlerFunc(handle))
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (r *Router) OPTIONS(path string, handle http.HandlerFunc) {
	r.Handle("OPTIONS", path, http.HandlerFunc(handle))
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (r *Router) POST(path string, handle http.HandlerFunc) {
	r.Handle("POST", path, http.HandlerFunc(handle))
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (r *Router) PUT(path string, handle http.HandlerFunc) {
	r.Handle("PUT", path, http.HandlerFunc(handle))
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (r *Router) PATCH(path string, handle http.HandlerFunc) {
	r.Handle("PATCH", path, http.HandlerFunc(handle))
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (r *Router) DELETE(path string, handle http.HandlerFunc) {
	r.Handle("DELETE", path, http.HandlerFunc(handle))
}

// Handle registers a new request handle with the given path and method.
func (r *Router) Handle(method, path string, handle http.Handler) {
	r.Router.Handle(method, path, stripParams(r.chain.Then(handle)))
}
