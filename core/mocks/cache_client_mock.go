// Code generated by mockery v2.28.0. DO NOT EDIT.

package mocks

import (
	memcache "github.com/rainycape/memcache"
	mock "github.com/stretchr/testify/mock"
)

// CacheClient is an autogenerated mock type for the CacheClient type
type CacheClient struct {
	mock.Mock
}

type CacheClient_Expecter struct {
	mock *mock.Mock
}

func (_m *CacheClient) EXPECT() *CacheClient_Expecter {
	return &CacheClient_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: key
func (_m *CacheClient) Get(key string) (*memcache.Item, error) {
	ret := _m.Called(key)

	var r0 *memcache.Item
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*memcache.Item, error)); ok {
		return rf(key)
	}
	if rf, ok := ret.Get(0).(func(string) *memcache.Item); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*memcache.Item)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CacheClient_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type CacheClient_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - key string
func (_e *CacheClient_Expecter) Get(key interface{}) *CacheClient_Get_Call {
	return &CacheClient_Get_Call{Call: _e.mock.On("Get", key)}
}

func (_c *CacheClient_Get_Call) Run(run func(key string)) *CacheClient_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *CacheClient_Get_Call) Return(_a0 *memcache.Item, _a1 error) *CacheClient_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CacheClient_Get_Call) RunAndReturn(run func(string) (*memcache.Item, error)) *CacheClient_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Set provides a mock function with given fields: item
func (_m *CacheClient) Set(item *memcache.Item) error {
	ret := _m.Called(item)

	var r0 error
	if rf, ok := ret.Get(0).(func(*memcache.Item) error); ok {
		r0 = rf(item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CacheClient_Set_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Set'
type CacheClient_Set_Call struct {
	*mock.Call
}

// Set is a helper method to define mock.On call
//   - item *memcache.Item
func (_e *CacheClient_Expecter) Set(item interface{}) *CacheClient_Set_Call {
	return &CacheClient_Set_Call{Call: _e.mock.On("Set", item)}
}

func (_c *CacheClient_Set_Call) Run(run func(item *memcache.Item)) *CacheClient_Set_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*memcache.Item))
	})
	return _c
}

func (_c *CacheClient_Set_Call) Return(_a0 error) *CacheClient_Set_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CacheClient_Set_Call) RunAndReturn(run func(*memcache.Item) error) *CacheClient_Set_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewCacheClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewCacheClient creates a new instance of CacheClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCacheClient(t mockConstructorTestingTNewCacheClient) *CacheClient {
	mock := &CacheClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}