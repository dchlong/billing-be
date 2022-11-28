package logger

import (
	"context"
	"time"

	gormlogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	logger        ILogger
	LogLevel      gormlogger.LogLevel
	SlowThreshold time.Duration
}

func NewGormLogger(logger ILogger) GormLogger {
	return GormLogger{
		logger:        logger,
		LogLevel:      gormlogger.Warn,
		SlowThreshold: 100 * time.Millisecond,
	}
}

func (l GormLogger) SetAsDefault() {
	gormlogger.Default = l
}

func (l GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return GormLogger{
		SlowThreshold: l.SlowThreshold,
		LogLevel:      level,
	}
}

func (l GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Info {
		return
	}

	l.logger.GetLogger(ctx).Infof(str, args...)
}

func (l GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Warn {
		return
	}

	l.logger.GetLogger(ctx).Warnf(str, args...)
}

func (l GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Error {
		return
	}

	l.logger.GetLogger(ctx).Errorf(str, args...)
}

func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}

	sql, _ := fc()
	l.logger.GetLogger(ctx).Infof("sql_query: %s", sql)
}
