// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/klimenkokayot/vk-internship/libs/logger (interfaces: Logger)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/klimenkokayot/vk-internship/libs/logger/domain"
)

// MockLogger is a mock of Logger interface.
type MockLogger struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerMockRecorder
}

// MockLoggerMockRecorder is the mock recorder for MockLogger.
type MockLoggerMockRecorder struct {
	mock *MockLogger
}

// NewMockLogger creates a new mock instance.
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {
	mock := &MockLogger{ctrl: ctrl}
	mock.recorder = &MockLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder {
	return m.recorder
}

// Debug mocks base method.
func (m *MockLogger) Debug(arg0 string, arg1 ...domain.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debug", varargs...)
}

// Debug indicates an expected call of Debug.
func (mr *MockLoggerMockRecorder) Debug(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*MockLogger)(nil).Debug), varargs...)
}

// Error mocks base method.
func (m *MockLogger) Error(arg0 string, arg1 ...domain.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Error", varargs...)
}

// Error indicates an expected call of Error.
func (mr *MockLoggerMockRecorder) Error(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockLogger)(nil).Error), varargs...)
}

// Fatal mocks base method.
func (m *MockLogger) Fatal(arg0 string, arg1 ...domain.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Fatal", varargs...)
}

// Fatal indicates an expected call of Fatal.
func (mr *MockLoggerMockRecorder) Fatal(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fatal", reflect.TypeOf((*MockLogger)(nil).Fatal), varargs...)
}

// Info mocks base method.
func (m *MockLogger) Info(arg0 string, arg1 ...domain.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Info", varargs...)
}

// Info indicates an expected call of Info.
func (mr *MockLoggerMockRecorder) Info(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockLogger)(nil).Info), varargs...)
}

// OK mocks base method.
func (m *MockLogger) OK(arg0 string, arg1 ...domain.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "OK", varargs...)
}

// OK indicates an expected call of OK.
func (mr *MockLoggerMockRecorder) OK(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OK", reflect.TypeOf((*MockLogger)(nil).OK), varargs...)
}

// Warn mocks base method.
func (m *MockLogger) Warn(arg0 string, arg1 ...domain.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warn", varargs...)
}

// Warn indicates an expected call of Warn.
func (mr *MockLoggerMockRecorder) Warn(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warn", reflect.TypeOf((*MockLogger)(nil).Warn), varargs...)
}

// WithFields mocks base method.
func (m *MockLogger) WithFields(arg0 ...domain.Field) domain.Logger {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "WithFields", varargs...)
	ret0, _ := ret[0].(domain.Logger)
	return ret0
}

// WithFields indicates an expected call of WithFields.
func (mr *MockLoggerMockRecorder) WithFields(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithFields", reflect.TypeOf((*MockLogger)(nil).WithFields), arg0...)
}

// WithLayer mocks base method.
func (m *MockLogger) WithLayer(arg0 string) domain.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithLayer", arg0)
	ret0, _ := ret[0].(domain.Logger)
	return ret0
}

// WithLayer indicates an expected call of WithLayer.
func (mr *MockLoggerMockRecorder) WithLayer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithLayer", reflect.TypeOf((*MockLogger)(nil).WithLayer), arg0)
}
