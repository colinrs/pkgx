package gin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Health ...
func Health(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}