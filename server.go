package vanilla

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	r   = gin.Default()
	srv *http.Server
)

// Start the server
func Start(addr string) {
	// serve the static files in the client's dist
	r.Use(static.ServeRoot("/", "../client/dist"))
	r.Use(static.ServeRoot("/dist", "../client/dist"))

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
