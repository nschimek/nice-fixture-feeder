// Code generated by mockery v2.28.0. DO NOT EDIT.

package mocks

import (
	model "github.com/nschimek/nice-fixture-feeder/model"
	scores "github.com/nschimek/nice-fixture-feeder/service/scores"
	mock "github.com/stretchr/testify/mock"
)

// Score is an autogenerated mock type for the Score type
type Score struct {
	mock.Mock
}

type Score_Expecter struct {
	mock *mock.Mock
}

func (_m *Score) EXPECT() *Score_Expecter {
	return &Score_Expecter{mock: &_m.Mock}
}

// CanScore provides a mock function with given fields: fixture
func (_m *Score) CanScore(fixture *model.Fixture) bool {
	ret := _m.Called(fixture)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*model.Fixture) bool); ok {
		r0 = rf(fixture)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Score_CanScore_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CanScore'
type Score_CanScore_Call struct {
	*mock.Call
}

// CanScore is a helper method to define mock.On call
//   - fixture *model.Fixture
func (_e *Score_Expecter) CanScore(fixture interface{}) *Score_CanScore_Call {
	return &Score_CanScore_Call{Call: _e.mock.On("CanScore", fixture)}
}

func (_c *Score_CanScore_Call) Run(run func(fixture *model.Fixture)) *Score_CanScore_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*model.Fixture))
	})
	return _c
}

func (_c *Score_CanScore_Call) Return(_a0 bool) *Score_CanScore_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Score_CanScore_Call) RunAndReturn(run func(*model.Fixture) bool) *Score_CanScore_Call {
	_c.Call.Return(run)
	return _c
}

// GetId provides a mock function with given fields:
func (_m *Score) GetId() scores.ScoreId {
	ret := _m.Called()

	var r0 scores.ScoreId
	if rf, ok := ret.Get(0).(func() scores.ScoreId); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(scores.ScoreId)
	}

	return r0
}

// Score_GetId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetId'
type Score_GetId_Call struct {
	*mock.Call
}

// GetId is a helper method to define mock.On call
func (_e *Score_Expecter) GetId() *Score_GetId_Call {
	return &Score_GetId_Call{Call: _e.mock.On("GetId")}
}

func (_c *Score_GetId_Call) Run(run func()) *Score_GetId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Score_GetId_Call) Return(_a0 scores.ScoreId) *Score_GetId_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Score_GetId_Call) RunAndReturn(run func() scores.ScoreId) *Score_GetId_Call {
	_c.Call.Return(run)
	return _c
}

// Score provides a mock function with given fields: fixture
func (_m *Score) Score(fixture *model.Fixture) (*model.FixtureScore, error) {
	ret := _m.Called(fixture)

	var r0 *model.FixtureScore
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.Fixture) (*model.FixtureScore, error)); ok {
		return rf(fixture)
	}
	if rf, ok := ret.Get(0).(func(*model.Fixture) *model.FixtureScore); ok {
		r0 = rf(fixture)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FixtureScore)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.Fixture) error); ok {
		r1 = rf(fixture)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Score_Score_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Score'
type Score_Score_Call struct {
	*mock.Call
}

// Score is a helper method to define mock.On call
//   - fixture *model.Fixture
func (_e *Score_Expecter) Score(fixture interface{}) *Score_Score_Call {
	return &Score_Score_Call{Call: _e.mock.On("Score", fixture)}
}

func (_c *Score_Score_Call) Run(run func(fixture *model.Fixture)) *Score_Score_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*model.Fixture))
	})
	return _c
}

func (_c *Score_Score_Call) Return(_a0 *model.FixtureScore, _a1 error) *Score_Score_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Score_Score_Call) RunAndReturn(run func(*model.Fixture) (*model.FixtureScore, error)) *Score_Score_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewScore interface {
	mock.TestingT
	Cleanup(func())
}

// NewScore creates a new instance of Score. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewScore(t mockConstructorTestingTNewScore) *Score {
	mock := &Score{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
