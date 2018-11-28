// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package server

import (
	"testing"
)

func TestHealthcheckProcessor(t *testing.T) {
	healthchecks := make(map[string]HealthCheckHandler)

	mysql := &MockHealthCheckHandler{}
	mysql.On("Check").Return(true)

	healthchecks["mysql"] = mysql

	response := healthCheckProcessor(healthchecks)

	if response.Status != true {
		t.Error("response.Status is not true")
	}

	if response.Services["mysql"] != true {
		t.Error(`response.Services["mysql"] is not true`)
	}
}

func TestHealthcheckProcessorWithFailedCheck(t *testing.T) {
	healthchecks := make(map[string]HealthCheckHandler)

	mysql := &MockHealthCheckHandler{}
	mysql.On("Check").Return(true)

	healthchecks["mysql"] = mysql

	redis := &MockHealthCheckHandler{}
	redis.On("Check").Return(false)

	healthchecks["redis"] = redis

	response := healthCheckProcessor(healthchecks)

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
