package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware - middleware для логирования запросов
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now() // Запоминаем время начала обработки запроса

		// Обрабатываем следующий middleware/handler
		c.Next()

		// После обработки запроса
		duration := time.Since(start) // Вычисляем время обработки
		status := c.Writer.Status()   // Получаем статус ответа
		method := c.Request.Method    // Получаем HTTP метод
		path := c.Request.URL.Path    // Получаем путь запроса

		// Логируем информацию о запросе
		log.Printf("%s %s %d %s", method, path, status, duration)
	}
}
