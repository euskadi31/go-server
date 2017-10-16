// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var (
	errContextIsNull     = errors.New("The context is null")
	errNotFountInContext = errors.New("The Params is not found in context")
)

type key int

const (
	paramsKey key = iota
)

// ParamsFromContext returns the httprouter.Params.
func ParamsFromContext(ctx context.Context) (httprouter.Params, error) {
	if ctx == nil {
		return nil, errContextIsNull
	}

	params, ok := ctx.Value(paramsKey).(httprouter.Params)
	if !ok {
		return nil, errNotFountInContext
	}

	return params, nil
}

func stripParams(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), paramsKey, ps)

		h.ServeHTTP(w, r.WithContext(ctx))
	}
}
