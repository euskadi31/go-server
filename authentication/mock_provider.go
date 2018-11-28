// Code generated by mockery v1.0.0. DO NOT EDIT.

package authentication

import http "net/http"
import mock "github.com/stretchr/testify/mock"

// MockProvider is an autogenerated mock type for the Provider type
type MockProvider struct {
	mock.Mock
}

// Validate provides a mock function with given fields: r
func (_m *MockProvider) Validate(r *http.Request) bool {
	ret := _m.Called(r)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*http.Request) bool); ok {
		r0 = rf(r)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}