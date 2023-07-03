package logger

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type recordingType int

type logger struct {
	ctx       context.Context
	config    *Config
	logger    *zap.SugaredLogger
	zapLogger *zap.Logger
	fields    []interface{}
	skipInit  bool
	tracing   recordingType
}

func newLogger(config *Config) *logger {
	config.buildZapConfig()

	zapLogger, err := config.zapConfig.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error on build zap logger (%s)", err)
		return nil
	}
	zapLogger = zapLogger.WithOptions(zap.AddCallerSkip(config.CallerSkip))

	return &logger{
		logger:    zapLogger.Sugar(),
		zapLogger: zapLogger,
		config:    config,
	}
}

func (l *logger) SetLevel(level Level) {
	l.config.Level = level
	l.config.zapConfig.Level.SetLevel(zapcore.Level(level))
}

func (l *logger) Level() Level {
	return l.config.Level
}

func (l *logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *logger) DPanic(args ...interface{}) {
	l.logger.DPanic(args...)
}

func (l *logger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *logger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args...)
}

func (l *logger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}

func (l *logger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}

func (l *logger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
}

func (l *logger) DPanicf(template string, args ...interface{}) {
	l.logger.DPanicf(template, args...)
}

func (l *logger) Panicf(template string, args ...interface{}) {
	l.logger.Panicf(template, args...)
}

func (l *logger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args...)
}

func (l *logger) Debugw(msg string, keysAndValues ...interface{}) {
	keysAndValues = l.tracingEvent(zapcore.DebugLevel, msg, keysAndValues...)
	l.logger.Debugw(msg, keysAndValues...)
}

func (l *logger) Infow(msg string, keysAndValues ...interface{}) {
	keysAndValues = l.tracingEvent(zapcore.InfoLevel, msg, keysAndValues...)
	l.logger.Infow(msg, keysAndValues...)
}

func (l *logger) Warnw(msg string, keysAndValues ...interface{}) {
	keysAndValues = l.tracingEvent(zapcore.WarnLevel, msg, keysAndValues...)
	l.logger.Warnw(msg, keysAndValues...)
}

func (l *logger) Errorw(msg string, keysAndValues ...interface{}) {
	keysAndValues = l.tracingEvent(zapcore.ErrorLevel, msg, keysAndValues...)
	l.logger.Errorw(msg, keysAndValues...)
}

func (l *logger) DPanicw(msg string, keysAndValues ...interface{}) {
	keysAndValues = l.tracingEvent(zapcore.DPanicLevel, msg, keysAndValues...)
	l.logger.DPanicw(msg, keysAndValues...)
}

func (l *logger) Panicw(msg string, keysAndValues ...interface{}) {
	keysAndValues = l.tracingEvent(zapcore.PanicLevel, msg, keysAndValues...)
	l.logger.Panicw(msg, keysAndValues...)
}

func (l *logger) Fatalw(msg string, keysAndValues ...interface{}) {
	keysAndValues = l.tracingEvent(zapcore.FatalLevel, msg, keysAndValues...)
	l.logger.Fatalw(msg, keysAndValues...)
}

func (l *logger) Log(level Level, args ...interface{}) {
	// nolint: exhaustive
	switch level {
	case DebugLevel:
		l.logger.Debug(args...)
	case InfoLevel:
		l.logger.Info(args...)
	case WarnLevel:
		l.logger.Warn(args...)
	case ErrorLevel:
		l.logger.Error(args...)
	default:
		l.logger.Info(args...)
	}
}

func (l *logger) Logf(level Level, template string, args ...interface{}) {
	// nolint: exhaustive
	switch level {
	case DebugLevel:
		l.logger.Debugf(template, args...)
	case InfoLevel:
		l.logger.Infof(template, args...)
	case WarnLevel:
		l.logger.Warnf(template, args...)
	case ErrorLevel:
		l.logger.Errorf(template, args...)
	default:
		l.logger.Infof(template, args...)
	}
}

func (l *logger) Logw(level Level, msg string, keysAndValues ...interface{}) {
	// nolint: exhaustive
	switch level {
	case DebugLevel:
		l.logger.Debugw(msg, keysAndValues...)
	case InfoLevel:
		l.logger.Infow(msg, keysAndValues...)
	case WarnLevel:
		l.logger.Warnw(msg, keysAndValues...)
	case ErrorLevel:
		l.logger.Errorw(msg, keysAndValues...)
	default:
		l.logger.Infow(msg, keysAndValues...)
	}
}

