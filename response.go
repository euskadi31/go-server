// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/rs/zerolog/log"
)

// ResponseStatus returns the HTTP response status.
// Remember that the status is only set by the server after WriteHeader has been called.
func ResponseStatus(w http.ResponseWriter) int {
	return int(httpResponseStruct(reflect.ValueOf(w)).FieldByName("status").Int())
}

// httpResponseStruct returns the response structure after going trough all the intermediary response writers.
func httpResponseStruct(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Type().String() == "http.response" {
		return v
	}

	return httpResponseStruct(v.FieldByName("ResponseWriter").Elem())
}

// NotFoundFailure response
func NotFoundFailure(w http.ResponseWriter, r *http.Request) {
	Failure(w, http.StatusNotFound, ErrorMessage{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf(`No route found for "%s %s"`, r.Method, r.URL.Path),
	})
}

// MethodNotAllowedFailure response
func MethodNotAllowedFailure(w http.ResponseWriter, r *http.Request) {
	Failure(w, http.StatusMethodNotAllowed, ErrorMessage{
		Code:    http.StatusMethodNotAllowed,
		Message: fmt.Sprintf(`Method "%s" not allowed for "%s"`, r.Method, r.URL.Path),
	})
}

// InternalServerFailure response
func InternalServerFailure(w http.ResponseWriter, r *http.Request, p interface{}) {
	Failure(w, http.StatusInternalServerError, ErrorMessage{
		Code:    http.StatusInternalServerError,
		Message: "Internal Server Error",
	})
}

// FailureFromError write ErrorMessage from error
func FailureFromError(w http.ResponseWriter, status int, err error) {
	Failure(w, status, ErrorMessage{
		Code:    status,
		Message: err.Error(),
	})
}

// Failure response
func Failure(w http.ResponseWriter, status int, err ErrorMessage) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)

	body := ErrorResponse{
		Error: err,
	}

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Error().Err(err)
	}
}

// JSON response
func JSON(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		FailureFromError(w, http.StatusInternalServerError, err)
	}
}
