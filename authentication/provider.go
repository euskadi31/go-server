// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package authentication

import "net/http"

// Provider interface
//go:generate mockery -case=underscore -inpkg -name=Provider
type Provider interface {
	Validate(r *http.Request) bool
}
