// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
)

// Server struct.
type Server struct {
	*Router
	cfg         *Configuration
	httpServer  *http.Server
	httpsServer *http.Server
}

// New Server.
func New(cfg *Configuration) *Server {
	return &Server{
		Router: NewRouter(),
		cfg:    cfg,
	}
}

func (s *Server) runHTTPServer() error {
	addr := s.cfg.HTTP.Addr()

	log.Info().Msgf("HTTP Server running on %s", addr)

	s.httpServer = &http.Server{
		Addr:              addr,
		Handler:           s.Router,
		ReadTimeout:       s.cfg.ReadTimeout,
		ReadHeaderTimeout: s.cfg.ReadHeaderTimeout,
		WriteTimeout:      s.cfg.WriteTimeout,
		IdleTimeout:       s.cfg.IdleTimeout,
	}

	return s.httpServer.ListenAndServe() // nolint: wrapcheck
}

func (s *Server) runHTTPSServer() error {
	addr := s.cfg.HTTPS.Addr()

	log.Info().Msgf("HTTPS Server running on %s", addr)

	s.httpsServer = &http.Server{
		Addr:              addr,
		Handler:           s.Router,
		TLSConfig:         s.cfg.HTTPS.TLSConfig,
		ReadTimeout:       s.cfg.ReadTimeout,
		ReadHeaderTimeout: s.cfg.ReadHeaderTimeout,
		WriteTimeout:      s.cfg.WriteTimeout,
		IdleTimeout:       s.cfg.IdleTimeout,
	}

	return s.httpsServer.ListenAndServeTLS(s.cfg.HTTPS.CertFile, s.cfg.HTTPS.KeyFile) // nolint: wrapcheck
}

// Run Server.
func (s *Server) Run() (err error) {
	if !s.cfg.IsEnabled("http") && !s.cfg.IsEnabled("https") {
		return errors.New("http or https server is not configured")
	}

	s.EnableRecovery()

	if s.cfg.HealthCheck {
		s.EnableHealthCheck()
	}

	if s.cfg.Metrics {
		s.EnableMetrics()
	}

	if s.cfg.Profiling {
		s.EnableProfiling()
	}

	if s.cfg.IsEnabled("http") {
		go func() {
			if e := s.runHTTPServer(); e != nil {
				err = e
			}
		}()
	}

	if s.cfg.IsEnabled("https") {
		go func() {
			if e := s.runHTTPSServer(); e != nil {
				err = e
			}
		}()
	}

	ch := make(chan struct{})
	<-ch

	return nil
}

// Shutdown server.
func (s *Server) Shutdown() (err error) {
	if s.httpServer != nil {
		log.Info().Msg("Shutting down HTTP server...")

		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()

		if e := s.httpServer.Shutdown(ctx); err != nil {
			err = e
		}
	}

	if s.httpsServer != nil {
		log.Info().Msg("Shutting down HTTPS server...")

		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()

		if e := s.httpsServer.Shutdown(ctx); err != nil {
			err = e
		}
	}

	return
}
