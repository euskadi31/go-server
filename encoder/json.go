// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package encoder

import (
	"encoding/json"
	"net/http"
)

type jsonEncoder struct {
}

// NewJSONEncoder construct
func NewJSONEncoder() Encoder {
	return &jsonEncoder{}
}

func (jsonEncoder) MimeType() string {
	return "application/json"
}

func (jsonEncoder) Encode(w http.ResponseWriter, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
}
