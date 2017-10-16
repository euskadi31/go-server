// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// NewHandler instanciates a new mysql HTTP handler.
func NewHandler() func(http.Handler) http.Handler {
	duration := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "The HTTP request latencies in seconds.",
			Buckets: []float64{0.01, 0.1, 0.3, 0.5, 1., 2., 5.},
		},
	)

	prometheus.MustRegister(duration)

	request := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "The count of request.",
		},
		[]string{},
	)

	prometheus.MustRegister(request)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			request.WithLabelValues().Add(1)

			ts := time.Now()

			next.ServeHTTP(w, r)

			duration.Observe(time.Since(ts).Seconds())
		})
	}
}
