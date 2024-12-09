// Code generated by MockGen. DO NOT EDIT.
// Source: contract.go
//
// Generated by this command:
//
//	mockgen -source contract.go -destination mocks_test.go -package create_test
//

// Package create_test is a generated GoMock package.
package create_test

import (
	model "contacts/internal/model"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// Mockstorage is a mock of storage interface.
type Mockstorage struct {
	ctrl     *gomock.Controller
	recorder *MockstorageMockRecorder
}

// MockstorageMockRecorder is the mock recorder for Mockstorage.
type MockstorageMockRecorder struct {
	mock *Mockstorage
}

// NewMockstorage creates a new mock instance.
func NewMockstorage(ctrl *gomock.Controller) *Mockstorage {
	mock := &Mockstorage{ctrl: ctrl}
	mock.recorder = &MockstorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockstorage) EXPECT() *MockstorageMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *Mockstorage) Create(contact model.Contact) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", contact)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockstorageMockRecorder) Create(contact any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*Mockstorage)(nil).Create), contact)
}

// Mockvalidator is a mock of validator interface.
type Mockvalidator struct {
	ctrl     *gomock.Controller
	recorder *MockvalidatorMockRecorder
}

// MockvalidatorMockRecorder is the mock recorder for Mockvalidator.
type MockvalidatorMockRecorder struct {
	mock *Mockvalidator
}

// NewMockvalidator creates a new mock instance.
func NewMockvalidator(ctrl *gomock.Controller) *Mockvalidator {
	mock := &Mockvalidator{ctrl: ctrl}
	mock.recorder = &MockvalidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockvalidator) EXPECT() *MockvalidatorMockRecorder {
	return m.recorder
}

// Validate mocks base method.
func (m *Mockvalidator) Validate(contact model.ContactForCreate) map[model.Field]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", contact)
	ret0, _ := ret[0].(map[model.Field]string)
	return ret0
}

// Validate indicates an expected call of Validate.
func (mr *MockvalidatorMockRecorder) Validate(contact any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*Mockvalidator)(nil).Validate), contact)
}

// Mockuuid is a mock of uuid interface.
type Mockuuid struct {
	ctrl     *gomock.Controller
	recorder *MockuuidMockRecorder
}

// MockuuidMockRecorder is the mock recorder for Mockuuid.
type MockuuidMockRecorder struct {
	mock *Mockuuid
}

// NewMockuuid creates a new mock instance.
func NewMockuuid(ctrl *gomock.Controller) *Mockuuid {
	mock := &Mockuuid{ctrl: ctrl}
	mock.recorder = &MockuuidMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockuuid) EXPECT() *MockuuidMockRecorder {
	return m.recorder
}

// NewString mocks base method.
func (m *Mockuuid) NewString() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewString")
	ret0, _ := ret[0].(string)
	return ret0
}

// NewString indicates an expected call of NewString.
func (mr *MockuuidMockRecorder) NewString() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewString", reflect.TypeOf((*Mockuuid)(nil).NewString))
}