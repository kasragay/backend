package utils

import (
	"context"
	"errors"
	"fmt"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

type Logger struct {
	*otelzap.Logger
	Fields []zap.Field
}

func NewLogger(customFields ...zap.Field) *Logger {
	config := zap.NewProductionConfig()
	config.DisableStacktrace = true
	z, err := config.Build()
	if err != nil {
		panic(err)
	}
	return &Logger{
		Logger: otelzap.New(z),
		Fields: customFields,
	}
}

func mergeFields(source1, source2 []zap.Field) []zap.Field {
	fields := make([]zap.Field, 0, len(source1)+len(source2))
	fields = append(fields, source1...)
	fields = append(fields, source2...)
	return fields
}

func (l Logger) Info(ctx context.Context, msg string) {
	l.Logger.Ctx(ctx).Info(msg, l.Fields...)
}

func (l Logger) Infof(ctx context.Context, msg string, args ...any) {
	l.Logger.Ctx(ctx).Info(fmt.Sprintf(msg, args...), l.Fields...)
}

func (l Logger) InfoFields(ctx context.Context, msg string, fields ...zap.Field) {
	l.Logger.Ctx(ctx).Info(msg, mergeFields(l.Fields, fields)...)
}

func (l Logger) Trace(ctx context.Context, msg string) {
	l.Logger.Ctx(ctx).Debug(msg, l.Fields...)
}

func (l Logger) Tracef(ctx context.Context, msg string, args ...any) {
	l.Logger.Ctx(ctx).Debug(fmt.Sprintf(msg, args...), l.Fields...)
}

func (l Logger) Debug(ctx context.Context, msg string) {
	l.Logger.Ctx(ctx).Debug(msg, l.Fields...)
}

func (l Logger) Debugf(ctx context.Context, msg string, args ...any) {
	l.Logger.Ctx(ctx).Debug(fmt.Sprintf(msg, args...), l.Fields...)
}

func (l Logger) Warn(ctx context.Context, msg string) {
	l.Logger.Ctx(ctx).Warn(msg, l.Fields...)
}

func (l Logger) Warnf(ctx context.Context, msg string, args ...any) {
	l.Logger.Ctx(ctx).Warn(fmt.Sprintf(msg, args...), l.Fields...)
}

func (l Logger) Error(ctx context.Context, err error, msg string) {
	var customErr *Error
	if errors.As(err, &customErr) {
		if msg == "" {
			l.Logger.
				Ctx(ctx).
				Error(
					customErr.message,
					zap.Int("code", customErr.code),
					zap.Int("appCode", int(customErr.appCode)),
					zap.Any("reasons", customErr.reasons),
					zap.Any("callers", customErr.callers),
				)
			return
		}
		l.Logger.
			Ctx(ctx).
			Error(
				msg,
				zap.Error(err),
				zap.Int("code", customErr.code),
				zap.Int("appCode", int(customErr.appCode)),
				zap.Any("reasons", customErr.reasons),
				zap.Any("callers", customErr.callers),
			)
		return
	}
	if msg == "" {
		l.Logger.Ctx(ctx).Error(err.Error(), l.Fields...)
		return
	}
	l.Logger.Ctx(ctx).Error(msg, mergeFields(l.Fields, []zap.Field{zap.Error(err)})...)
}

func (l Logger) Errorf(ctx context.Context, err error, msg string, args ...any) {
	var customErr *Error
	if errors.As(err, &customErr) {
		l.Logger.
			Ctx(ctx).
			Error(
				fmt.Sprintf(msg, args...),
				zap.Error(err),
				zap.Int("code", customErr.code),
				zap.Int("appCode", int(customErr.appCode)),
				zap.Any("reasons", customErr.reasons),
				zap.Any("callers", customErr.callers),
			)
		return
	}
	l.Logger.Ctx(ctx).Error(fmt.Sprintf(msg, args...), mergeFields(l.Fields, []zap.Field{zap.Error(err)})...)
}

func (l Logger) ErrorFields(ctx context.Context, err error, msg string, fields ...zap.Field) {
	var customErr *Error
	if errors.As(err, &customErr) {
		fields = mergeFields(l.Fields, fields)
		fields = append(
			fields,
			zap.Int("code", customErr.code),
			zap.Int("appCode", int(customErr.appCode)),
			zap.Any("reasons", customErr.reasons),
			zap.Any("callers", customErr.callers),
		)
		if msg == "" {
			l.Logger.
				Ctx(ctx).
				Error(
					customErr.message,
					fields...,
				)
			return
		}
		fields = append(fields, zap.Error(err))
		l.Logger.
			Ctx(ctx).
			Error(
				msg,
				fields...,
			)
		return
	}
	if msg == "" {
		l.Logger.Ctx(ctx).Error(err.Error(), fields...)
		return
	}
	fields = append(fields, zap.Error(err))
	l.Logger.Ctx(ctx).Error(msg, fields...)
}

func (l Logger) Panic(ctx context.Context, msg string) {
	l.Logger.Ctx(ctx).Panic(msg)
}

func (l Logger) Panicf(ctx context.Context, msg string, args ...any) {
	l.Logger.Ctx(ctx).Panic(fmt.Sprintf(msg, args...), l.Fields...)
}

func (l Logger) Fatal(ctx context.Context, msg string) {
	l.Logger.Ctx(ctx).Fatal(msg, l.Fields...)
}

func (l Logger) Fatalf(ctx context.Context, msg string, args ...any) {
	l.Logger.Ctx(ctx).Fatal(fmt.Sprintf(msg, args...), l.Fields...)
}
