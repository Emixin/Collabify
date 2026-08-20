package middlewares

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

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

// TODO: Add a JWT Middleware
func JWTMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString, err := context.Cookie("token")
		if err != nil {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// TODO: Complete Here
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {

			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("JWT Method is not valid!")
			}

			expirationTime, err := token.Claims.GetExpirationTime()
			if err != nil {
				return nil, errors.New("Failed to fetch token!")
			}

			if expirationTime.Time.After(time.Now()) {
				return nil, errors.New("Token expired!")
			}

			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				return nil, errors.New("JWT secret is not configured!")
			}

			return []byte(secret), nil
		})

		if err != nil {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		context.Next()
	}
}
