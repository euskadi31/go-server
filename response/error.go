// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package response

// ErrorMessageInterface interface
type ErrorMessageInterface interface {
	GetCode() int
	GetMessage() string
	Error() string
}

// ErrorMessage struct
type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error implements ErrorMessageInterface
func (e ErrorMessage) Error() string {
	return e.Message
}

// GetCode implements ErrorMessageInterface
func (e ErrorMessage) GetCode() int {
	return e.Code
}

// GetMessage implements ErrorMessageInterface
func (e ErrorMessage) GetMessage() string {
	return e.Message
}

// ValidatorError struct
type ValidatorError struct {
	Code    int32         `json:"code,omitempty"`
	Name    string        `json:"name,omitempty"`
	In      string        `json:"in,omitempty"`
	Value   interface{}   `json:"value,omitempty"`
	Message string        `json:"message"`
	Values  []interface{} `json:"values,omitempty"`
}

func (e ValidatorError) Error() string {
	return e.Message
}

// ErrorResponse struct
type ErrorResponse struct {
	Error ErrorMessage `json:"error"`
}

// ErrorsResponse struct
type ErrorsResponse struct {
	Errors []error `json:"errors"`
}
