// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package encoder

import (
	"net/http"

	"github.com/golang/gddo/httputil"
	"github.com/rs/zerolog/log"
)

const defaultMediatype = "application/json"

// Encode data to HTTP response
func Encode(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	offers := make([]string, 0, len(encoders))

	for mime := range encoders {
		offers = append(offers, mime)
	}

	mediatype := httputil.NegotiateContentType(r, offers, defaultMediatype)

	encoder := encoders[mediatype]

	if err := encoder.Encode(w, data); err != nil {
		log.Error().Err(err).Msg("encode content failed")

		status = http.StatusInternalServerError

		//@TODO: write error response
	}

	w.Header().Set("Content-Type", encoder.MimeType()+";charset=utf-8")
	w.WriteHeader(status)
}
