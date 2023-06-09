// Code generated by mockery v2.28.0. DO NOT EDIT.

package repository

import (
	model "github.com/nschimek/nice-fixture-feeder/model"
	mock "github.com/stretchr/testify/mock"
)

// MockFixtureStatusRepository is an autogenerated mock type for the FixtureStatusRepository type
type MockFixtureStatusRepository struct {
	mock.Mock
}

type MockFixtureStatusRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockFixtureStatusRepository) EXPECT() *MockFixtureStatusRepository_Expecter {
	return &MockFixtureStatusRepository_Expecter{mock: &_m.Mock}
}

// GetAll provides a mock function with given fields:
func (_m *MockFixtureStatusRepository) GetAll() ([]model.FixtureStatus, error) {
	ret := _m.Called()

	var r0 []model.FixtureStatus
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]model.FixtureStatus, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []model.FixtureStatus); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.FixtureStatus)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockFixtureStatusRepository_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type MockFixtureStatusRepository_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
func (_e *MockFixtureStatusRepository_Expecter) GetAll() *MockFixtureStatusRepository_GetAll_Call {
	return &MockFixtureStatusRepository_GetAll_Call{Call: _e.mock.On("GetAll")}
}

func (_c *MockFixtureStatusRepository_GetAll_Call) Run(run func()) *MockFixtureStatusRepository_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockFixtureStatusRepository_GetAll_Call) Return(_a0 []model.FixtureStatus, _a1 error) *MockFixtureStatusRepository_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockFixtureStatusRepository_GetAll_Call) RunAndReturn(run func() ([]model.FixtureStatus, error)) *MockFixtureStatusRepository_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockFixtureStatusRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockFixtureStatusRepository creates a new instance of MockFixtureStatusRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockFixtureStatusRepository(t mockConstructorTestingTNewMockFixtureStatusRepository) *MockFixtureStatusRepository {
	mock := &MockFixtureStatusRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
