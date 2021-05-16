// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Logger is an autogenerated mock type for the Logger type
type Logger struct {
	mock.Mock
}

// Debug provides a mock function with given fields: msg, pkg
func (_m *Logger) Debug(msg string, pkg string) {
	_m.Called(msg, pkg)
}

// Error provides a mock function with given fields: msg
func (_m *Logger) Error(msg string) {
	_m.Called(msg)
}

// Info provides a mock function with given fields: msg
func (_m *Logger) Info(msg string) {
	_m.Called(msg)
}

// Warn provides a mock function with given fields: msg
func (_m *Logger) Warn(msg string) {
	_m.Called(msg)
}
