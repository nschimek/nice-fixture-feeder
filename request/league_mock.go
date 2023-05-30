// Code generated by mockery v2.28.0. DO NOT EDIT.

package request

import mock "github.com/stretchr/testify/mock"

// MockLeagueRequest is an autogenerated mock type for the LeagueRequest type
type MockLeagueRequest struct {
	mock.Mock
}

type MockLeagueRequest_Expecter struct {
	mock *mock.Mock
}

func (_m *MockLeagueRequest) EXPECT() *MockLeagueRequest_Expecter {
	return &MockLeagueRequest_Expecter{mock: &_m.Mock}
}

// Persist provides a mock function with given fields:
func (_m *MockLeagueRequest) Persist() {
	_m.Called()
}

// MockLeagueRequest_Persist_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Persist'
type MockLeagueRequest_Persist_Call struct {
	*mock.Call
}

// Persist is a helper method to define mock.On call
func (_e *MockLeagueRequest_Expecter) Persist() *MockLeagueRequest_Persist_Call {
	return &MockLeagueRequest_Persist_Call{Call: _e.mock.On("Persist")}
}

func (_c *MockLeagueRequest_Persist_Call) Run(run func()) *MockLeagueRequest_Persist_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockLeagueRequest_Persist_Call) Return() *MockLeagueRequest_Persist_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockLeagueRequest_Persist_Call) RunAndReturn(run func()) *MockLeagueRequest_Persist_Call {
	_c.Call.Return(run)
	return _c
}

// PostPersist provides a mock function with given fields:
func (_m *MockLeagueRequest) PostPersist() {
	_m.Called()
}

// MockLeagueRequest_PostPersist_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PostPersist'
type MockLeagueRequest_PostPersist_Call struct {
	*mock.Call
}

// PostPersist is a helper method to define mock.On call
func (_e *MockLeagueRequest_Expecter) PostPersist() *MockLeagueRequest_PostPersist_Call {
	return &MockLeagueRequest_PostPersist_Call{Call: _e.mock.On("PostPersist")}
}

func (_c *MockLeagueRequest_PostPersist_Call) Run(run func()) *MockLeagueRequest_PostPersist_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockLeagueRequest_PostPersist_Call) Return() *MockLeagueRequest_PostPersist_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockLeagueRequest_PostPersist_Call) RunAndReturn(run func()) *MockLeagueRequest_PostPersist_Call {
	_c.Call.Return(run)
	return _c
}

// Request provides a mock function with given fields: idMap
func (_m *MockLeagueRequest) Request(idMap map[string]struct{}) {
	_m.Called(idMap)
}

// MockLeagueRequest_Request_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Request'
type MockLeagueRequest_Request_Call struct {
	*mock.Call
}

// Request is a helper method to define mock.On call
//   - idMap map[string]struct{}
func (_e *MockLeagueRequest_Expecter) Request(idMap interface{}) *MockLeagueRequest_Request_Call {
	return &MockLeagueRequest_Request_Call{Call: _e.mock.On("Request", idMap)}
}

func (_c *MockLeagueRequest_Request_Call) Run(run func(idMap map[string]struct{})) *MockLeagueRequest_Request_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(map[string]struct{}))
	})
	return _c
}

func (_c *MockLeagueRequest_Request_Call) Return() *MockLeagueRequest_Request_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockLeagueRequest_Request_Call) RunAndReturn(run func(map[string]struct{})) *MockLeagueRequest_Request_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockLeagueRequest interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockLeagueRequest creates a new instance of MockLeagueRequest. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockLeagueRequest(t mockConstructorTestingTNewMockLeagueRequest) *MockLeagueRequest {
	mock := &MockLeagueRequest{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}