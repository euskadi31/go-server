// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package response

import (
	"net/http"

	"github.com/euskadi31/go-server/response/encoder"
	"github.com/golang/gddo/httputil"
)

const defaultMediatype = "application/json"

func init() {
	Register(encoder.JSON())
}

// Encode data to HTTP response
func Encode(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	offers := make([]string, 0, len(encoders))

	for mime := range encoders {
		offers = append(offers, mime)
	}

	mediatype := httputil.NegotiateContentType(r, offers, defaultMediatype)

	encoder := encoders[mediatype]

	w.Header().Set("Content-Type", encoder.MimeType()+"; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)

	if err := encoder.Encode(w, data); err != nil {
		FailureFromError(w, http.StatusInternalServerError, err)

		return
	}
}
