// Code generated by mockery v2.28.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// S3Client is an autogenerated mock type for the S3Client type
type S3Client struct {
	mock.Mock
}

type S3Client_Expecter struct {
	mock *mock.Mock
}

func (_m *S3Client) EXPECT() *S3Client_Expecter {
	return &S3Client_Expecter{mock: &_m.Mock}
}

// Exists provides a mock function with given fields: bucket, key
func (_m *S3Client) Exists(bucket string, key string) (bool, error) {
	ret := _m.Called(bucket, key)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (bool, error)); ok {
		return rf(bucket, key)
	}
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(bucket, key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(bucket, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// S3Client_Exists_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exists'
type S3Client_Exists_Call struct {
	*mock.Call
}

// Exists is a helper method to define mock.On call
//   - bucket string
//   - key string
func (_e *S3Client_Expecter) Exists(bucket interface{}, key interface{}) *S3Client_Exists_Call {
	return &S3Client_Exists_Call{Call: _e.mock.On("Exists", bucket, key)}
}

func (_c *S3Client_Exists_Call) Run(run func(bucket string, key string)) *S3Client_Exists_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *S3Client_Exists_Call) Return(_a0 bool, _a1 error) *S3Client_Exists_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *S3Client_Exists_Call) RunAndReturn(run func(string, string) (bool, error)) *S3Client_Exists_Call {
	_c.Call.Return(run)
	return _c
}

// Upload provides a mock function with given fields: data, bucket, key
func (_m *S3Client) Upload(data []byte, bucket string, key string) error {
	ret := _m.Called(data, bucket, key)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte, string, string) error); ok {
		r0 = rf(data, bucket, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// S3Client_Upload_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Upload'
type S3Client_Upload_Call struct {
	*mock.Call
}

// Upload is a helper method to define mock.On call
//   - data []byte
//   - bucket string
//   - key string
func (_e *S3Client_Expecter) Upload(data interface{}, bucket interface{}, key interface{}) *S3Client_Upload_Call {
	return &S3Client_Upload_Call{Call: _e.mock.On("Upload", data, bucket, key)}
}

func (_c *S3Client_Upload_Call) Run(run func(data []byte, bucket string, key string)) *S3Client_Upload_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *S3Client_Upload_Call) Return(_a0 error) *S3Client_Upload_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *S3Client_Upload_Call) RunAndReturn(run func([]byte, string, string) error) *S3Client_Upload_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewS3Client interface {
	mock.TestingT
	Cleanup(func())
}

// NewS3Client creates a new instance of S3Client. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewS3Client(t mockConstructorTestingTNewS3Client) *S3Client {
	mock := &S3Client{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
