// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Hasher is an autogenerated mock type for the Hasher type
type Hasher struct {
	mock.Mock
}

// HashData provides a mock function with given fields: data
func (_m *Hasher) HashData(data []byte) []byte {
	ret := _m.Called(data)

	var r0 []byte
	if rf, ok := ret.Get(0).(func([]byte) []byte); ok {
		r0 = rf(data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

type mockConstructorTestingTNewHasher interface {
	mock.TestingT
	Cleanup(func())
}

// NewHasher creates a new instance of Hasher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewHasher(t mockConstructorTestingTNewHasher) *Hasher {
	mock := &Hasher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
