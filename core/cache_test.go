package core_test // use a different package here to avoid cycle errors

import (
	"encoding/json"
	"testing"

	. "github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/core/mocks"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/rainycape/memcache"
	"github.com/stretchr/testify/suite"
)

type cacheTestSuite struct {
	suite.Suite
	mockCacheClient *mocks.CacheClient
	cache Cache[model.TeamLeagueSeason]
}

func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(cacheTestSuite))
}

func (s *cacheTestSuite) SetupTest() {
	s.mockCacheClient = &mocks.CacheClient{}
	s.cache = NewCache[model.TeamLeagueSeason](s.mockCacheClient, "tls")
}

func (s *cacheTestSuite) TestGetHit() {
	key := model.TeamLeagueSeasonId{TeamId: 1, LeagueId: 37, Season: 2022}
	value := model.TeamLeagueSeason{Id: key, MaxFixtureId: 100}

	kb, _ := json.Marshal(key)
	vb, _ := json.Marshal(value)

	s.mockCacheClient.EXPECT().Get("tls-" + string(kb)).Return(&memcache.Item{
		Value: vb,
	}, nil)

	res, err := s.cache.Get(key)

	s.Equal(value, *res)
	s.Nil(err)
}

func (s *cacheTestSuite) TestGetMiss() {
	key := model.TeamLeagueSeasonId{TeamId: 1, LeagueId: 37, Season: 2022}
	kb, _ := json.Marshal(key)

	s.mockCacheClient.EXPECT().Get("tls-" + string(kb)).Return(nil, memcache.ErrCacheMiss)

	res, err := s.cache.Get(key)
	
	s.Nil(res)
	s.Nil(err)
}

func (s *cacheTestSuite) TestGetError() {
	key := model.TeamLeagueSeasonId{TeamId: 1, LeagueId: 37, Season: 2022}
	kb, _ := json.Marshal(key)

	s.mockCacheClient.EXPECT().Get("tls-" + string(kb)).Return(nil, memcache.ErrBadMagic)

	res, err := s.cache.Get(key)
	
	s.Nil(res)
	s.Error(err, memcache.ErrBadMagic)
}

func (s *cacheTestSuite) TestGetErrorUnmarshall() {
	key := model.TeamLeagueSeasonId{TeamId: 1, LeagueId: 37, Season: 2022}
	kb, _ := json.Marshal(key)

	s.mockCacheClient.EXPECT().Get("tls-" + string(kb)).Return(&memcache.Item{
		Value: []byte("invalid json"),
	}, nil)

	res, err := s.cache.Get(key)
	
	s.Nil(res)
	s.Error(err, "Unmarshal")
}

func (s *cacheTestSuite) TestSet() {
	key := model.TeamLeagueSeasonId{TeamId: 1, LeagueId: 37, Season: 2022}
	kb, _ := json.Marshal(key)

	value := model.TeamLeagueSeason{Id: key, MaxFixtureId: 100}
	vb, _ := json.Marshal(value)

	s.mockCacheClient.EXPECT().Set(&memcache.Item{
		Key: "tls-" + string(kb),
		Value: vb,
		Expiration: 900,
	}).Return(nil)

	err := s.cache.Set(key, &value)

	s.Nil(err)
	s.mockCacheClient.AssertExpectations(s.T())
}

func (s *cacheTestSuite) TestSetError() {
	key := model.TeamLeagueSeasonId{TeamId: 1, LeagueId: 37, Season: 2022}
	kb, _ := json.Marshal(key)

	value := model.TeamLeagueSeason{Id: key, MaxFixtureId: 100}
	vb, _ := json.Marshal(value)

	s.mockCacheClient.EXPECT().Set(&memcache.Item{
		Key: "tls-" + string(kb),
		Value: vb,
		Expiration: 900,
	}).Return(memcache.ErrBadMagic)

	err := s.cache.Set(key, &value)

	s.Error(err, memcache.ErrBadMagic)
}