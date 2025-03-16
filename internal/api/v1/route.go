package v1

import "github.com/gin-gonic/gin"

func MountRoute(r *gin.RouterGroup) {
	r.GET("/healthz", healthHandler)

	// tasks
	r.GET("/tasks", getTasksHandler)
	r.POST("/tasks", createTaskHandler)
	r.PUT("/tasks/:id", putTaskHandler)
	r.DELETE("/tasks/:id", deleteTaskHandler)
}
