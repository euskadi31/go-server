// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package opentracing

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/rs/zerolog/log"
)

// Handler opentracing
func Handler(tracer opentracing.Tracer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			wireContext, err := tracer.Extract(
				opentracing.TextMap,
				opentracing.HTTPHeadersCarrier(req.Header),
			)
			if err != nil {
				log.Debug().Err(err).Msg("error encountered while trying to extract span")
			}

			span := tracer.StartSpan(req.URL.Path, ext.RPCServerOption(wireContext))
			defer span.Finish()

			ctx := opentracing.ContextWithSpan(req.Context(), span)

			req = req.WithContext(ctx)

			next.ServeHTTP(w, req)
		})
	}
}
