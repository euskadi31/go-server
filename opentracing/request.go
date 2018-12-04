// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package opentracing

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/rs/zerolog/log"
)

// WrapRequest injects an OpenTracing Span found in
// context into the HTTP Headers. If no such Span can be found, WrapRequest
// is a noop.
func WrapRequest(ctx context.Context, req *http.Request) *http.Request {
	// Retrieve the Span from context.
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return req
	}

	// We are going to use this span in a client request, so mark as such.
	ext.SpanKindRPCClient.Set(span)

	// Add some standard OpenTracing tags, useful in an HTTP request.
	ext.HTTPMethod.Set(span, req.Method)
	ext.HTTPUrl.Set(
		span,
		fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.URL.Host, req.URL.Path),
	)

	// Add information on the peer service we're about to contact.
	if host, portString, err := net.SplitHostPort(req.URL.Host); err == nil {
		ext.PeerHostname.Set(span, host)
		if port, err := strconv.Atoi(portString); err == nil {
			ext.PeerPort.Set(span, uint16(port))
		}
	} else {
		ext.PeerHostname.Set(span, req.URL.Host)
	}

	// Inject the Span context into the outgoing HTTP Request.
	if err := opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	); err != nil {
		log.Error().Err(err).Msg("trying to inject span")
	}

	return req
}
