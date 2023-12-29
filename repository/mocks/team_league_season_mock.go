// Code generated by mockery v2.28.0. DO NOT EDIT.

package mocks

import (
	model "github.com/nschimek/nice-fixture-feeder/model"
	mock "github.com/stretchr/testify/mock"
)

// TeamLeagueSeason is an autogenerated mock type for the TeamLeagueSeason type
type TeamLeagueSeason struct {
	mock.Mock
}

type TeamLeagueSeason_Expecter struct {
	mock *mock.Mock
}

func (_m *TeamLeagueSeason) EXPECT() *TeamLeagueSeason_Expecter {
	return &TeamLeagueSeason_Expecter{mock: &_m.Mock}
}

// GetById provides a mock function with given fields: id
func (_m *TeamLeagueSeason) GetById(id model.TeamLeagueSeasonId) (*model.TeamLeagueSeason, error) {
	ret := _m.Called(id)

	var r0 *model.TeamLeagueSeason
	var r1 error
	if rf, ok := ret.Get(0).(func(model.TeamLeagueSeasonId) (*model.TeamLeagueSeason, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(model.TeamLeagueSeasonId) *model.TeamLeagueSeason); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.TeamLeagueSeason)
		}
	}

	if rf, ok := ret.Get(1).(func(model.TeamLeagueSeasonId) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TeamLeagueSeason_GetById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetById'
type TeamLeagueSeason_GetById_Call struct {
	*mock.Call
}

// GetById is a helper method to define mock.On call
//   - id model.TeamLeagueSeasonId
func (_e *TeamLeagueSeason_Expecter) GetById(id interface{}) *TeamLeagueSeason_GetById_Call {
	return &TeamLeagueSeason_GetById_Call{Call: _e.mock.On("GetById", id)}
}

func (_c *TeamLeagueSeason_GetById_Call) Run(run func(id model.TeamLeagueSeasonId)) *TeamLeagueSeason_GetById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(model.TeamLeagueSeasonId))
	})
	return _c
}

func (_c *TeamLeagueSeason_GetById_Call) Return(_a0 *model.TeamLeagueSeason, _a1 error) *TeamLeagueSeason_GetById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TeamLeagueSeason_GetById_Call) RunAndReturn(run func(model.TeamLeagueSeasonId) (*model.TeamLeagueSeason, error)) *TeamLeagueSeason_GetById_Call {
	_c.Call.Return(run)
	return _c
}

// Upsert provides a mock function with given fields: entities
func (_m *TeamLeagueSeason) Upsert(entities []model.TeamLeagueSeason) ([]model.TeamLeagueSeason, error) {
	ret := _m.Called(entities)

	var r0 []model.TeamLeagueSeason
	var r1 error
	if rf, ok := ret.Get(0).(func([]model.TeamLeagueSeason) ([]model.TeamLeagueSeason, error)); ok {
		return rf(entities)
	}
	if rf, ok := ret.Get(0).(func([]model.TeamLeagueSeason) []model.TeamLeagueSeason); ok {
		r0 = rf(entities)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.TeamLeagueSeason)
		}
	}

	if rf, ok := ret.Get(1).(func([]model.TeamLeagueSeason) error); ok {
		r1 = rf(entities)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TeamLeagueSeason_Upsert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Upsert'
type TeamLeagueSeason_Upsert_Call struct {
	*mock.Call
}

// Upsert is a helper method to define mock.On call
//   - entities []model.TeamLeagueSeason
func (_e *TeamLeagueSeason_Expecter) Upsert(entities interface{}) *TeamLeagueSeason_Upsert_Call {
	return &TeamLeagueSeason_Upsert_Call{Call: _e.mock.On("Upsert", entities)}
}

func (_c *TeamLeagueSeason_Upsert_Call) Run(run func(entities []model.TeamLeagueSeason)) *TeamLeagueSeason_Upsert_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]model.TeamLeagueSeason))
	})
	return _c
}

func (_c *TeamLeagueSeason_Upsert_Call) Return(_a0 []model.TeamLeagueSeason, _a1 error) *TeamLeagueSeason_Upsert_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TeamLeagueSeason_Upsert_Call) RunAndReturn(run func([]model.TeamLeagueSeason) ([]model.TeamLeagueSeason, error)) *TeamLeagueSeason_Upsert_Call {
	_c.Call.Return(run)
	return _c
}

// UpsertOne provides a mock function with given fields: entity
func (_m *TeamLeagueSeason) UpsertOne(entity model.TeamLeagueSeason) (*model.TeamLeagueSeason, error) {
	ret := _m.Called(entity)

	var r0 *model.TeamLeagueSeason
	var r1 error
	if rf, ok := ret.Get(0).(func(model.TeamLeagueSeason) (*model.TeamLeagueSeason, error)); ok {
		return rf(entity)
	}
	if rf, ok := ret.Get(0).(func(model.TeamLeagueSeason) *model.TeamLeagueSeason); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.TeamLeagueSeason)
		}
	}

	if rf, ok := ret.Get(1).(func(model.TeamLeagueSeason) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TeamLeagueSeason_UpsertOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpsertOne'
type TeamLeagueSeason_UpsertOne_Call struct {
	*mock.Call
}

// UpsertOne is a helper method to define mock.On call
//   - entity model.TeamLeagueSeason
func (_e *TeamLeagueSeason_Expecter) UpsertOne(entity interface{}) *TeamLeagueSeason_UpsertOne_Call {
	return &TeamLeagueSeason_UpsertOne_Call{Call: _e.mock.On("UpsertOne", entity)}
}

func (_c *TeamLeagueSeason_UpsertOne_Call) Run(run func(entity model.TeamLeagueSeason)) *TeamLeagueSeason_UpsertOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(model.TeamLeagueSeason))
	})
	return _c
}

func (_c *TeamLeagueSeason_UpsertOne_Call) Return(_a0 *model.TeamLeagueSeason, _a1 error) *TeamLeagueSeason_UpsertOne_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TeamLeagueSeason_UpsertOne_Call) RunAndReturn(run func(model.TeamLeagueSeason) (*model.TeamLeagueSeason, error)) *TeamLeagueSeason_UpsertOne_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewTeamLeagueSeason interface {
	mock.TestingT
	Cleanup(func())
}

// NewTeamLeagueSeason creates a new instance of TeamLeagueSeason. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTeamLeagueSeason(t mockConstructorTestingTNewTeamLeagueSeason) *TeamLeagueSeason {
	mock := &TeamLeagueSeason{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
