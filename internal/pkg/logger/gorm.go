package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

func NewGormLogger() GormLogger {
	return GormLogger{
		LogLevel:      logger.Info,
		SlowThreshold: 100 * time.Millisecond,
	}
}

type GormLogger struct {
	LogLevel      logger.LogLevel
	SlowThreshold time.Duration
}

func (l GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return GormLogger{
		SlowThreshold: l.SlowThreshold,
		LogLevel:      level,
	}
}

func (l GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	Logger(ctx).Sugar().Infof(str, args...)
}

func (l GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	Logger(ctx).Sugar().Warnf(str, args...)
}

func (l GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	Logger(ctx).Sugar().Errorf(str, args...)
}

func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	switch {
	case l.LogLevel >= logger.Error && err != nil:
		latency := time.Since(begin)
		sql, rows := fc()
		Logger(ctx).Error("trace", zap.Error(err), zap.Duration("latency", latency), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.LogLevel >= logger.Warn && l.SlowThreshold != 0:
		latency := time.Since(begin)
		if latency > l.SlowThreshold {
			sql, rows := fc()
			Logger(ctx).Warn("trace", zap.Duration("latency", latency), zap.Int64("rows", rows), zap.String("sql", sql))
		}
	case l.LogLevel >= logger.Info:
		latency := time.Since(begin)
		sql, rows := fc()
		Logger(ctx).Debug("trace", zap.Duration("latency", latency), zap.Int64("rows", rows), zap.String("sql", sql))
	}
}
