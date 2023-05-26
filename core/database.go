package core

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

var DB *database

var updateAll = clause.OnConflict{UpdateAll: true}

//go:generate mockery --name Database
type Database interface {
	Upsert(value interface{}) DatabaseResult
	UpsertWithOmit(value interface{}, omitColumns ...string) DatabaseResult
}

type DatabaseResult struct {
	RowsAffected int64
	Error error
}

type database struct {
	config *Config
	gorm *gorm.DB
}

const (
	dsnFormat = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
)

func SetupDatabase(config *Config) {
	db := &database{
		config: config,
	}
	db.connect()
	DB = db
}

func (db *database) Upsert(value interface{}) DatabaseResult {
	return gormReturn(db.gorm.Clauses(updateAll).Create(value))
}

func (db *database) UpsertWithOmit(value interface{}, omitColumns ...string) DatabaseResult {
	return gormReturn(db.gorm.Clauses(updateAll).Omit(omitColumns...).Create(value))
}

func (db *database) connect() {
	Log.WithFields(logrus.Fields{
		"name":     db.config.Database.Name,
		"location": db.config.Database.Location,
		"port":     db.config.Database.Port,
	}).Info("Connecting to database...")
	dsn := fmt.Sprintf(dsnFormat, db.config.Database.User, db.config.Database.Password,
		db.config.Database.Location, db.config.Database.Port, db.config.Database.Name)
	gorm, err := gorm.Open(mysql.Open(dsn), db.getGormConfig())

	if err != nil {
		Log.Fatal(err)
	}

	db.gorm = gorm
}

// simplifies our DAL return signatures to just what we need from GORM
func gormReturn(gorm *gorm.DB) DatabaseResult {
	return DatabaseResult{
		RowsAffected: gorm.RowsAffected,
		Error: gorm.Error,
	}
}

func (db *database) getGormConfig() *gorm.Config {
	var logMode logger.LogLevel

	if db.config.Debug {
		logMode = logger.Info
	} else {
		logMode = logger.Error
	}

	return &gorm.Config{Logger: NewLogger(Log).LogMode(logMode)}
}

type gormLogger struct {
	log   *logrus.Logger
	debug bool
}

func NewLogger(l *logrus.Logger) *gormLogger {
	return &gormLogger{
		log: l,
	}
}

// Implementation of the gorm logger.Interface methods
func (l *gormLogger) LogMode(logLevel logger.LogLevel) logger.Interface {
	if logLevel == logger.Info {
		l.debug = true
	}
	return l
}

func (l *gormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	l.log.WithContext(ctx).Infof(s, args...)
}

func (l *gormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.log.WithContext(ctx).Warnf(s, args...)
}

func (l *gormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	l.log.WithContext(ctx).Errorf(s, args...)
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	fields := logrus.Fields{}
	fields["loc"] = utils.FileWithLineNum()
	fields["rows"] = rows
	fields["ms"] = time.Since(begin)

	if err != nil {
		fields[logrus.ErrorKey] = err
		l.log.WithContext(ctx).WithFields(fields).Errorf(sql)
		return
	}

	if l.debug {
		l.log.WithContext(ctx).WithFields(fields).Debugf(sql)
	}
}