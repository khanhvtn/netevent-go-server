package middlewares

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di"
)

func InjectContainerMiddleware(container di.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("container", container)
		c.Next()
	}
}

// PanicRecoveryMiddleware handles the panic in the handlers.
func PanicRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				// log the error
				log.Fatalln(fmt.Sprint(rec))

				// write the error response
				c.JSON(500, map[string]interface{}{
					"error": "Internal Error",
				})
			}
		}()
		c.Next()
	}
}
