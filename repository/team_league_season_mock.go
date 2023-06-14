// Code generated by mockery v2.28.0. DO NOT EDIT.

package repository

import (
	model "github.com/nschimek/nice-fixture-feeder/model"
	mock "github.com/stretchr/testify/mock"
)

// MockTeamLeagueSeasonRepository is an autogenerated mock type for the TeamLeagueSeasonRepository type
type MockTeamLeagueSeasonRepository struct {
	mock.Mock
}

type MockTeamLeagueSeasonRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTeamLeagueSeasonRepository) EXPECT() *MockTeamLeagueSeasonRepository_Expecter {
	return &MockTeamLeagueSeasonRepository_Expecter{mock: &_m.Mock}
}

// GetById provides a mock function with given fields: id
func (_m *MockTeamLeagueSeasonRepository) GetById(id model.TeamLeagueSeason) *model.TeamLeagueSeason {
	ret := _m.Called(id)

	var r0 *model.TeamLeagueSeason
	if rf, ok := ret.Get(0).(func(model.TeamLeagueSeason) *model.TeamLeagueSeason); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.TeamLeagueSeason)
		}
	}

	return r0
}

// MockTeamLeagueSeasonRepository_GetById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetById'
type MockTeamLeagueSeasonRepository_GetById_Call struct {
	*mock.Call
}

// GetById is a helper method to define mock.On call
//   - id model.TeamLeagueSeason
func (_e *MockTeamLeagueSeasonRepository_Expecter) GetById(id interface{}) *MockTeamLeagueSeasonRepository_GetById_Call {
	return &MockTeamLeagueSeasonRepository_GetById_Call{Call: _e.mock.On("GetById", id)}
}

func (_c *MockTeamLeagueSeasonRepository_GetById_Call) Run(run func(id model.TeamLeagueSeason)) *MockTeamLeagueSeasonRepository_GetById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(model.TeamLeagueSeason))
	})
	return _c
}

func (_c *MockTeamLeagueSeasonRepository_GetById_Call) Return(_a0 *model.TeamLeagueSeason) *MockTeamLeagueSeasonRepository_GetById_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTeamLeagueSeasonRepository_GetById_Call) RunAndReturn(run func(model.TeamLeagueSeason) *model.TeamLeagueSeason) *MockTeamLeagueSeasonRepository_GetById_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockTeamLeagueSeasonRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockTeamLeagueSeasonRepository creates a new instance of MockTeamLeagueSeasonRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockTeamLeagueSeasonRepository(t mockConstructorTestingTNewMockTeamLeagueSeasonRepository) *MockTeamLeagueSeasonRepository {
	mock := &MockTeamLeagueSeasonRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}