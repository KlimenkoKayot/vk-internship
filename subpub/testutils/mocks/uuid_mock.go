// Code generated by MockGen. DO NOT EDIT.
// Source: ../internal/infrastructure/uuid/domain/interfaces.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUUIDGenerator is a mock of UUIDGenerator interface.
type MockUUIDGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockUUIDGeneratorMockRecorder
}

// MockUUIDGeneratorMockRecorder is the mock recorder for MockUUIDGenerator.
type MockUUIDGeneratorMockRecorder struct {
	mock *MockUUIDGenerator
}

// NewMockUUIDGenerator creates a new mock instance.
func NewMockUUIDGenerator(ctrl *gomock.Controller) *MockUUIDGenerator {
	mock := &MockUUIDGenerator{ctrl: ctrl}
	mock.recorder = &MockUUIDGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUUIDGenerator) EXPECT() *MockUUIDGeneratorMockRecorder {
	return m.recorder
}

// NewString mocks base method.
func (m *MockUUIDGenerator) NewString() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewString")
	ret0, _ := ret[0].(string)
	return ret0
}

// NewString indicates an expected call of NewString.
func (mr *MockUUIDGeneratorMockRecorder) NewString() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewString", reflect.TypeOf((*MockUUIDGenerator)(nil).NewString))
}
