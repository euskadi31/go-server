// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.
package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	err := ErrorMessage{
		Code:    100,
		Message: "foo",
	}

	assert.Equal(t, 100, err.GetCode())
	assert.Equal(t, "foo", err.GetMessage())
	assert.Equal(t, "foo", err.Error())
}
