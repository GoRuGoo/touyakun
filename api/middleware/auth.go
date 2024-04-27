package middleware

import (
	"github.com/gin-gonic/gin"
)

package middleware

import (
"encoding/json"
"io"
"log"
"net/http"
"net/url"
"os"
"strings"

"github.com/gin-gonic/gin"
)
type responseBody struct {
	Sub string `json:"sub"`
	Err string `json:"error_description"`
}

func AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// #23でLINEのAuthサーバーに問い合わせるmiddlewareを実装予定
		c.Set("auth_key", "test_auth")
		c.Next()
	}
}
