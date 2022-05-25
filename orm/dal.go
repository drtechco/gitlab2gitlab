package orm

import (
	"context"
	"drtech.co/gl2gl/core/configs"
	"drtech.co/gl2gl/orm/query"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var ormLogger = logrus.WithField("Name", "ORM")
var _forkDb *gorm.DB

func Setup() error {

	var err error
	_forkDb, err = gorm.Open(sqlite.Open(configs.SqliteDsn),
		&gorm.Config{
			CreateBatchSize: 1000,
			Logger:          &DbLogger{xlogger: ormLogger},
		})
	if err != nil {
		return err
	}
	err = InitDb()
	return err
}

func GetDb() *gorm.DB {
	return _forkDb
}

func InitDb() error {
	return nil
}

func DbQuery() *query.Query {
	return query.Use(GetDb())
}
func MakeQuery(db *gorm.DB) *query.Query {
	return query.Use(db)
}

type DbLogger struct {
	xlogger *logrus.Entry
}

func (l *DbLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}
func (l *DbLogger) Info(ctx context.Context, message string, args ...interface{}) {
	l.xlogger.Info(message, args)
}
func (l *DbLogger) Warn(ctx context.Context, message string, args ...interface{}) {
	l.xlogger.Warn(message, args)
}
func (l *DbLogger) Error(ctx context.Context, message string, args ...interface{}) {
	l.xlogger.Error(message, args)
}
func (l *DbLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	microseconds := time.Now().Sub(begin).Microseconds()
	sql, rowsAffected := fc()
	info := map[string]interface{}{
		"runMilliseconds": float64(microseconds) / float64(1000),
		"begin":           begin.Format("2006-01-02 15:04:05.000000"),
		"sql":             sql,
		"rowsAffected":    rowsAffected,
	}
	if microseconds < 100*1000 {
		l.xlogger.Trace(info)
	} else {
		l.xlogger.Warn(info)
	}

	if err != nil {
		l.xlogger.
			Error(map[string]interface{}{
				"begin": begin.Format("2006-01-02 15:04:05.000000"),
			}, err)
	}

}
