// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"testing"
)

func TestHealthcheckProcessor(t *testing.T) {
	healthchecks := make(map[string]HealthCheckHandler)

	healthchecks["mysql"] = func(ctx context.Context) bool {
		return true
	}

	response := healthCheckProcessor(context.Background(), healthchecks)

	if response.Status != true {
		t.Error("response.Status is not true")
	}

	if response.Services["mysql"] != true {
		t.Error(`response.Services["mysql"] is not true`)
	}
}

func TestHealthcheckProcessorWithFailedCheck(t *testing.T) {
	healthchecks := make(map[string]HealthCheckHandler)

	healthchecks["mysql"] = func(ctx context.Context) bool {
		return true
	}

	healthchecks["redis"] = func(ctx context.Context) bool {
		return false
	}

	response := healthCheckProcessor(context.Background(), healthchecks)

	if response.Status != false {
		t.Error("response.Status is not false")
	}

	if response.Services["mysql"] != true {
		t.Error(`response.Services["mysql"] is not true`)
	}

	if response.Services["redis"] != false {
		t.Error(`response.Services["redis"] is not false`)
	}
}
