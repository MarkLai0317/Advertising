// Code generated by mockery v2.40.1. DO NOT EDIT.

package ad

import (
	ad "github.com/MarkLai0317/Advertising/ad"
	mock "github.com/stretchr/testify/mock"
)

// UseCase is an autogenerated mock type for the UseCase type
type UseCase struct {
	mock.Mock
}

type UseCase_Expecter struct {
	mock *mock.Mock
}

func (_m *UseCase) EXPECT() *UseCase_Expecter {
	return &UseCase_Expecter{mock: &_m.Mock}
}

// Advertise provides a mock function with given fields: client
func (_m *UseCase) Advertise(client *ad.Client) ([]ad.Advertisement, error) {
	ret := _m.Called(client)

	if len(ret) == 0 {
		panic("no return value specified for Advertise")
	}

	var r0 []ad.Advertisement
	var r1 error
	if rf, ok := ret.Get(0).(func(*ad.Client) ([]ad.Advertisement, error)); ok {
		return rf(client)
	}
	if rf, ok := ret.Get(0).(func(*ad.Client) []ad.Advertisement); ok {
		r0 = rf(client)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ad.Advertisement)
		}
	}

	if rf, ok := ret.Get(1).(func(*ad.Client) error); ok {
		r1 = rf(client)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UseCase_Advertise_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Advertise'
type UseCase_Advertise_Call struct {
	*mock.Call
}

// Advertise is a helper method to define mock.On call
//   - client *ad.Client
func (_e *UseCase_Expecter) Advertise(client interface{}) *UseCase_Advertise_Call {
	return &UseCase_Advertise_Call{Call: _e.mock.On("Advertise", client)}
}

func (_c *UseCase_Advertise_Call) Run(run func(client *ad.Client)) *UseCase_Advertise_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*ad.Client))
	})
	return _c
}

func (_c *UseCase_Advertise_Call) Return(_a0 []ad.Advertisement, _a1 error) *UseCase_Advertise_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UseCase_Advertise_Call) RunAndReturn(run func(*ad.Client) ([]ad.Advertisement, error)) *UseCase_Advertise_Call {
	_c.Call.Return(run)
	return _c
}

// CreateAd provides a mock function with given fields: _a0
func (_m *UseCase) CreateAd(_a0 *ad.Advertisement) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for CreateAd")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*ad.Advertisement) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UseCase_CreateAd_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateAd'
type UseCase_CreateAd_Call struct {
	*mock.Call
}

// CreateAd is a helper method to define mock.On call
//   - _a0 *ad.Advertisement
func (_e *UseCase_Expecter) CreateAd(_a0 interface{}) *UseCase_CreateAd_Call {
	return &UseCase_CreateAd_Call{Call: _e.mock.On("CreateAd", _a0)}
}

func (_c *UseCase_CreateAd_Call) Run(run func(_a0 *ad.Advertisement)) *UseCase_CreateAd_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*ad.Advertisement))
	})
	return _c
}

func (_c *UseCase_CreateAd_Call) Return(_a0 error) *UseCase_CreateAd_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UseCase_CreateAd_Call) RunAndReturn(run func(*ad.Advertisement) error) *UseCase_CreateAd_Call {
	_c.Call.Return(run)
	return _c
}

// NewUseCase creates a new instance of UseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *UseCase {
	mock := &UseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
