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

	ks, _ := json.Marshal(key)
	vs, _ := json.Marshal(value)

	s.mockCacheClient.EXPECT().Get("tls-" + string(ks)).Return(&memcache.Item{
		Value: vs,
	}, nil)

	res, err := s.cache.Get(key)

	s.Equal(value, *res)
	s.Nil(err)
}

func (s *cacheTestSuite) TestGetMiss() {
	key := model.TeamLeagueSeasonId{TeamId: 1, LeagueId: 37, Season: 2022}
	ks, _ := json.Marshal(key)

	s.mockCacheClient.EXPECT().Get("tls-" + string(ks)).Return(nil, memcache.ErrCacheMiss)

	res, err := s.cache.Get(key)
	
	s.Nil(res)
	s.Nil(err)
}

func (s *cacheTestSuite) TestGetError() {
	key := model.TeamLeagueSeasonId{TeamId: 1, LeagueId: 37, Season: 2022}
	ks, _ := json.Marshal(key)

	s.mockCacheClient.EXPECT().Get("tls-" + string(ks)).Return(nil, memcache.ErrBadMagic)

	res, err := s.cache.Get(key)
	
	s.Nil(res)
	s.Error(err, memcache.ErrBadMagic)
}