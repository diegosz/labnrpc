// Code generated by mockery v2.33.2. DO NOT EDIT.

package provisioningpb

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockProvisioningServiceServer is an autogenerated mock type for the ProvisioningServiceServer type
type MockProvisioningServiceServer struct {
	mock.Mock
}

// SayHello provides a mock function with given fields: _a0, _a1
func (_m *MockProvisioningServiceServer) SayHello(_a0 context.Context, _a1 *SayHelloRequest) (*SayHelloResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *SayHelloResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *SayHelloRequest) (*SayHelloResponse, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *SayHelloRequest) *SayHelloResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*SayHelloResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *SayHelloRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mustEmbedUnimplementedProvisioningServiceServer provides a mock function with given fields:
func (_m *MockProvisioningServiceServer) mustEmbedUnimplementedProvisioningServiceServer() {
	_m.Called()
}

// NewMockProvisioningServiceServer creates a new instance of MockProvisioningServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockProvisioningServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockProvisioningServiceServer {
	mock := &MockProvisioningServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
