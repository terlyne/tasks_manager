package logger

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

var (
	EnvDev   = "dev"
	EnvLocal = "local"
	EnvProd  = "prod"
)
var log *slog.Logger

func InitLogger(env string) {
	var handler slog.Handler

	switch env {
	case EnvDev:
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	case EnvLocal:
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	case EnvProd:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelError,
		})
	}

	log = slog.New(handler)
}

func GetLogger() *slog.Logger {
	return log
}

// Middleware для логирования HTTP-запросов
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Начало запроса
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Обработка запроса
		c.Next()

		// Конец запроса
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		// Логируем информацию о запросе
		log.Info("HTTP Request",
			"status", statusCode,
			"latency", latency,
			"client_ip", clientIP,
			"method", method,
			"path", path,
			"query", query,
		)
	}
}
