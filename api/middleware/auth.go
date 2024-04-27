package middleware

import (
	"github.com/gin-gonic/gin"
)

func AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// #23でLINEのAuthサーバーに問い合わせるmiddlewareを実装予定
		c.Set("auth_key", "test_auth")
		c.Next()
	}
}
