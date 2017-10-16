// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestParamsFromContext(t *testing.T) {
	_, err := ParamsFromContext(nil)
	assert.EqualError(t, err, errContextIsNull.Error())

	_, err = ParamsFromContext(context.Background())
	assert.EqualError(t, err, errNotFountInContext.Error())

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	handler := stripParams(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := ParamsFromContext(r.Context())

		assert.NoError(t, err)
	}))

	handler(w, req, httprouter.Params{})
}
