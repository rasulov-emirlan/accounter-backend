package logging

import (
	"fmt"

	"go.uber.org/zap"
)

type GooseLogger struct {
	*zap.Logger
}

func (l GooseLogger) Fatal(v ...interface{}) {
	l.Logger.Fatal("goose migrations", zap.Any("msg", v))
}

func (l GooseLogger) Fatalf(format string, v ...interface{}) {
	l.Logger.Fatal("goose migrations", zap.String("msg", fmt.Sprintf(format, v...)))
}

func (l GooseLogger) Print(v ...interface{}) {
	l.Logger.Info("goose migrations", zap.Any("msg", v))
}

func (l GooseLogger) Println(v ...interface{}) {
	l.Logger.Info("goose migrations", zap.Any("msg", v))
}

func (l GooseLogger) Printf(format string, v ...interface{}) {
	l.Logger.Info("goose migrations", zap.String("msg", fmt.Sprintf(format, v...)))
}
