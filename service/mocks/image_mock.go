// Code generated by mockery v2.28.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Image is an autogenerated mock type for the Image type
type Image struct {
	mock.Mock
}

type Image_Expecter struct {
	mock *mock.Mock
}

func (_m *Image) EXPECT() *Image_Expecter {
	return &Image_Expecter{mock: &_m.Mock}
}

// TransferURL provides a mock function with given fields: url, bucket, keyFormat
func (_m *Image) TransferURL(url string, bucket string, keyFormat string) bool {
	ret := _m.Called(url, bucket, keyFormat)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string, string) bool); ok {
		r0 = rf(url, bucket, keyFormat)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Image_TransferURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TransferURL'
type Image_TransferURL_Call struct {
	*mock.Call
}

// TransferURL is a helper method to define mock.On call
//   - url string
//   - bucket string
//   - keyFormat string
func (_e *Image_Expecter) TransferURL(url interface{}, bucket interface{}, keyFormat interface{}) *Image_TransferURL_Call {
	return &Image_TransferURL_Call{Call: _e.mock.On("TransferURL", url, bucket, keyFormat)}
}

func (_c *Image_TransferURL_Call) Run(run func(url string, bucket string, keyFormat string)) *Image_TransferURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *Image_TransferURL_Call) Return(_a0 bool) *Image_TransferURL_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Image_TransferURL_Call) RunAndReturn(run func(string, string, string) bool) *Image_TransferURL_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewImage interface {
	mock.TestingT
	Cleanup(func())
}

// NewImage creates a new instance of Image. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewImage(t mockConstructorTestingTNewImage) *Image {
	mock := &Image{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
