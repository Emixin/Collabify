package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// TODO: Add a require auth middleware and use sessions here!
func RequireAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		session := sessions.Default(context)
		_, ok := session.Get("user_id").(int)

		if !ok {
			context.Redirect(http.StatusFound, "/login")
			return
		}

		context.Next()
	}
}
