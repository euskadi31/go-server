// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package authentication

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/euskadi31/go-server"
	"github.com/rs/zerolog/log"
)

// Handler authentication
func Handler(config *Configuration, provider Provider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Skip public endpoints
			if r.URL.Path == "/health" || r.URL.Path == "/metrics" {
				next.ServeHTTP(w, r)

				return
			}

			if !provider.Validate(r) {
				w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Bearer realm="%s"`, config.Realm))

				log.Error().Msg("Access token invalid or expired")

				server.FailureFromError(w, http.StatusUnauthorized, errors.New("Unauthorized"))

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
