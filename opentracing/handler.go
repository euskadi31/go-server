// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package opentracing

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/rs/zerolog/log"
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

			wireContext, err := tracer.Extract(
				opentracing.TextMap,
				opentracing.HTTPHeadersCarrier(req.Header),
			)
			if err != nil {
				log.Debug().Err(err).Msg("trying to extract span")
			}

			path := routeRequestName(req)

			span := tracer.StartSpan(path, ext.RPCServerOption(wireContext))
			defer span.Finish()

			ext.HTTPUrl.Set(span, req.URL.Path)
			ext.HTTPMethod.Set(span, req.Method)

			params := mux.Vars(req)

			fields := []otlog.Field{}
			for k, v := range params {
				if i, err := strconv.ParseInt(v, 10, 64); err == nil {
					fields = append(fields, otlog.Int64(k, i))
				} else {
					fields = append(fields, otlog.String(k, v))
				}
			}

			if len(fields) > 0 {
				span.LogFields(fields...)
			}

			ctx := opentracing.ContextWithSpan(req.Context(), span)

			req = req.WithContext(ctx)

			next.ServeHTTP(w, req)
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
