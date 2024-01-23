// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// WisdomQuotesGetter is an autogenerated mock type for the WisdomQuotesGetter type
type WisdomQuotesGetter struct {
	mock.Mock
}

// GetWisdomQuote provides a mock function with given fields:
func (_m *WisdomQuotesGetter) GetWisdomQuote() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

type mockConstructorTestingTNewWisdomQuotesGetter interface {
	mock.TestingT
	Cleanup(func())
}

// NewWisdomQuotesGetter creates a new instance of WisdomQuotesGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewWisdomQuotesGetter(t mockConstructorTestingTNewWisdomQuotesGetter) *WisdomQuotesGetter {
	mock := &WisdomQuotesGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
