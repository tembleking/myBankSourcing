// Code generated by MockGen. DO NOT EDIT.
// Source: append_only_store.go
//
// Generated by this command:
//
//	mockgen -source=append_only_store.go -destination=mocks/append_only_store.go -package=mocks
//
// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	persistence "github.com/tembleking/myBankSourcing/pkg/persistence"
	gomock "go.uber.org/mock/gomock"
)

// MockAppendOnlyStore is a mock of AppendOnlyStore interface.
type MockAppendOnlyStore struct {
	ctrl     *gomock.Controller
	recorder *MockAppendOnlyStoreMockRecorder
}

// MockAppendOnlyStoreMockRecorder is the mock recorder for MockAppendOnlyStore.
type MockAppendOnlyStoreMockRecorder struct {
	mock *MockAppendOnlyStore
}

// NewMockAppendOnlyStore creates a new mock instance.
func NewMockAppendOnlyStore(ctrl *gomock.Controller) *MockAppendOnlyStore {
	mock := &MockAppendOnlyStore{ctrl: ctrl}
	mock.recorder = &MockAppendOnlyStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAppendOnlyStore) EXPECT() *MockAppendOnlyStoreMockRecorder {
	return m.recorder
}

// AfterEventID mocks base method.
func (m *MockAppendOnlyStore) AfterEventID(eventID string) persistence.ReadOnlyStore {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AfterEventID", eventID)
	ret0, _ := ret[0].(persistence.ReadOnlyStore)
	return ret0
}

// AfterEventID indicates an expected call of AfterEventID.
func (mr *MockAppendOnlyStoreMockRecorder) AfterEventID(eventID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AfterEventID", reflect.TypeOf((*MockAppendOnlyStore)(nil).AfterEventID), eventID)
}

// Append mocks base method.
func (m *MockAppendOnlyStore) Append(ctx context.Context, events ...persistence.StoredStreamEvent) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx}
	for _, a := range events {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Append", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Append indicates an expected call of Append.
func (mr *MockAppendOnlyStoreMockRecorder) Append(ctx any, events ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx}, events...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Append", reflect.TypeOf((*MockAppendOnlyStore)(nil).Append), varargs...)
}

// Limit mocks base method.
func (m *MockAppendOnlyStore) Limit(limit int) persistence.ReadOnlyStore {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Limit", limit)
	ret0, _ := ret[0].(persistence.ReadOnlyStore)
	return ret0
}

// Limit indicates an expected call of Limit.
func (mr *MockAppendOnlyStoreMockRecorder) Limit(limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Limit", reflect.TypeOf((*MockAppendOnlyStore)(nil).Limit), limit)
}

// ReadAllRecords mocks base method.
func (m *MockAppendOnlyStore) ReadAllRecords(ctx context.Context) ([]persistence.StoredStreamEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadAllRecords", ctx)
	ret0, _ := ret[0].([]persistence.StoredStreamEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadAllRecords indicates an expected call of ReadAllRecords.
func (mr *MockAppendOnlyStoreMockRecorder) ReadAllRecords(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadAllRecords", reflect.TypeOf((*MockAppendOnlyStore)(nil).ReadAllRecords), ctx)
}

// ReadRecords mocks base method.
func (m *MockAppendOnlyStore) ReadRecords(ctx context.Context, streamName string) ([]persistence.StoredStreamEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadRecords", ctx, streamName)
	ret0, _ := ret[0].([]persistence.StoredStreamEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadRecords indicates an expected call of ReadRecords.
func (mr *MockAppendOnlyStoreMockRecorder) ReadRecords(ctx, streamName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadRecords", reflect.TypeOf((*MockAppendOnlyStore)(nil).ReadRecords), ctx, streamName)
}

// MockReadOnlyStore is a mock of ReadOnlyStore interface.
type MockReadOnlyStore struct {
	ctrl     *gomock.Controller
	recorder *MockReadOnlyStoreMockRecorder
}

// MockReadOnlyStoreMockRecorder is the mock recorder for MockReadOnlyStore.
type MockReadOnlyStoreMockRecorder struct {
	mock *MockReadOnlyStore
}

// NewMockReadOnlyStore creates a new mock instance.
func NewMockReadOnlyStore(ctrl *gomock.Controller) *MockReadOnlyStore {
	mock := &MockReadOnlyStore{ctrl: ctrl}
	mock.recorder = &MockReadOnlyStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReadOnlyStore) EXPECT() *MockReadOnlyStoreMockRecorder {
	return m.recorder
}

// AfterEventID mocks base method.
func (m *MockReadOnlyStore) AfterEventID(eventID string) persistence.ReadOnlyStore {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AfterEventID", eventID)
	ret0, _ := ret[0].(persistence.ReadOnlyStore)
	return ret0
}

// AfterEventID indicates an expected call of AfterEventID.
func (mr *MockReadOnlyStoreMockRecorder) AfterEventID(eventID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AfterEventID", reflect.TypeOf((*MockReadOnlyStore)(nil).AfterEventID), eventID)
}

// Limit mocks base method.
func (m *MockReadOnlyStore) Limit(limit int) persistence.ReadOnlyStore {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Limit", limit)
	ret0, _ := ret[0].(persistence.ReadOnlyStore)
	return ret0
}

// Limit indicates an expected call of Limit.
func (mr *MockReadOnlyStoreMockRecorder) Limit(limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Limit", reflect.TypeOf((*MockReadOnlyStore)(nil).Limit), limit)
}

// ReadAllRecords mocks base method.
func (m *MockReadOnlyStore) ReadAllRecords(ctx context.Context) ([]persistence.StoredStreamEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadAllRecords", ctx)
	ret0, _ := ret[0].([]persistence.StoredStreamEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadAllRecords indicates an expected call of ReadAllRecords.
func (mr *MockReadOnlyStoreMockRecorder) ReadAllRecords(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadAllRecords", reflect.TypeOf((*MockReadOnlyStore)(nil).ReadAllRecords), ctx)
}

// ReadRecords mocks base method.
func (m *MockReadOnlyStore) ReadRecords(ctx context.Context, streamName string) ([]persistence.StoredStreamEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadRecords", ctx, streamName)
	ret0, _ := ret[0].([]persistence.StoredStreamEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadRecords indicates an expected call of ReadRecords.
func (mr *MockReadOnlyStoreMockRecorder) ReadRecords(ctx, streamName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadRecords", reflect.TypeOf((*MockReadOnlyStore)(nil).ReadRecords), ctx, streamName)
}
