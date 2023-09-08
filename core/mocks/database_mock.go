// Code generated by mockery v2.28.0. DO NOT EDIT.

package mocks

import (
	core "github.com/nschimek/nice-fixture-feeder/core"
	mock "github.com/stretchr/testify/mock"
)

// Database is an autogenerated mock type for the Database type
type Database struct {
	mock.Mock
}

type Database_Expecter struct {
	mock *mock.Mock
}

func (_m *Database) EXPECT() *Database_Expecter {
	return &Database_Expecter{mock: &_m.Mock}
}

// GetAll provides a mock function with given fields: entities
func (_m *Database) GetAll(entities interface{}) core.DatabaseResult {
	ret := _m.Called(entities)

	var r0 core.DatabaseResult
	if rf, ok := ret.Get(0).(func(interface{}) core.DatabaseResult); ok {
		r0 = rf(entities)
	} else {
		r0 = ret.Get(0).(core.DatabaseResult)
	}

	return r0
}

// Database_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type Database_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - entities interface{}
func (_e *Database_Expecter) GetAll(entities interface{}) *Database_GetAll_Call {
	return &Database_GetAll_Call{Call: _e.mock.On("GetAll", entities)}
}

func (_c *Database_GetAll_Call) Run(run func(entities interface{})) *Database_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *Database_GetAll_Call) Return(_a0 core.DatabaseResult) *Database_GetAll_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Database_GetAll_Call) RunAndReturn(run func(interface{}) core.DatabaseResult) *Database_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetById provides a mock function with given fields: id, dest
func (_m *Database) GetById(id interface{}, dest interface{}) core.DatabaseResult {
	ret := _m.Called(id, dest)

	var r0 core.DatabaseResult
	if rf, ok := ret.Get(0).(func(interface{}, interface{}) core.DatabaseResult); ok {
		r0 = rf(id, dest)
	} else {
		r0 = ret.Get(0).(core.DatabaseResult)
	}

	return r0
}

// Database_GetById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetById'
type Database_GetById_Call struct {
	*mock.Call
}

// GetById is a helper method to define mock.On call
//   - id interface{}
//   - dest interface{}
func (_e *Database_Expecter) GetById(id interface{}, dest interface{}) *Database_GetById_Call {
	return &Database_GetById_Call{Call: _e.mock.On("GetById", id, dest)}
}

func (_c *Database_GetById_Call) Run(run func(id interface{}, dest interface{})) *Database_GetById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}), args[1].(interface{}))
	})
	return _c
}

func (_c *Database_GetById_Call) Return(_a0 core.DatabaseResult) *Database_GetById_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Database_GetById_Call) RunAndReturn(run func(interface{}, interface{}) core.DatabaseResult) *Database_GetById_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields: value
func (_m *Database) Save(value interface{}) core.DatabaseResult {
	ret := _m.Called(value)

	var r0 core.DatabaseResult
	if rf, ok := ret.Get(0).(func(interface{}) core.DatabaseResult); ok {
		r0 = rf(value)
	} else {
		r0 = ret.Get(0).(core.DatabaseResult)
	}

	return r0
}

// Database_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type Database_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - value interface{}
func (_e *Database_Expecter) Save(value interface{}) *Database_Save_Call {
	return &Database_Save_Call{Call: _e.mock.On("Save", value)}
}

func (_c *Database_Save_Call) Run(run func(value interface{})) *Database_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *Database_Save_Call) Return(_a0 core.DatabaseResult) *Database_Save_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Database_Save_Call) RunAndReturn(run func(interface{}) core.DatabaseResult) *Database_Save_Call {
	_c.Call.Return(run)
	return _c
}

// Upsert provides a mock function with given fields: value
func (_m *Database) Upsert(value interface{}) core.DatabaseResult {
	ret := _m.Called(value)

	var r0 core.DatabaseResult
	if rf, ok := ret.Get(0).(func(interface{}) core.DatabaseResult); ok {
		r0 = rf(value)
	} else {
		r0 = ret.Get(0).(core.DatabaseResult)
	}

	return r0
}

// Database_Upsert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Upsert'
type Database_Upsert_Call struct {
	*mock.Call
}

// Upsert is a helper method to define mock.On call
//   - value interface{}
func (_e *Database_Expecter) Upsert(value interface{}) *Database_Upsert_Call {
	return &Database_Upsert_Call{Call: _e.mock.On("Upsert", value)}
}

func (_c *Database_Upsert_Call) Run(run func(value interface{})) *Database_Upsert_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *Database_Upsert_Call) Return(_a0 core.DatabaseResult) *Database_Upsert_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Database_Upsert_Call) RunAndReturn(run func(interface{}) core.DatabaseResult) *Database_Upsert_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewDatabase interface {
	mock.TestingT
	Cleanup(func())
}

// NewDatabase creates a new instance of Database. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDatabase(t mockConstructorTestingTNewDatabase) *Database {
	mock := &Database{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}