// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package encoder

import (
	"mime"
	"net/http"

	"github.com/rs/zerolog/log"
)

const defaultMediatype = "application/json"

// Encode data to HTTP response
func Encode(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	accept := r.Header.Get("Accept")

	mediatype, _, err := mime.ParseMediaType(accept)
	if err != nil {
		log.Warn().Msgf("ParseMediaType failed: %s, using default type %s", accept, defaultMediatype)

		mediatype = defaultMediatype
	}

	encoder, ok := encoders[mediatype]
	if !ok {
		log.Warn().Msgf("invalid accept type %s, using default type %s", mediatype, defaultMediatype)

		encoder = encoders[defaultMediatype]
	}

	if err := encoder.Encode(w, data); err != nil {
		log.Error().Err(err).Msg("encode content failed")

		status = http.StatusInternalServerError

		//@TODO: write error response
	}

	w.Header().Set("Content-Type", encoder.MimeType()+";charset=utf-8")
	w.WriteHeader(status)
}
