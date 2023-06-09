// Code generated by mockery v2.28.0. DO NOT EDIT.

package core

import mock "github.com/stretchr/testify/mock"

// MockDatabase is an autogenerated mock type for the Database type
type MockDatabase struct {
	mock.Mock
}

type MockDatabase_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDatabase) EXPECT() *MockDatabase_Expecter {
	return &MockDatabase_Expecter{mock: &_m.Mock}
}

// GetAll provides a mock function with given fields: entities
func (_m *MockDatabase) GetAll(entities interface{}) DatabaseResult {
	ret := _m.Called(entities)

	var r0 DatabaseResult
	if rf, ok := ret.Get(0).(func(interface{}) DatabaseResult); ok {
		r0 = rf(entities)
	} else {
		r0 = ret.Get(0).(DatabaseResult)
	}

	return r0
}

// MockDatabase_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type MockDatabase_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - entities interface{}
func (_e *MockDatabase_Expecter) GetAll(entities interface{}) *MockDatabase_GetAll_Call {
	return &MockDatabase_GetAll_Call{Call: _e.mock.On("GetAll", entities)}
}

func (_c *MockDatabase_GetAll_Call) Run(run func(entities interface{})) *MockDatabase_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *MockDatabase_GetAll_Call) Return(_a0 DatabaseResult) *MockDatabase_GetAll_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_GetAll_Call) RunAndReturn(run func(interface{}) DatabaseResult) *MockDatabase_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetById provides a mock function with given fields: id, dest
func (_m *MockDatabase) GetById(id interface{}, dest interface{}) DatabaseResult {
	ret := _m.Called(id, dest)

	var r0 DatabaseResult
	if rf, ok := ret.Get(0).(func(interface{}, interface{}) DatabaseResult); ok {
		r0 = rf(id, dest)
	} else {
		r0 = ret.Get(0).(DatabaseResult)
	}

	return r0
}

// MockDatabase_GetById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetById'
type MockDatabase_GetById_Call struct {
	*mock.Call
}

// GetById is a helper method to define mock.On call
//   - id interface{}
//   - dest interface{}
func (_e *MockDatabase_Expecter) GetById(id interface{}, dest interface{}) *MockDatabase_GetById_Call {
	return &MockDatabase_GetById_Call{Call: _e.mock.On("GetById", id, dest)}
}

func (_c *MockDatabase_GetById_Call) Run(run func(id interface{}, dest interface{})) *MockDatabase_GetById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}), args[1].(interface{}))
	})
	return _c
}

func (_c *MockDatabase_GetById_Call) Return(_a0 DatabaseResult) *MockDatabase_GetById_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_GetById_Call) RunAndReturn(run func(interface{}, interface{}) DatabaseResult) *MockDatabase_GetById_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields: value
func (_m *MockDatabase) Save(value interface{}) DatabaseResult {
	ret := _m.Called(value)

	var r0 DatabaseResult
	if rf, ok := ret.Get(0).(func(interface{}) DatabaseResult); ok {
		r0 = rf(value)
	} else {
		r0 = ret.Get(0).(DatabaseResult)
	}

	return r0
}

// MockDatabase_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type MockDatabase_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - value interface{}
func (_e *MockDatabase_Expecter) Save(value interface{}) *MockDatabase_Save_Call {
	return &MockDatabase_Save_Call{Call: _e.mock.On("Save", value)}
}

func (_c *MockDatabase_Save_Call) Run(run func(value interface{})) *MockDatabase_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *MockDatabase_Save_Call) Return(_a0 DatabaseResult) *MockDatabase_Save_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_Save_Call) RunAndReturn(run func(interface{}) DatabaseResult) *MockDatabase_Save_Call {
	_c.Call.Return(run)
	return _c
}

// Upsert provides a mock function with given fields: value
func (_m *MockDatabase) Upsert(value interface{}) DatabaseResult {
	ret := _m.Called(value)

	var r0 DatabaseResult
	if rf, ok := ret.Get(0).(func(interface{}) DatabaseResult); ok {
		r0 = rf(value)
	} else {
		r0 = ret.Get(0).(DatabaseResult)
	}

	return r0
}

// MockDatabase_Upsert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Upsert'
type MockDatabase_Upsert_Call struct {
	*mock.Call
}

// Upsert is a helper method to define mock.On call
//   - value interface{}
func (_e *MockDatabase_Expecter) Upsert(value interface{}) *MockDatabase_Upsert_Call {
	return &MockDatabase_Upsert_Call{Call: _e.mock.On("Upsert", value)}
}

func (_c *MockDatabase_Upsert_Call) Run(run func(value interface{})) *MockDatabase_Upsert_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *MockDatabase_Upsert_Call) Return(_a0 DatabaseResult) *MockDatabase_Upsert_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_Upsert_Call) RunAndReturn(run func(interface{}) DatabaseResult) *MockDatabase_Upsert_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockDatabase interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockDatabase creates a new instance of MockDatabase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockDatabase(t mockConstructorTestingTNewMockDatabase) *MockDatabase {
	mock := &MockDatabase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
