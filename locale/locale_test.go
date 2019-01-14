// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package locale

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocaleString(t *testing.T) {
	l := Locale{
		Language: "fr",
		Region:   "FR",
	}

	assert.Equal(t, "fr-FR", l.String())
}
