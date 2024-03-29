// Code generated by mockery v2.28.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Team is an autogenerated mock type for the Team type
type Team struct {
	mock.Mock
}

type Team_Expecter struct {
	mock *mock.Mock
}

func (_m *Team) EXPECT() *Team_Expecter {
	return &Team_Expecter{mock: &_m.Mock}
}

// Persist provides a mock function with given fields:
func (_m *Team) Persist() {
	_m.Called()
}

// Team_Persist_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Persist'
type Team_Persist_Call struct {
	*mock.Call
}

// Persist is a helper method to define mock.On call
func (_e *Team_Expecter) Persist() *Team_Persist_Call {
	return &Team_Persist_Call{Call: _e.mock.On("Persist")}
}

func (_c *Team_Persist_Call) Run(run func()) *Team_Persist_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Team_Persist_Call) Return() *Team_Persist_Call {
	_c.Call.Return()
	return _c
}

func (_c *Team_Persist_Call) RunAndReturn(run func()) *Team_Persist_Call {
	_c.Call.Return(run)
	return _c
}

// Request provides a mock function with given fields:
func (_m *Team) Request() {
	_m.Called()
}

// Team_Request_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Request'
type Team_Request_Call struct {
	*mock.Call
}

// Request is a helper method to define mock.On call
func (_e *Team_Expecter) Request() *Team_Request_Call {
	return &Team_Request_Call{Call: _e.mock.On("Request")}
}

func (_c *Team_Request_Call) Run(run func()) *Team_Request_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Team_Request_Call) Return() *Team_Request_Call {
	_c.Call.Return()
	return _c
}

func (_c *Team_Request_Call) RunAndReturn(run func()) *Team_Request_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewTeam interface {
	mock.TestingT
	Cleanup(func())
}

// NewTeam creates a new instance of Team. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTeam(t mockConstructorTestingTNewTeam) *Team {
	mock := &Team{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
