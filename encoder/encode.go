// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package encoder

import (
	"net/http"
)

// Encode data to HTTP response
func Encode(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	//@TODO: use Accept header for choose encoder
	encoder, ok := encoders["application/json"]
	if !ok {
		//@TODO: error
	}

	if err := encoder.Encode(w, data); err != nil {
		//@TODO: error
	}

	w.WriteHeader(status)
}
