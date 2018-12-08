// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package server

import (
	"crypto/tls"
	"net"
	"strconv"
	"time"
)

var (
	// DefaultCurvePreferences defines the recommended elliptic curves for modern TLS
	DefaultCurvePreferences = []tls.CurveID{
		tls.CurveP256,
		tls.X25519, // Go 1.8 only
	}

	// DefaultCipherSuites defines the recommended cipher suites for modern TLS
	DefaultCipherSuites = []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

		// Best disabled, as they don't provide Forward Secrecy,
		// but might be necessary for some clients
		// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	}

	// DefaultMinVersion defines the recommended minimum version to use for the TLS protocol (1.2)
	DefaultMinVersion uint16 = tls.VersionTLS12

	// DefaultReadTimeout sets the maximum time a client has to fully stream a request (5s)
	DefaultReadTimeout = 5 * time.Second
	// DefaultWriteTimeout sets the maximum amount of time a handler has to fully process a request (10s)
	DefaultWriteTimeout = 10 * time.Second
	// DefaultIdleTimeout sets the maximum amount of time a Keep-Alive connection can remain idle before
	// being recycled (120s)
	DefaultIdleTimeout = 120 * time.Second

	// DefaultShutdownTimeout sets the maximum amount of time a shutdown
	DefaultShutdownTimeout = 2 * time.Second
)

// Configuration struct
type Configuration struct {
	HTTP              *HTTPConfiguration
	HTTPS             *HTTPSConfiguration
	ShutdownTimeout   time.Duration
	WriteTimeout      time.Duration
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	IdleTimeout       time.Duration
	Profiling         bool
	Metrics           bool
	HealthCheck       bool
}

// ConfigurationWithDefault return Configuration with default parameters
func ConfigurationWithDefault(cfg *Configuration) *Configuration {
	if cfg == nil {
		cfg = &Configuration{}
	}

	if cfg.HTTPS == nil {
		cfg.HTTPS = &HTTPSConfiguration{}
	}

	if cfg.HTTPS.TLSConfig == nil {
		cfg.HTTPS.TLSConfig = &tls.Config{}
	}

	cfg.HTTPS.TLSConfig.PreferServerCipherSuites = true
	cfg.HTTPS.TLSConfig.MinVersion = DefaultMinVersion
	cfg.HTTPS.TLSConfig.CurvePreferences = DefaultCurvePreferences
	cfg.HTTPS.TLSConfig.CipherSuites = DefaultCipherSuites

	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = DefaultReadTimeout
	}

	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = DefaultWriteTimeout
	}

	if cfg.IdleTimeout == 0 {
		cfg.IdleTimeout = DefaultIdleTimeout
	}

	if cfg.ShutdownTimeout == 0 {
		cfg.ShutdownTimeout = DefaultShutdownTimeout
	}

	return cfg
}

// IsEnabled check if protocol is enabled
func (c Configuration) IsEnabled(protocol string) bool {
	switch protocol {
	case "http":
		return c.HTTP != nil && c.HTTP.IsEnabled()
	case "https":
		return c.HTTPS != nil && c.HTTPS.IsEnabled()
	default:
		return false
	}
}

// HTTPConfiguration struct
type HTTPConfiguration struct {
	Host string
	Port int
}

// Addr string
func (c HTTPConfiguration) Addr() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

// IsEnabled check if HTTP is enabled
func (c HTTPConfiguration) IsEnabled() bool {
	return c.Port > 0 && c.Port < 65535
}

// HTTPSConfiguration struct
type HTTPSConfiguration struct {
	Host      string
	Port      int
	TLSConfig *tls.Config
	CertFile  string
	KeyFile   string
}

// Addr string
func (c HTTPSConfiguration) Addr() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

// IsEnabled check if HTTP is enabled
func (c HTTPSConfiguration) IsEnabled() bool {
	return c.Port > 0 && c.Port < 65535 && c.CertFile != "" && c.KeyFile != ""
}
