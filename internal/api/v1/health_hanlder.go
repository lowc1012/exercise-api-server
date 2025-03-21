package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthHandler(c *gin.Context) {
	type resp struct {
		Status string `json:"status"`
	}
	c.JSON(http.StatusOK, &resp{Status: "health"})
}
