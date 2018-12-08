// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigurationIsEnabled(t *testing.T) {
	c := &Configuration{}

	assert.False(t, c.IsEnabled("http"))
	assert.False(t, c.IsEnabled("https"))
	assert.False(t, c.IsEnabled("bad"))

	c.HTTP = &HTTPConfiguration{
		Port: 8080,
	}

	assert.True(t, c.IsEnabled("http"))
	assert.False(t, c.IsEnabled("https"))

	c.HTTPS = &HTTPSConfiguration{
		Port:     443,
		CertFile: "foo.cert",
		KeyFile:  "foo.key",
	}

	assert.True(t, c.IsEnabled("http"))
	assert.True(t, c.IsEnabled("https"))
}

func TestHTTPConfigurationAddr(t *testing.T) {
	c := &HTTPConfiguration{
		Host: "127.0.0.1",
		Port: 8080,
	}

	assert.Equal(t, "127.0.0.1:8080", c.Addr())
}

func TestHTTPSConfigurationAddr(t *testing.T) {
	c := &HTTPSConfiguration{
		Host: "127.0.0.1",
		Port: 8080,
	}

	assert.Equal(t, "127.0.0.1:8080", c.Addr())
}

func TestConfigurationWithDefault(t *testing.T) {

	cfg := ConfigurationWithDefault(nil)

	assert.Equal(t, DefaultIdleTimeout, cfg.IdleTimeout)
	assert.Equal(t, DefaultReadTimeout, cfg.ReadTimeout)
	assert.Equal(t, DefaultShutdownTimeout, cfg.ShutdownTimeout)
	assert.Equal(t, DefaultWriteTimeout, cfg.WriteTimeout)
	assert.NotNil(t, cfg.HTTPS)
	assert.NotNil(t, cfg.HTTPS.TLSConfig)
	assert.True(t, cfg.HTTPS.TLSConfig.PreferServerCipherSuites)
	assert.Equal(t, DefaultMinVersion, cfg.HTTPS.TLSConfig.MinVersion)
}
