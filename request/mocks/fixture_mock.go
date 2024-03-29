// Code generated by mockery v2.28.0. DO NOT EDIT.

package mocks

import (
	model "github.com/nschimek/nice-fixture-feeder/model"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// Fixture is an autogenerated mock type for the Fixture type
type Fixture struct {
	mock.Mock
}

type Fixture_Expecter struct {
	mock *mock.Mock
}

func (_m *Fixture) EXPECT() *Fixture_Expecter {
	return &Fixture_Expecter{mock: &_m.Mock}
}

// GetIds provides a mock function with given fields:
func (_m *Fixture) GetIds() []int {
	ret := _m.Called()

	var r0 []int
	if rf, ok := ret.Get(0).(func() []int); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int)
		}
	}

	return r0
}

// Fixture_GetIds_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetIds'
type Fixture_GetIds_Call struct {
	*mock.Call
}

// GetIds is a helper method to define mock.On call
func (_e *Fixture_Expecter) GetIds() *Fixture_GetIds_Call {
	return &Fixture_GetIds_Call{Call: _e.mock.On("GetIds")}
}

func (_c *Fixture_GetIds_Call) Run(run func()) *Fixture_GetIds_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Fixture_GetIds_Call) Return(_a0 []int) *Fixture_GetIds_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Fixture_GetIds_Call) RunAndReturn(run func() []int) *Fixture_GetIds_Call {
	_c.Call.Return(run)
	return _c
}

// GetMap provides a mock function with given fields:
func (_m *Fixture) GetMap() map[int]model.Fixture {
	ret := _m.Called()

	var r0 map[int]model.Fixture
	if rf, ok := ret.Get(0).(func() map[int]model.Fixture); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[int]model.Fixture)
		}
	}

	return r0
}

// Fixture_GetMap_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMap'
type Fixture_GetMap_Call struct {
	*mock.Call
}

// GetMap is a helper method to define mock.On call
func (_e *Fixture_Expecter) GetMap() *Fixture_GetMap_Call {
	return &Fixture_GetMap_Call{Call: _e.mock.On("GetMap")}
}

func (_c *Fixture_GetMap_Call) Run(run func()) *Fixture_GetMap_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Fixture_GetMap_Call) Return(_a0 map[int]model.Fixture) *Fixture_GetMap_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Fixture_GetMap_Call) RunAndReturn(run func() map[int]model.Fixture) *Fixture_GetMap_Call {
	_c.Call.Return(run)
	return _c
}

// Persist provides a mock function with given fields:
func (_m *Fixture) Persist() {
	_m.Called()
}

// Fixture_Persist_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Persist'
type Fixture_Persist_Call struct {
	*mock.Call
}

// Persist is a helper method to define mock.On call
func (_e *Fixture_Expecter) Persist() *Fixture_Persist_Call {
	return &Fixture_Persist_Call{Call: _e.mock.On("Persist")}
}

func (_c *Fixture_Persist_Call) Run(run func()) *Fixture_Persist_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Fixture_Persist_Call) Return() *Fixture_Persist_Call {
	_c.Call.Return()
	return _c
}

func (_c *Fixture_Persist_Call) RunAndReturn(run func()) *Fixture_Persist_Call {
	_c.Call.Return(run)
	return _c
}

// Request provides a mock function with given fields:
func (_m *Fixture) Request() {
	_m.Called()
}

// Fixture_Request_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Request'
type Fixture_Request_Call struct {
	*mock.Call
}

// Request is a helper method to define mock.On call
func (_e *Fixture_Expecter) Request() *Fixture_Request_Call {
	return &Fixture_Request_Call{Call: _e.mock.On("Request")}
}

func (_c *Fixture_Request_Call) Run(run func()) *Fixture_Request_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Fixture_Request_Call) Return() *Fixture_Request_Call {
	_c.Call.Return()
	return _c
}

func (_c *Fixture_Request_Call) RunAndReturn(run func()) *Fixture_Request_Call {
	_c.Call.Return(run)
	return _c
}

// RequestDateRange provides a mock function with given fields: startDate, endDate
func (_m *Fixture) RequestDateRange(startDate time.Time, endDate time.Time) {
	_m.Called(startDate, endDate)
}

// Fixture_RequestDateRange_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RequestDateRange'
type Fixture_RequestDateRange_Call struct {
	*mock.Call
}

// RequestDateRange is a helper method to define mock.On call
//   - startDate time.Time
//   - endDate time.Time
func (_e *Fixture_Expecter) RequestDateRange(startDate interface{}, endDate interface{}) *Fixture_RequestDateRange_Call {
	return &Fixture_RequestDateRange_Call{Call: _e.mock.On("RequestDateRange", startDate, endDate)}
}

func (_c *Fixture_RequestDateRange_Call) Run(run func(startDate time.Time, endDate time.Time)) *Fixture_RequestDateRange_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(time.Time), args[1].(time.Time))
	})
	return _c
}

func (_c *Fixture_RequestDateRange_Call) Return() *Fixture_RequestDateRange_Call {
	_c.Call.Return()
	return _c
}

func (_c *Fixture_RequestDateRange_Call) RunAndReturn(run func(time.Time, time.Time)) *Fixture_RequestDateRange_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewFixture interface {
	mock.TestingT
	Cleanup(func())
}

// NewFixture creates a new instance of Fixture. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFixture(t mockConstructorTestingTNewFixture) *Fixture {
	mock := &Fixture{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
