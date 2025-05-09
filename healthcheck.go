// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package server

import (
	"sync"
)

// HealthCheckHandler type.
//
//go:generate mockery -case=underscore -inpkg -name=HealthCheckHandler
type HealthCheckHandler interface {
	Check() bool
}

// HealthCheckHandlerFunc handler.
type HealthCheckHandlerFunc func() bool

// Check calls f().
func (f HealthCheckHandlerFunc) Check() bool {
	return f()
}

// HealthCheckResponse struct.
type HealthCheckResponse struct {
	Status   bool            `json:"status"`
	Services map[string]bool `json:"services"`
}

func healthCheckProcessor(healthchecks map[string]HealthCheckHandler) HealthCheckResponse {
	response := HealthCheckResponse{
		Status:   true,
		Services: make(map[string]bool, len(healthchecks)),
	}

	var (
		wg    = &sync.WaitGroup{}
		mutex = &sync.Mutex{}
	)

	for name, handle := range healthchecks {
		wg.Add(1)

		go func(n string, h HealthCheckHandler) {
			defer wg.Done()

			mutex.Lock()
			defer mutex.Unlock()

			s := h.Check()

			response.Services[n] = s

			if !s {
				response.Status = false
			}
		}(name, handle)
	}

	wg.Wait()

	return response
}
