// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"github.com/zenazn/goji/web/mutil"
)

// Handler instanciates a new mysql HTTP handler.
func Handler() func(http.Handler) http.Handler {
	duration := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "The HTTP request latencies in seconds.",
			Buckets: []float64{0.01, 0.1, 0.3, 0.5, 1., 2., 5.},
		},
	)

	if err := prometheus.Register(duration); err != nil {
		log.Debug().Err(err).Msg("prometheus register duration")
	}

	request := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "The count of request.",
		},
		[]string{"status", "method"},
	)

	if err := prometheus.Register(request); err != nil {
		log.Debug().Err(err).Msg("prometheus register request")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lw := mutil.WrapWriter(w)

			ts := time.Now()

			next.ServeHTTP(lw, r)

			duration.Observe(time.Since(ts).Seconds())

			request.With(prometheus.Labels{
				"status": strconv.FormatInt(int64(lw.Status()), 10),
				"method": r.Method,
			}).Add(1)
		})
	}
}
