// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package response

import (
	"github.com/euskadi31/go-server/response/encoder"
)

var encoders = map[string]encoder.Encoder{}

// Register encoder provider.
func Register(encoder encoder.Encoder) {
	encoders[encoder.MimeType()] = encoder
}
