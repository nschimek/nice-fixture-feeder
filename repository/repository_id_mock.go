// Code generated by mockery v2.28.0. DO NOT EDIT.

package repository

import mock "github.com/stretchr/testify/mock"

// MockGetByIdRepository is an autogenerated mock type for the GetByIdRepository type
type MockGetByIdRepository[T interface{}, I interface{}] struct {
	mock.Mock
}

type MockGetByIdRepository_Expecter[T interface{}, I interface{}] struct {
	mock *mock.Mock
}

func (_m *MockGetByIdRepository[T, I]) EXPECT() *MockGetByIdRepository_Expecter[T, I] {
	return &MockGetByIdRepository_Expecter[T, I]{mock: &_m.Mock}
}

// GetById provides a mock function with given fields: id
func (_m *MockGetByIdRepository[T, I]) GetById(id I) (*T, error) {
	ret := _m.Called(id)

	var r0 *T
	var r1 error
	if rf, ok := ret.Get(0).(func(I) (*T, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(I) *T); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*T)
		}
	}

	if rf, ok := ret.Get(1).(func(I) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGetByIdRepository_GetById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetById'
type MockGetByIdRepository_GetById_Call[T interface{}, I interface{}] struct {
	*mock.Call
}

// GetById is a helper method to define mock.On call
//   - id I
func (_e *MockGetByIdRepository_Expecter[T, I]) GetById(id interface{}) *MockGetByIdRepository_GetById_Call[T, I] {
	return &MockGetByIdRepository_GetById_Call[T, I]{Call: _e.mock.On("GetById", id)}
}

func (_c *MockGetByIdRepository_GetById_Call[T, I]) Run(run func(id I)) *MockGetByIdRepository_GetById_Call[T, I] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(I))
	})
	return _c
}

func (_c *MockGetByIdRepository_GetById_Call[T, I]) Return(_a0 *T, _a1 error) *MockGetByIdRepository_GetById_Call[T, I] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGetByIdRepository_GetById_Call[T, I]) RunAndReturn(run func(I) (*T, error)) *MockGetByIdRepository_GetById_Call[T, I] {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockGetByIdRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockGetByIdRepository creates a new instance of MockGetByIdRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockGetByIdRepository[T interface{}, I interface{}](t mockConstructorTestingTNewMockGetByIdRepository) *MockGetByIdRepository[T, I] {
	mock := &MockGetByIdRepository[T, I]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
