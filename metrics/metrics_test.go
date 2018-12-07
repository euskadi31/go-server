// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/justinas/alice"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	_ = Handler()

	middleware := alice.New(Handler()).ThenFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
