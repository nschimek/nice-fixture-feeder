package request

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/nschimek/nice-fixture-feeder/model"
	"github.com/stretchr/testify/suite"
)

const (
	validResponseJson = `{"get":"leagues","parameters":{"id":"39","season":"2022"},"errors":[],"results":1,"response":[{"league":{"id":39,"name":"Premier League","type":"League","logo":"https:\/\/media-2.api-sports.io\/football\/leagues\/39.png"},"country":{"name":"England","code":"GB","flag":"https:\/\/media-2.api-sports.io\/flags\/gb.svg"},"seasons":[{"year":2022,"start":"2022-08-05","end":"2023-05-28","current":true,"coverage":{"fixtures":{"events":true,"lineups":true,"statistics_fixtures":true,"statistics_players":true},"standings":true,"players":true,"top_scorers":true,"top_assists":true,"top_cards":true,"injuries":true,"predictions":true,"odds":true}}]}]}`
	emptyResponseJson = `{"get":"leagues","parameters":{"id":"39","season":"2022"},"errors":[],"results":1,"response":[]}`
)

type requesterTestSuite struct {
	suite.Suite
	requester *requester[model.League]
}

func TestRequesterTestSuite(t *testing.T) {
	suite.Run(t, new(requesterTestSuite))
}

func (s *requesterTestSuite) SetupTest() {
	s.requester = &requester[model.League]{config: core.MockConfig}
}

func (s *requesterTestSuite) TestValidWithParams() {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Equal("id=39", r.URL.RawQuery)
		w.Write([]byte(validResponseJson))
	}))
	defer svr.Close()
	s.requester.BaseUrl = svr.URL

	r, err := s.requester.Get("test", url.Values{"id": {"39"}})

	s.Nil(err)
	s.Len(r.Response, 1)
	s.Equal(r.Response[0].League.Id, 39)
	s.Len(r.Response[0].Seasons, 1)
	s.Equal(r.Response[0].Seasons[0].Season, 2022)
}

func (s *requesterTestSuite) TestValidWithoutParams() {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Equal("", r.URL.RawQuery)
		w.Write([]byte(validResponseJson))
	}))
	defer svr.Close()
	s.requester.BaseUrl = svr.URL

	r, err := s.requester.Get("test", nil)

	s.Nil(err)
	s.Len(r.Response, 1)
}

func (s *requesterTestSuite) TestHttpError() {
	s.requester.BaseUrl = "invalid"
	r, err := s.requester.Get("test", nil)

	s.Nil(r)
	s.ErrorContains(err, "unsupported protocol scheme")
}

func (s *requesterTestSuite) TestNon200() {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "test", http.StatusNoContent)
	}))
	defer svr.Close()
	s.requester.BaseUrl = svr.URL

	r, err := s.requester.Get("test", nil)

	s.Nil(r)
	s.ErrorContains(err, "received non-200 response code")
}

func (s *requesterTestSuite) TestReadError() {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1")
	}))
	defer svr.Close()
	s.requester.BaseUrl = svr.URL

	r, err := s.requester.Get("test", nil)

	s.Nil(r)
	s.ErrorContains(err, "unexpected EOF")
}

func (s *requesterTestSuite) TestInvalidJson() {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{test}"))
	}))
	defer svr.Close()
	s.requester.BaseUrl = svr.URL

	r, err := s.requester.Get("test", nil)

	s.Nil(r)
	s.ErrorContains(err, "invalid character")
}

func (s *requesterTestSuite) TestEmptyResponse() {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(emptyResponseJson))
	}))
	defer svr.Close()
	s.requester.BaseUrl = svr.URL

	r, err := s.requester.Get("test", nil)

	s.Nil(err)
	s.Len(r.Response, 0)
}