func (l *logger) Sync() error {
	return l.logger.Sync()
}

func (l *logger) With(keyValues ...interface{}) Logger {
	return l.WithCallerSkip(l.ctx, defaultCallerSkip, l.tracing, keyValues...)
}

func (l *logger) WithSkip(callerSkip int, keyValues ...interface{}) Logger {
	return l.WithCallerSkip(l.ctx, callerSkip, l.tracing, keyValues...)
}

const (
	TraceIDOnly recordingType = 0
	TraceEvent  recordingType = 1
)

const defaultCallerSkip = -1

func (l *logger) WithCallerSkip(ctx context.Context, callerSkip int, tracing recordingType, keyValues ...interface{}) Logger {
	newFields := make([]interface{}, len(l.fields), len(l.fields)+len(keyValues))
	copy(newFields, l.fields)
	newFields = append(newFields, keyValues...)

	sugar := l.logger.With(keyValues...)
	zaplogger := l.zapLogger

	// only first With need skip caller, be aware DO NOT affect parent logger
	if !l.skipInit && callerSkip != 0 {
		zaplogger = l.logger.With(keyValues...).Desugar().WithOptions(zap.AddCallerSkip(callerSkip))
		sugar = zaplogger.Sugar()
	}

	newLogger := &logger{
		ctx:       ctx,
		fields:    newFields,
		config:    l.config.clone(),
		logger:    sugar,
		zapLogger: zaplogger,
		skipInit:  true,
		tracing:   tracing,
	}
	return newLogger
}

func (l *logger) WithTraceID(ctx context.Context, keyValues ...interface{}) Logger {
	return l.WithCallerSkip(ctx, defaultCallerSkip, TraceIDOnly, keyValues...)
}

func (l *logger) Ctx(ctx context.Context) Logger {
	return l.WithCallerSkip(ctx, defaultCallerSkip, TraceEvent)
}

func (l *logger) GetZapLogger() *zap.Logger {
	return l.zapLogger
}

func (l *logger) tracingEvent(lvl zapcore.Level, msg string, keysAndValues ...interface{}) []interface{} {
	if l.ctx == nil {
		return keysAndValues
	}

	// skip handling tracing if current logging level is not enabled
	if !l.config.zapConfig.Level.Level().Enabled(lvl) {
		return keysAndValues
	}

	span := trace.SpanFromContext(l.ctx)
	if span.IsRecording() {
		if l.tracing == TraceEvent {
			if zapcore.ErrorLevel.Enabled(lvl) {
				span.SetStatus(codes.Error, msg)
			}

			// Allocate enough space for the worst case; if users pass only structured
			// fields, we shouldn't penalize them with extra allocations.
			attrs := make([]attribute.KeyValue, 0)

			attrs = append(attrs, logSeverityKey.String(levelString(lvl)))
			attrs = append(attrs, logMessageKey.String(msg))

			for i := 0; i < len(l.fields); {
				// Make sure this element isn't a dangling key.
				if i == len(l.fields)-1 {
					break
				}

				// Consume this value and the next, treating them as a key-value pair. If the
				// key isn't a string, add this pair to the slice of invalid pairs.
				key, val := l.fields[i], l.fields[i+1]
				if keyStr, ok := key.(string); ok {
					f := zap.Any(keyStr, val)
					attrs = appendField(attrs, f)
				}
				i += 2
			}

			for i := 0; i < len(keysAndValues); {
				// Make sure this element isn't a dangling key.
				if i == len(keysAndValues)-1 {
					break
				}

				// Consume this value and the next, treating them as a key-value pair. If the
				// key isn't a string, add this pair to the slice of invalid pairs.
				key, val := keysAndValues[i], keysAndValues[i+1]
				if keyStr, ok := key.(string); ok {
					f := zap.Any(keyStr, val)
					attrs = appendField(attrs, f)
				}
				i += 2
			}

			span.AddEvent("log", trace.WithAttributes(attrs...))
		}

		if s := span.SpanContext(); s.HasTraceID() {
			// keysAndValues = append([]interface{}{"trace_id", s.TraceID().String()}, keysAndValues...)
			keysAndValues = append(keysAndValues, "trace_id", s.TraceID().String())
		}
	}
	return keysAndValues
}
