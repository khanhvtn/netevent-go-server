package middlewares

import (
	"context"

	"github.com/gin-gonic/gin"
)

func ContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "gincontext", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
