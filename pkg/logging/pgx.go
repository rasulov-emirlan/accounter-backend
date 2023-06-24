package logging

import (
	"context"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// this is a compile-time check, to ensure that Logger implements pgx.QueryTracer
var _ pgx.QueryTracer = (*Logger)(nil)

// Function bellow are used to implement pgx.QueryTracer interface.
// In future it can be used to log queries and their execution time.
// Or even add an actual tracer to our logs.

func (l *Logger) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	l.logger.Debug("pgx start", zap.String("sql", data.SQL), zap.Any("args", data.Args))
	return ctx
}

func (l *Logger) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	if data.Err != nil {
		l.logger.Error("pgx end", zap.Error(data.Err))
	} else {
		l.logger.Debug("pgx end")
	}
}
