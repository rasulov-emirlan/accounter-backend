package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Logger struct {
		logger *zap.Logger
		closer func() error
	}

	ErrCloser struct {
		Errors []error
	}

	Field = zap.Field
)

var logLevels = map[string]zapcore.Level{
	"dev":    zapcore.DebugLevel,
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func (e *ErrCloser) Error() string {
	err := ""
	for i, v := range e.Errors {
		err += v.Error()
		if i != len(e.Errors)-1 {
			err += ","
		}
	}
	return err
}

// If logLevel is not valid, it will be set to debug.
func NewLogger(logLevel string) (*Logger, error) {
	var (
		logger Logger
		err    error
	)
	config := zap.NewDevelopmentConfig()

	lvl := zapcore.DebugLevel
	if l, ok := logLevels[logLevel]; ok {
		lvl = l
	}

	config.EncoderConfig.StacktraceKey = zapcore.OmitKey
	config.EncoderConfig.CallerKey = zapcore.OmitKey
	if logLevel == "dev" {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	config.Level = zap.NewAtomicLevelAt(lvl)

	logger.closer = func() error {
		// If this logger will have any connections to external services, close them here.
		// If there are no connections.
		return nil
	}

	logger.logger, err = config.Build()
	if err != nil {
		if errr := logger.closer(); errr != nil {
			return nil, &ErrCloser{
				Errors: []error{errr, err},
			}
		}
		return nil, err
	}

	return &logger, nil
}

func (l *Logger) Close() error {
	return l.closer()
}

func (l Logger) Goosed() GooseLogger {
	return GooseLogger{Logger: l.logger}
}

func String(key, data string) Field {
	return zap.String(key, data)
}

func Int(key string, data int) Field {
	return zap.Int(key, data)
}

func Int64(key string, data int64) Field {
	return zap.Int64(key, data)
}

func Uint(key string, data uint) Field {
	return zap.Uint(key, data)
}

func Uint64(key string, data uint64) Field {
	return zap.Uint64(key, data)
}

func Float64(key string, data float64) Field {
	return zap.Float64(key, data)
}

func Bool(key string, data bool) Field {
	return zap.Bool(key, data)
}

func Any(key string, data any) Field {
	return zap.Any(key, data)
}

func Error(key string, data error) Field {
	return zap.NamedError(key, data)
}

// Sync flushes any buffered log entries.
// It might be a good idea to call Sync before exiting your program.
// Or defer it in the beginning of your handlers.
func (l *Logger) Sync() error {
	return l.logger.Sync()
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	fields = append(fields, zap.Stack("stacktrace"))
	l.logger.Fatal(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...Field) {
	fields = append(fields, zap.Stack("stacktrace"))
	l.logger.Panic(msg, fields...)
}

func (l *Logger) DPanic(msg string, fields ...Field) {
	fields = append(fields, zap.Stack("stacktrace"))
	l.logger.DPanic(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, fields...)
}
