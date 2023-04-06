// Code generated by MockGen. DO NOT EDIT.
// Source: event_store.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	persistence "github.com/tembleking/myBankSourcing/pkg/persistence"
)

// MockEventDispatcher is a mock of EventDispatcher interface.
type MockEventDispatcher struct {
	ctrl     *gomock.Controller
	recorder *MockEventDispatcherMockRecorder
}

// MockEventDispatcherMockRecorder is the mock recorder for MockEventDispatcher.
type MockEventDispatcherMockRecorder struct {
	mock *MockEventDispatcher
}

// NewMockEventDispatcher creates a new mock instance.
func NewMockEventDispatcher(ctrl *gomock.Controller) *MockEventDispatcher {
	mock := &MockEventDispatcher{ctrl: ctrl}
	mock.recorder = &MockEventDispatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventDispatcher) EXPECT() *MockEventDispatcherMockRecorder {
	return m.recorder
}

// Dispatch mocks base method.
func (m *MockEventDispatcher) Dispatch(ctx context.Context, events ...persistence.StreamEvent) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range events {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Dispatch", varargs...)
}

// Dispatch indicates an expected call of Dispatch.
func (mr *MockEventDispatcherMockRecorder) Dispatch(ctx interface{}, events ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, events...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dispatch", reflect.TypeOf((*MockEventDispatcher)(nil).Dispatch), varargs...)
}

// MockClock is a mock of Clock interface.
type MockClock struct {
	ctrl     *gomock.Controller
	recorder *MockClockMockRecorder
}

// MockClockMockRecorder is the mock recorder for MockClock.
type MockClockMockRecorder struct {
	mock *MockClock
}

// NewMockClock creates a new mock instance.
func NewMockClock(ctrl *gomock.Controller) *MockClock {
	mock := &MockClock{ctrl: ctrl}
	mock.recorder = &MockClockMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClock) EXPECT() *MockClockMockRecorder {
	return m.recorder
}

// Now mocks base method.
func (m *MockClock) Now() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Now")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// Now indicates an expected call of Now.
func (mr *MockClockMockRecorder) Now() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Now", reflect.TypeOf((*MockClock)(nil).Now))
}