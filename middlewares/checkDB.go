package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khanhvtn/netevent-go/database"
)

/*CheckDB : Check the DB connection before to execute a handle func*/
func CheckDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !database.ConnectionOK() {
			c.String(http.StatusInternalServerError, "Cannot connect to database")
		}
		c.Next()
	}
}
