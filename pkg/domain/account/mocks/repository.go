// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	account "github.com/tembleking/myBankSourcing/pkg/domain/account"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// GetAccount mocks base method.
func (m *MockRepository) GetAccount(ctx context.Context, accountID account.ID) (*account.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", ctx, accountID)
	ret0, _ := ret[0].(*account.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockRepositoryMockRecorder) GetAccount(ctx, accountID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockRepository)(nil).GetAccount), ctx, accountID)
}

// ListAccounts mocks base method.
func (m *MockRepository) ListAccounts(ctx context.Context) ([]*account.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAccounts", ctx)
	ret0, _ := ret[0].([]*account.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAccounts indicates an expected call of ListAccounts.
func (mr *MockRepositoryMockRecorder) ListAccounts(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAccounts", reflect.TypeOf((*MockRepository)(nil).ListAccounts), ctx)
}

// SaveAccount mocks base method.
func (m *MockRepository) SaveAccount(ctx context.Context, account *account.Account) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveAccount", ctx, account)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveAccount indicates an expected call of SaveAccount.
func (mr *MockRepositoryMockRecorder) SaveAccount(ctx, account interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveAccount", reflect.TypeOf((*MockRepository)(nil).SaveAccount), ctx, account)
}
