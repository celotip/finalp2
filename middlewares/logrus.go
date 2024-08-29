package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func LogrusMiddleware(logger *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//log request
			logger.WithFields(logrus.Fields{
				"method": c.Request().Method,
				"url":    c.Request().URL.String(),
				"ip":     c.RealIP(),
			}).Info("Request Received")

			//proses ke next middleware/handler
			err := next(c)

			//log response status
			logger.WithFields(logrus.Fields{
				"status": c.Response().Status,
			}).Info("Response sent")

			return err
		}
	}
}
