// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package request

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestErrSchemaFileFormatNotSupported(t *testing.T) {
	err := NewErrSchemaFileFormatNotSupported("yml")

	assert.EqualError(t, err, "yml file schema is not supported")
}

func TestErrSchemaNotFound(t *testing.T) {
	err := NewErrSchemaNotFound("foo")

	assert.EqualError(t, err, `schema "foo" not found`)
}

func TestValidatorWithNoSchema(t *testing.T) {
	v := NewValidator()

	result := v.Validate("foo", nil)

	assert.Equal(t, 1, len(result.Errors))
	assert.EqualError(t, result.Errors[0], `schema "foo" not found`)
	assert.False(t, result.IsValid())
}

func TestValidatorAddSchemaFromFileFailed(t *testing.T) {
	v := NewValidator()

	err := v.AddSchemaFromFile("test", "testdata/bad")
	assert.Error(t, err)
}

func TestValidatorAddSchemaFromFileJSON(t *testing.T) {
	v := NewValidator()

	err := v.AddSchemaFromFile("test", "testdata/test.json")
	assert.NoError(t, err)

	req := map[string]interface{}{
		"name": "foo",
	}

	result := v.Validate("test", req)
	assert.True(t, result.IsValid())
}

func TestValidatorAddSchemaFromFileYML(t *testing.T) {
	v := NewValidator()

	err := v.AddSchemaFromFile("test", "testdata/test.yml")
	assert.NoError(t, err)

	req := map[string]interface{}{
		"name": "foo",
	}

	result := v.Validate("test", req)
	assert.True(t, result.IsValid())
}

func TestValidatorAddSchemaFromFileWithBadExt(t *testing.T) {
	v := NewValidator()

	err := v.AddSchemaFromFile("test", "testdata/test.bad")
	assert.Error(t, err)
}

func TestValidatorAddSchemaFromReaderWithBadInput(t *testing.T) {
	r := &mockIoReader{}

	r.On("Read", mock.Anything).Return(0, errors.New("fail"))

	v := NewValidator()

	err := v.AddSchemaFromReader("test", "json", r)
	assert.Error(t, err)
}

func TestValidatorValidateWithBadInput(t *testing.T) {
	v := NewValidator()

	err := v.AddSchemaFromFile("test", "testdata/test.json")
	assert.NoError(t, err)

	req := map[string]interface{}{
		"name": "11",
	}

	result := v.Validate("test", req)
	assert.False(t, result.IsValid())
}

func TestValidatorAddSchemaFromJSONWithBadJSON(t *testing.T) {
	v := NewValidator()

	err := v.AddSchemaFromJSON("test", json.RawMessage(`{`))
	assert.Error(t, err)
}

func TestValidatorAddSchemaFromYMLWithBadYML(t *testing.T) {
	v := NewValidator()

	err := v.AddSchemaFromYAML("test", []byte(`{`))
	assert.Error(t, err)
}

func TestValidatorAddSchemFromObject(t *testing.T) {
	schemaMock := &MockSchemaValidator{}

	schemaMock.On("SchemaValidator").Return(&spec.Schema{})

	v := NewValidator()

	err := v.AddSchemFromObject(schemaMock)
	assert.NoError(t, err)
}

func BenchmarkValidator(b *testing.B) {
	v := NewValidator()

	err := v.AddSchemaFromFile("test", "testdata/test.yml")
	assert.NoError(b, err)

	req := map[string]interface{}{
		"name": "foo",
	}

	for n := 0; n < b.N; n++ {
		_ = v.Validate("test", req)
	}
}
