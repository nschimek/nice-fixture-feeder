// Code generated by mockery v2.28.0. DO NOT EDIT.

package repository

import mock "github.com/stretchr/testify/mock"

// MockUpsertRepository is an autogenerated mock type for the UpsertRepository type
type MockUpsertRepository[T interface{}] struct {
	mock.Mock
}

type MockUpsertRepository_Expecter[T interface{}] struct {
	mock *mock.Mock
}

func (_m *MockUpsertRepository[T]) EXPECT() *MockUpsertRepository_Expecter[T] {
	return &MockUpsertRepository_Expecter[T]{mock: &_m.Mock}
}

// Upsert provides a mock function with given fields: entities
func (_m *MockUpsertRepository[T]) Upsert(entities []T) ([]T, error) {
	ret := _m.Called(entities)

	var r0 []T
	var r1 error
	if rf, ok := ret.Get(0).(func([]T) ([]T, error)); ok {
		return rf(entities)
	}
	if rf, ok := ret.Get(0).(func([]T) []T); ok {
		r0 = rf(entities)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]T)
		}
	}

	if rf, ok := ret.Get(1).(func([]T) error); ok {
		r1 = rf(entities)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUpsertRepository_Upsert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Upsert'
type MockUpsertRepository_Upsert_Call[T interface{}] struct {
	*mock.Call
}

// Upsert is a helper method to define mock.On call
//   - entities []T
func (_e *MockUpsertRepository_Expecter[T]) Upsert(entities interface{}) *MockUpsertRepository_Upsert_Call[T] {
	return &MockUpsertRepository_Upsert_Call[T]{Call: _e.mock.On("Upsert", entities)}
}

func (_c *MockUpsertRepository_Upsert_Call[T]) Run(run func(entities []T)) *MockUpsertRepository_Upsert_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]T))
	})
	return _c
}

func (_c *MockUpsertRepository_Upsert_Call[T]) Return(_a0 []T, _a1 error) *MockUpsertRepository_Upsert_Call[T] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUpsertRepository_Upsert_Call[T]) RunAndReturn(run func([]T) ([]T, error)) *MockUpsertRepository_Upsert_Call[T] {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockUpsertRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockUpsertRepository creates a new instance of MockUpsertRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockUpsertRepository[T interface{}](t mockConstructorTestingTNewMockUpsertRepository) *MockUpsertRepository[T] {
	mock := &MockUpsertRepository[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
