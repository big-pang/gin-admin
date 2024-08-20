package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
)

var csrfMd func(http.Handler) http.Handler

func CSRF() gin.HandlerFunc {
	csrfMd = csrf.Protect([]byte("2vqp1fzk3nb9ki3c0wi08fki4knnw5rd"),
		csrf.HttpOnly(true),
		csrf.Secure(false),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"message": "Forbidden - CSRF token invalid"}`))
		})),
	)
	return adapter.Wrap(csrfMd)
}

func SCRFToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := csrf.Token(c.Request)
		c.Header("X-CSRF-Token", token)
		c.Set("csrf_token", token)
		c.Next()
	}
}
