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
