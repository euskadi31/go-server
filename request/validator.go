// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package request

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
	"github.com/go-yaml/yaml"
	"github.com/rs/zerolog/log"
)

// ioReader interface for testing
// Hack for generate mock
//go:generate mockery -case=underscore -inpkg -name=ioReader
// nolint: deadcode,megacheck
type ioReader interface {
	io.Reader
}

// ErrSchemaFileFormatNotSupported type
type ErrSchemaFileFormatNotSupported struct {
	Ext string
}

// NewErrSchemaFileFormatNotSupported error
func NewErrSchemaFileFormatNotSupported(ext string) error {
	return &ErrSchemaFileFormatNotSupported{
		Ext: ext,
	}
}

func (e *ErrSchemaFileFormatNotSupported) Error() string {
	return fmt.Sprintf("%s file schema is not supported", e.Ext)
}

// ErrSchemaNotFound type
type ErrSchemaNotFound struct {
	Name string
}

// NewErrSchemaNotFound error
func NewErrSchemaNotFound(name string) error {
	return &ErrSchemaNotFound{
		Name: name,
	}
}

func (e *ErrSchemaNotFound) Error() string {
	return fmt.Sprintf(`schema "%s" not found`, e.Name)
}

// SchemaValidator interface experimental
//go:generate mockery -case=underscore -inpkg -name=SchemaValidator
type SchemaValidator interface {
	SchemaValidator() *spec.Schema
}

// Validator struct
type Validator struct {
	schemas map[string]*validate.SchemaValidator
}

// NewValidator constructor
func NewValidator() *Validator {
	return &Validator{
		schemas: make(map[string]*validate.SchemaValidator),
	}
}

// AddSchemaFromFile name
func (v *Validator) AddSchemaFromFile(name string, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Error().Err(err).Msgf("Close %s file", filename)
		}
	}()

	return v.AddSchemaFromReader(name, strings.Trim(filepath.Ext(filename), "."), file)
}

// AddSchemaFromReader func
func (v *Validator) AddSchemaFromReader(name string, format string, reader io.Reader) error {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	switch format {
	case "json":
		return v.AddSchemaFromJSON(name, b)
	case "yml", "yaml":
		return v.AddSchemaFromYAML(name, b)
	default:
		return NewErrSchemaFileFormatNotSupported(format)
	}
}

// AddSchemaFromJSON string
func (v *Validator) AddSchemaFromJSON(name string, content json.RawMessage) error {
	var schema spec.Schema

	if err := json.Unmarshal(content, &schema); err != nil {
		return err
	}

	return v.AddSchema(name, &schema)
}

// AddSchemaFromYAML string
func (v *Validator) AddSchemaFromYAML(name string, content []byte) error {
	var schema spec.Schema

	if err := yaml.Unmarshal(content, &schema); err != nil {
		return err
	}

	return v.AddSchema(name, &schema)
}

// AddSchema by name
func (v *Validator) AddSchema(name string, schema *spec.Schema) error {
	validator := validate.NewSchemaValidator(schema, nil, "", strfmt.Default)

	v.schemas[name] = validator

	return nil
}

// AddSchemFromObject experimental
func (v *Validator) AddSchemFromObject(object SchemaValidator) error {
	rt := reflect.TypeOf(object)

	return v.AddSchemFromObjectName(rt.Name(), object)
}

// AddSchemFromObjectName experimental
func (v *Validator) AddSchemFromObjectName(name string, object SchemaValidator) error {
	return v.AddSchema(name, object.SchemaValidator())
}

// Validate data
func (v Validator) Validate(name string, data interface{}) *validate.Result {
	result := &validate.Result{}

	validator, ok := v.schemas[name]
	if !ok {
		result.AddErrors(NewErrSchemaNotFound(name))

		return result
	}

	result = validator.Validate(data)

	return result
}
