package logging

import (
	"time"

	"github.com/labstack/echo/v4"
)

const reqLogMsg = "transport_log"

func (l *Logger) NewEchoMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		begin := time.Now()

		err := next(c)

		if err != nil {
			l.Error(
				reqLogMsg,
				String("method", c.Request().Method),
				String("path", c.Request().URL.Path),
				String("duration", time.Since(begin).String()),
				String("error", err.Error()),
			)
			return err
		}

		l.Info(
			reqLogMsg,
			String("method", c.Request().Method),
			String("path", c.Request().URL.Path),
			String("duration", time.Since(begin).String()),
		)

		return err
	}
}
