// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package encoder

var encoders map[string]Encoder

func init() {
	encoders = make(map[string]Encoder)
}

// Register encoder provider
func Register(encoder Encoder) {
	encoders[encoder.MimeType()] = encoder
}
