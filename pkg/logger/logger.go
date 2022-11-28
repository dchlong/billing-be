package logger

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type RequestIDContextKey struct {
}

type ContextKey struct {
}

type ILogger interface {
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	NewContext(ctx context.Context) context.Context
	GetLogger(ctx context.Context) ILogger
}

type Logger struct {
	zapLogger *zap.Logger
}

func (l *Logger) Infof(s string, i ...interface{}) {
	l.zapLogger.Sugar().Infof(s, i...)
}

func (l *Logger) Warnf(s string, i ...interface{}) {
	l.zapLogger.Sugar().Warnf(s, i...)
}

func (l *Logger) Errorf(s string, i ...interface{}) {
	l.zapLogger.Sugar().Errorf(s, i...)
}

func (l *Logger) WithContext(ctx context.Context) ILogger {
	reqID, _ := ctx.Value(RequestIDContextKey{}).(string)
	if reqID == "" {
		return l
	}

	return &Logger{
		zapLogger: l.zapLogger.With(zap.String("request_id", reqID)),
	}
}

func (l *Logger) NewContext(ctx context.Context) context.Context {
	reqID, _ := ctx.Value(RequestIDContextKey{}).(string)
	if reqID == "" {
		uid, _ := uuid.NewUUID()
		reqID = uid.String()
	}

	newLogger := &Logger{
		zapLogger: l.zapLogger.With(zap.String("request_id", reqID)),
	}

	return context.WithValue(ctx, ContextKey{}, newLogger)
}

func (l *Logger) GetLogger(ctx context.Context) ILogger {
	clogger, ok := ctx.Value(ContextKey{}).(*Logger)
	if !ok {
		return l
	}

	return clogger
}

func ProvideLogger() (ILogger, func(), error) {
	logger, _ := zap.NewDevelopment()
	return &Logger{
			zapLogger: logger,
		}, func() {
			_ = logger.Sync()
		}, nil
}
