// Code generated by mockery v2.16.0. DO NOT EDIT.

package application

import (
	context "context"
	domain "zenport/gates/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// MockApp is an autogenerated mock type for the App type
type MockApp struct {
	mock.Mock
}

// GetTime provides a mock function with given fields: ctx, r
func (_m *MockApp) GetTime(ctx context.Context, r TimeRequest) (*domain.Time, error) {
	ret := _m.Called(ctx, r)

	var r0 *domain.Time
	if rf, ok := ret.Get(0).(func(context.Context, TimeRequest) *domain.Time); ok {
		r0 = rf(ctx, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Time)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, TimeRequest) error); ok {
		r1 = rf(ctx, r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMockApp interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockApp creates a new instance of MockApp. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockApp(t mockConstructorTestingTNewMockApp) *MockApp {
	mock := &MockApp{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
