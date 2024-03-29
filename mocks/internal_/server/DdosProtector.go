// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	pow "github.com/nikvakhrameev/pow_tcp_server/internal/pow"
	mock "github.com/stretchr/testify/mock"
)

// DdosProtector is an autogenerated mock type for the DdosProtector type
type DdosProtector struct {
	mock.Mock
}

// CheckSolution provides a mock function with given fields: challenge, nonce
func (_m *DdosProtector) CheckSolution(challenge pow.Challenge, nonce uint64) (bool, error) {
	ret := _m.Called(challenge, nonce)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(pow.Challenge, uint64) (bool, error)); ok {
		return rf(challenge, nonce)
	}
	if rf, ok := ret.Get(0).(func(pow.Challenge, uint64) bool); ok {
		r0 = rf(challenge, nonce)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(pow.Challenge, uint64) error); ok {
		r1 = rf(challenge, nonce)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateChallenge provides a mock function with given fields:
func (_m *DdosProtector) GenerateChallenge() (pow.Challenge, error) {
	ret := _m.Called()

	var r0 pow.Challenge
	var r1 error
	if rf, ok := ret.Get(0).(func() (pow.Challenge, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() pow.Challenge); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(pow.Challenge)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewDdosProtector interface {
	mock.TestingT
	Cleanup(func())
}

// NewDdosProtector creates a new instance of DdosProtector. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDdosProtector(t mockConstructorTestingTNewDdosProtector) *DdosProtector {
	mock := &DdosProtector{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
