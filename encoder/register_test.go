// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package encoder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {

	assert.Equal(t, 0, len(encoders))

	provider := &MockEncoder{}

	provider.On("MimeType").Return("application/json")

	Register(provider)

	assert.Equal(t, 1, len(encoders))

	encoders = make(map[string]Encoder)
}
