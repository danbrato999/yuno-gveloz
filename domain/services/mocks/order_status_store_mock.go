// Code generated by MockGen. DO NOT EDIT.
// Source: order_status_store.go
//
// Generated by this command:
//
//	mockgen -source=order_status_store.go -destination mocks/order_status_store_mock.go -package mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	domain "github.com/danbrato999/yuno-gveloz/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockOrderStatusStore is a mock of OrderStatusStore interface.
type MockOrderStatusStore struct {
	ctrl     *gomock.Controller
	recorder *MockOrderStatusStoreMockRecorder
	isgomock struct{}
}

// MockOrderStatusStoreMockRecorder is the mock recorder for MockOrderStatusStore.
type MockOrderStatusStoreMockRecorder struct {
	mock *MockOrderStatusStore
}

// NewMockOrderStatusStore creates a new mock instance.
func NewMockOrderStatusStore(ctrl *gomock.Controller) *MockOrderStatusStore {
	mock := &MockOrderStatusStore{ctrl: ctrl}
	mock.recorder = &MockOrderStatusStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderStatusStore) EXPECT() *MockOrderStatusStoreMockRecorder {
	return m.recorder
}

// AddCurrentStatus mocks base method.
func (m *MockOrderStatusStore) AddCurrentStatus(order *domain.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCurrentStatus", order)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCurrentStatus indicates an expected call of AddCurrentStatus.
func (mr *MockOrderStatusStoreMockRecorder) AddCurrentStatus(order any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCurrentStatus", reflect.TypeOf((*MockOrderStatusStore)(nil).AddCurrentStatus), order)
}

// GetHistory mocks base method.
func (m *MockOrderStatusStore) GetHistory(id uint) ([]domain.OrderStatusHistory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHistory", id)
	ret0, _ := ret[0].([]domain.OrderStatusHistory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHistory indicates an expected call of GetHistory.
func (mr *MockOrderStatusStoreMockRecorder) GetHistory(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHistory", reflect.TypeOf((*MockOrderStatusStore)(nil).GetHistory), id)
}
