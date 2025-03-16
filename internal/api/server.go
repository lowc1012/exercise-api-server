package api

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	v1 "github.com/lowc1012/exercise-api-server/internal/api/v1"
)

func StartAsync() (*http.Server, error) {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE"},
	}))

	router.GET("/healthz", v1.HealthHandler)

	// tasks
	v1.InitTaskHandler()
	router.GET("/tasks", v1.FetchAllTasksHandler)
	router.GET("/tasks/:id", v1.GetTaskHandler)
	router.POST("/tasks", v1.CreateTaskHandler)
	router.PUT("/tasks/:id", v1.PutTaskHandler)
	router.DELETE("/tasks/:id", v1.DeleteTaskHandler)

	// TODO: make more configurable
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Print("HTTP API server closed")
			} else {
				log.Fatalf("error: %v", err.Error())
			}
		}
	}()

	return srv, nil
}
