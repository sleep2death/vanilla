package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sleep2death/vanilla"
)

func main() {
	go func() { vanilla.Run(":8082") }()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	vanilla.Stop()
}
