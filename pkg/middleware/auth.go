package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware - middleware для проверки JWT
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// Проверяем, начинается ли токен с "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // Если токен не начинается с "Bearer "
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header"})
			c.Abort()
			return
		}

		// Проверяем токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Проверяем, что токен использует ожидаемый алгоритм
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorUnverifiable)
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Если токен действителен, продолжаем выполнение запроса
		c.Next()
	}
}

// WrapAuth - обертка для обработчиков с использованием аутентификации
func WrapAuth(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		AuthMiddleware(os.Getenv("JWT_SECRET"))(c) // Применяем middleware аутентификации
		if c.IsAborted() {                         // Если запрос был прерван (например, из-за неудачной аутентификации)
			return
		}
		handler(c) // Вызываем обернутый обработчик
	}
}
