// Automatically generated by MockGen. DO NOT EDIT!
// Source: ../command/command.go

package mocks

import (
	"github.com/golang/mock/gomock"
)

// Mock of ClientFactory interface
type MockClientFactory struct {
	ctrl     *gomock.Controller
	recorder *_MockClientFactoryRecorder
}

// Recorder for MockClientFactory (not exported)
type _MockClientFactoryRecorder struct {
	mock *MockClientFactory
}

func NewMockClientFactory(ctrl *gomock.Controller) *MockClientFactory {
	mock := &MockClientFactory{ctrl: ctrl}
	mock.recorder = &_MockClientFactoryRecorder{mock}
	return mock
}

func (_m *MockClientFactory) EXPECT() *_MockClientFactoryRecorder {
	return _m.recorder
}

//func (_m *MockClientFactory) NewClient() proto.UserAuthenticationClient {
//	ret := _m.ctrl.Call(_m, "NewClient")
//	ret0, _ := ret[0].(proto.UserAuthenticationClient)
//	return ret0
//}

func (_mr *_MockClientFactoryRecorder) NewClient() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "NewClient")
}
