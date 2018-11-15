// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package opentracing

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/rs/zerolog/log"
	"github.com/zenazn/goji/web/mutil"
)

// RequestIgnorerFunc is the type of a function for use in
// WithServerRequestIgnorer.
type RequestIgnorerFunc func(*http.Request) bool

// IgnoreNone is a RequestIgnorerFunc which ignores no requests.
func IgnoreNone(*http.Request) bool {
	return false
}

// Handler opentracing
func Handler(tracer opentracing.Tracer, ignore RequestIgnorerFunc) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if ignore(req) {
				next.ServeHTTP(w, req)

				return
			}

			lw := mutil.WrapWriter(w)

			wireContext, err := tracer.Extract(
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(req.Header),
			)
			if err != nil {
				log.Debug().Err(err).Msg("trying to extract span")
			}

			path := routeRequestName(req)

			span := tracer.StartSpan(path, ext.RPCServerOption(wireContext))
			defer span.Finish()

			ext.HTTPUrl.Set(span, fmt.Sprintf("http://%s%s", req.Host, req.URL.Path))
			ext.HTTPMethod.Set(span, req.Method)

			params := mux.Vars(req)
			for k, v := range params {
				span.SetTag(k, v)
			}

			ctx := req.Context()
			ctx = opentracing.ContextWithSpan(ctx, span)

			req = req.WithContext(ctx)

			next.ServeHTTP(lw, req)

			ext.HTTPStatusCode.Set(span, uint16(lw.Status()))

			span.SetTag("result", statusCodeResult(lw.Status()))
		})
	}
}

func routeRequestName(req *http.Request) string {
	route := mux.CurrentRoute(req)
	if route != nil {
		tpl, err := route.GetPathTemplate()
		if err == nil {
			return req.Method + " " + massageTemplate(tpl)
		}
	}

	return serverRequestName(req)
}

func serverRequestName(req *http.Request) string {
	var b strings.Builder
	b.Grow(len(req.Method) + len(req.URL.Path) + 1)
	b.WriteString(req.Method)
	b.WriteByte(' ')
	b.WriteString(req.URL.Path)
	return b.String()
}

// massageTemplate removes the regexp patterns from template variables.
func massageTemplate(tpl string) string {
	braces := braceIndices(tpl)
	if len(braces) == 0 {
		return tpl
	}

	buf := bytes.NewBuffer(make([]byte, 0, len(tpl)))
	for i := 0; i < len(tpl); {
		var j int
		if i < braces[0] {
			j = braces[0]
			buf.WriteString(tpl[i:j])
		} else {
			j = braces[1]
			field := tpl[i:j]

			if colon := strings.IndexRune(field, ':'); colon >= 0 {
				buf.WriteString(field[:colon])
				buf.WriteRune('}')
			} else {
				buf.WriteString(field)
			}

			braces = braces[2:]
			if len(braces) == 0 {
				buf.WriteString(tpl[j:])
				break
			}
		}

		i = j
	}

	return buf.String()
}

// Copied/adapted from gorilla/mux. The original version checks
// that the braces are matched up correctly; we assume they are,
// as otherwise the path wouldn't have been registered correctly.
func braceIndices(s string) []int {
	var level, idx int
	var idxs []int

	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '{':
			if level++; level == 1 {
				idx = i
			}
		case '}':
			if level--; level == 0 {
				idxs = append(idxs, idx, i+1)
			}
		}
	}

	return idxs
}

var standardStatusCodeResults = [...]string{
	"HTTP 1xx",
	"HTTP 2xx",
	"HTTP 3xx",
	"HTTP 4xx",
	"HTTP 5xx",
}

// statusCodeResult returns the transaction result value to use for the given
// status code.
func statusCodeResult(statusCode int) string {
	switch i := statusCode / 100; i {
	case 1, 2, 3, 4, 5:
		return standardStatusCodeResults[i-1]
	}

	return fmt.Sprintf("HTTP %d", statusCode)
}
