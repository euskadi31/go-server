// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package locale

import (
	"context"
	"net/http"

	"golang.org/x/text/language"
)

// DefaultLanguage var
var DefaultLanguage = "fr"

// DefaultRegion var
var DefaultRegion = "FR"

// DefaultSupported language
var DefaultSupported = []string{
	DefaultLanguage, // fr: first language is fallback
}

type key int

const (
	contextKey key = iota
)

// Handler middleware
func Handler() func(next http.Handler) http.Handler {
	return HandlerWithConfig(DefaultSupported)
}

// HandlerWithConfig middleware
func HandlerWithConfig(languages []string) func(next http.Handler) http.Handler {
	supported := []language.Tag{}

	for _, lang := range languages {
		tag := language.MustParse(lang)

		supported = append(supported, tag)
	}

	var matcher = language.NewMatcher(supported)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			tag, _ := language.MatchStrings(matcher, r.Header.Get("Accept-Language"))

			region, _ := tag.Region()

			locale := Locale{
				Language: tag.String(),
				Region:   region.String(),
			}

			ctx = ToContext(ctx, locale)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ToContext add Locale to Context
func ToContext(ctx context.Context, locale Locale) context.Context {
	return context.WithValue(ctx, contextKey, locale)
}

// FromContext returns Locale from Context
func FromContext(ctx context.Context) Locale {
	value, ok := ctx.Value(contextKey).(Locale)
	if ok {
		return value
	}

	return Locale{
		Language: DefaultLanguage,
		Region:   DefaultRegion,
	}
}
