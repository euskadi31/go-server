// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package encoder

import (
	"net/http"
)

// Encoder provider
//go:generate mockery -case=underscore -inpkg -name=Encoder
type Encoder interface {
	MimeType() string
	Encode(w http.ResponseWriter, data interface{}) error
}
