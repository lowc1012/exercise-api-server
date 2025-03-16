package v1

import "github.com/gin-gonic/gin"

func MountRoute(r *gin.RouterGroup) {
	r.GET("/healthz", healthHandler)

	// tasks
	InitTaskHandler()
	r.GET("/tasks", fetchAllTasksHandler)
	r.GET("/tasks/:id", getTaskHandler)
	r.POST("/tasks", createTaskHandler)
	r.PUT("/tasks/:id", putTaskHandler)
	r.DELETE("/tasks/:id", deleteTaskHandler)
}
