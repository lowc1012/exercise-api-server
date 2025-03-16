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

	v1Group := router.Group("/v1")
	v1.MountRoute(v1Group)

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
