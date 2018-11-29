// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package encoder

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONEncoder(t *testing.T) {
	encoder := JSON()

	assert.Equal(t, "application/json", encoder.MimeType())

	w := httptest.NewRecorder()

	assert.NoError(t, encoder.Encode(w, true))

	body, err := ioutil.ReadAll(w.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "true\n", string(body))
}
