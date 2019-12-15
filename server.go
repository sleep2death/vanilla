package vanilla

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	srv *http.Server
)

type user struct {
	name string
}

func setupRouter() (*gin.Engine, error) {
	db, err := initDB("mongodb://localhost:27017")
	if err != nil {
		return nil, err
	}

	router := gin.Default()
	router.Use(CORSMiddleware())

	router.POST("/login", getLoginHandler(db))
	router.POST("/register", getRegisterHandler(db))

	api := router.Group("/api")
	api.Use(authMiddleware())
	api.GET("/ping", getPingHandler())

	return router, nil
}

// Start the server
func Start(addr string) {
	r, err := setupRouter()
	if err != nil {
		log.Fatal(err)
	}
	srv = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

// Stop the server
func Stop() {
	if srv != nil {
		log.Println("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
		// catching ctx.Done(). timeout of 1 seconds.
		select {
		case <-ctx.Done():
			log.Println("timeout of 1 seconds.")
		}
		log.Println("Server exiting")
	}
}
