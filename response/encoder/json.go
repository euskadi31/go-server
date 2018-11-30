// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package encoder

import (
	"encoding/json"
	"io"
)

type jsonEncoder struct {
}

// JSON construct
func JSON() Encoder {
	return &jsonEncoder{}
}

func (jsonEncoder) MimeType() string {
	return "application/json"
}

func (jsonEncoder) Encode(w io.Writer, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
}
