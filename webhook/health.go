package webhook

import (
	"github.com/edimarlnx/rinha-de-backend-2024-q1-edimarlnx/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func wsHealth(e *gin.RouterGroup) {
	e.GET("/_health", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, &types.Response{
			Status:  http.StatusOK,
			Message: "up",
		})
	})
}
