package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lowc1012/exercise-api-server/internal/domain"
	"github.com/stretchr/testify/assert"
)

type mockTaskService struct {
	tasks map[string]domain.Task
}

func newMockTaskService() *mockTaskService {
	return &mockTaskService{
		tasks: make(map[string]domain.Task),
	}
}

func (m *mockTaskService) FetchAll(ctx context.Context) ([]domain.Task, error) {
	var tasks []domain.Task
	for _, task := range m.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (m *mockTaskService) GetByID(ctx context.Context, id string) (domain.Task, error) {
	task, exists := m.tasks[id]
	if !exists {
		return domain.Task{}, nil
	}
	return task, nil
}

func (m *mockTaskService) Store(ctx context.Context, t domain.Task) error {
	t.ID = "1" // Simulating auto-increment
	m.tasks[t.ID] = t
	return nil
}

func (m *mockTaskService) Update(ctx context.Context, t domain.Task) error {
	if _, exists := m.tasks[t.ID]; !exists {
		return nil
	}
	m.tasks[t.ID] = t
	return nil
}

func (m *mockTaskService) Delete(ctx context.Context, id string) error {
	delete(m.tasks, id)
	return nil
}

func setupTest() (*gin.Engine, *mockTaskService) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mock := newMockTaskService()
	taskService = mock
	return r, mock
}

func TestFetchAllTasksHandler(t *testing.T) {
	r, mock := setupTest()
	r.GET("/tasks", FetchAllTasksHandler)

	mock.tasks["1"] = domain.Task{ID: "1", Name: "Test Task"}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])

	tasks := response["tasks"].([]interface{})
	assert.Equal(t, 1, len(tasks))
}

func TestGetTaskHandler(t *testing.T) {
	r, mock := setupTest()
	r.GET("/tasks/:id", GetTaskHandler)

	mock.tasks["1"] = domain.Task{ID: "1", Name: "Test Task"}

	t.Run("existing task", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tasks/1", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "success", response["status"])
	})

	t.Run("non-existent task", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tasks/999", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestCreateTaskHandler(t *testing.T) {
	r, _ := setupTest()
	r.POST("/tasks", CreateTaskHandler)

	t.Run("valid task", func(t *testing.T) {
		task := domain.Task{Name: "New Task"}
		body, _ := json.Marshal(task)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "success", response["status"])
	})

	t.Run("invalid task", func(t *testing.T) {
		task := domain.Task{} // Empty name
		body, _ := json.Marshal(task)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPutTaskHandler(t *testing.T) {
	r, mock := setupTest()
	r.PUT("/tasks/:id", PutTaskHandler)

	mock.tasks["1"] = domain.Task{ID: "1", Name: "Original Task"}

	t.Run("update existing task", func(t *testing.T) {
		task := domain.Task{Name: "Updated Task"}
		body, _ := json.Marshal(task)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("update non-existent task", func(t *testing.T) {
		task := domain.Task{Name: "New Task"}
		body, _ := json.Marshal(task)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tasks/999", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDeleteTaskHandler(t *testing.T) {
	r, mock := setupTest()
	r.DELETE("/tasks/:id", DeleteTaskHandler)

	mock.tasks["1"] = domain.Task{ID: "1", Name: "Test Task"}

	t.Run("delete existing task", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("delete non-existent task", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tasks/999", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
