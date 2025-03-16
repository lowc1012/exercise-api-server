package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lowc1012/exercise-api-server/internal/domain"
	"github.com/lowc1012/exercise-api-server/internal/repository/memory"
	"github.com/lowc1012/exercise-api-server/internal/task"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

type TaskService interface {
	FetchAll(ctx context.Context) ([]domain.Task, error)
	GetByID(ctx context.Context, id string) (domain.Task, error)
	Update(ctx context.Context, t domain.Task) error
	Store(ctx context.Context, t domain.Task) error
	Delete(ctx context.Context, id string) error
}

var taskService TaskService

// InitTaskHandler initializes the task service
func InitTaskHandler() {
	tr := memory.NewTaskRepository()
	taskService = task.NewService(tr)
}

func FetchAllTasksHandler(c *gin.Context) {
	tasks, err := taskService.FetchAll(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks":  tasks,
		"status": "success",
	})
}

func GetTaskHandler(c *gin.Context) {
	id := c.Param("id")
	t, err := taskService.GetByID(c, id)
	// TODO: improve validation
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
		return
	}
	if t.ID == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, ResponseError{Message: "task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"task":   t,
		"status": "success",
	})
}

func CreateTaskHandler(c *gin.Context) {
	var t domain.Task
	var err error
	if err = c.ShouldBindJSON(&t); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseError{Message: "invalid request body"})
		return
	}
	// TODO: improve validation
	if t.Name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseError{Message: "name is required"})
		return
	}

	if t.Status != 0 && t.Status != 1 {
		t.Status = 0
	}

	if err := taskService.Store(c.Request.Context(), t); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ResponseError{Message: "failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"task":    t,
		"status":  "success",
		"message": "task created successfully",
	})
}

func PutTaskHandler(c *gin.Context) {
	id := c.Param("id")
	existedTask, err := taskService.GetByID(c, id)
	// TODO: improve validation
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
		return
	}
	if existedTask.ID == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, ResponseError{Message: "task not found"})
		return
	}
	var t domain.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseError{Message: "invalid request body"})
		return
	}
	existedTask = t
	err = taskService.Update(c, existedTask)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ResponseError{Message: "failed to create task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task":   t,
		"status": "success",
	})
}

func DeleteTaskHandler(c *gin.Context) {
	id := c.Param("id")
	existedTask, err := taskService.GetByID(c, id)
	// TODO: improve validation
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
		return
	}
	if existedTask.ID == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, ResponseError{Message: "task not found"})
		return
	}
	err = taskService.Delete(c, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
