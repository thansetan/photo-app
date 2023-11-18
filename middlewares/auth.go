package middlewares

import (
	"net/http"
	"photo-app/helpers"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var re = regexp.MustCompile(`^[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*$`)

func AuthMiddleware(strict bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("Authorization")
		if token == "" {
			if !strict {
				ctx.Next()
				return
			}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "empty token",
			})
			return
		}

		if !strings.Contains(token, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": `invalid token format: please use "Bearer <your-token-here>`,
			})
			return
		}

		token = strings.Split(token, " ")[1]
		if !re.Match([]byte(token)) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": `invalid token`,
			})
			return
		}

		claims, err := helpers.ParseJWT(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.Set("id", claims.ID)
		ctx.Next()
	}
}
