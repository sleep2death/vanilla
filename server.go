package vanilla

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	server *http.Server
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
	api.GET("/playerinfo", getPlayerInfoHandler(db))

	ws := router.Group("/ws")
	ws.GET("", getWSHandler(db))

	return router, nil
}

// Run the server
func Run(addr string) {
	router, err := setupRouter()
	if err != nil {
		log.Fatal(err)
	}

	server = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("gin server error: %s", err)
		}
	}()
}

// Stop the server
func Stop() {
	if server != nil {
		log.Println("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
		// catching ctx.Done(). timeout of 1 seconds.
		select {
		case <-ctx.Done():
			log.Println("timeout of 1 second.")
		}
		log.Println("Server exiting")
	}
}
