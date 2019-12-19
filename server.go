package vanilla

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
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

	return router, nil
}

func setupIO(router *gin.Engine) (*socketio.Server, error) {
	server, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})
	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})
	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})
	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})
	server.OnError("/", func(e error) {
		fmt.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	ioGroup := router.Group("/io")
	ioGroup.Use(authMiddleware())
	ioGroup.GET("/*any", gin.WrapH(server))

	return server, nil
}

// Run the server
func Run(addr string) {
	router, err := setupRouter()
	if err != nil {
		log.Fatal(err)
	}

	io, err := setupIO(router)
	if err != nil {
		log.Fatal(err)
	}

	defer io.Close()
	go io.Serve()

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

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
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
