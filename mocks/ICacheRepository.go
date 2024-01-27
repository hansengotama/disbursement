// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	cacherepo "github.com/hansengotama/disbursement/internal/repository/cache"
	mock "github.com/stretchr/testify/mock"
)

// ICacheRepository is an autogenerated mock type for the ICacheRepository type
type ICacheRepository struct {
	mock.Mock
}

// Get provides a mock function with given fields: param
func (_m *ICacheRepository) Get(param cacherepo.GetParam) (string, error) {
	ret := _m.Called(param)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(cacherepo.GetParam) (string, error)); ok {
		return rf(param)
	}
	if rf, ok := ret.Get(0).(func(cacherepo.GetParam) string); ok {
		r0 = rf(param)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(cacherepo.GetParam) error); ok {
		r1 = rf(param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Set provides a mock function with given fields: param
func (_m *ICacheRepository) Set(param cacherepo.SetParam) (string, error) {
	ret := _m.Called(param)

	if len(ret) == 0 {
		panic("no return value specified for Set")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(cacherepo.SetParam) (string, error)); ok {
		return rf(param)
	}
	if rf, ok := ret.Get(0).(func(cacherepo.SetParam) string); ok {
		r0 = rf(param)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(cacherepo.SetParam) error); ok {
		r1 = rf(param)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewICacheRepository creates a new instance of ICacheRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewICacheRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ICacheRepository {
	mock := &ICacheRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
