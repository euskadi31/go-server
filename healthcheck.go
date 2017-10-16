// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"sync"
)

// HealthCheckHandler type
type HealthCheckHandler func(context.Context) bool

// HealthCheckResponse struct
type HealthCheckResponse struct {
	Status   bool            `json:"status"`
	Services map[string]bool `json:"services"`
}

func healthCheckProcessor(ctx context.Context, healthchecks map[string]HealthCheckHandler) HealthCheckResponse {
	response := HealthCheckResponse{
		Status:   true,
		Services: make(map[string]bool),
	}

	var wg = &sync.WaitGroup{}
	var mutex = &sync.Mutex{}
	for name, handle := range healthchecks {
		wg.Add(1)
		go func(n string, h HealthCheckHandler) {
			defer wg.Done()
			mutex.Lock()
			defer mutex.Unlock()

			s := h(ctx)

			response.Services[n] = s

			if s == false {
				response.Status = false
			}
		}(name, handle)
	}

	wg.Wait()

	return response
}
