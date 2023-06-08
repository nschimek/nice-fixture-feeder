package repository

import (
	"github.com/nschimek/nice-fixture-feeder/core"
	"github.com/sirupsen/logrus"
)

type ResultStats struct {
	Success map[string]int
	Error map[string]int
}

func NewResultStats() *ResultStats {
	return &ResultStats{
		Success: make(map[string]int),
		Error: make(map[string]int),
	}
}

func (rs *ResultStats) LogErrors() {
	if len(rs.Error) > 0 {
		core.Log.WithFields(logStatsMap(rs.Error)).Error("Issues during persistence")
	}
}

func (rs *ResultStats) HasErrors() bool {
	if len(rs.Error) > 0 {
		core.Log.WithFields(logStatsMap(rs.Error)).Error("Issues during persistence")
		return true
	}
	return false
}

func (rs *ResultStats) LogSuccesses() {
	if len(rs.Success) > 0 {
		core.Log.WithFields(logStatsMap(rs.Success)).Info("Persistence successful!")
	}
}

func logStatsMap(stats map[string]int) logrus.Fields {
	f := logrus.Fields{}
	for key, value := range stats {
		if value > 0 {
			f[key] = value
		}
	}
	return f
}