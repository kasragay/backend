package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/kasragay/backend/internal/ports"
	"github.com/kasragay/backend/internal/server/post"
)

var port = os.Getenv("PORT")

func init() {
	if port == "" {
		port = "8083"
	}
}

func gracefulShutdown(name string, server ports.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Printf("shutting %s down gracefully, press Ctrl+C again to force", name)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.App().ShutdownWithContext(ctx); err != nil {
		log.Printf("%s forced to shutdown with error: %v", name, err)
	}
	log.Printf("%s exiting", ``)
	done <- true
}

func main() {
	name := "post"
	post := post.New()
	post.RegisterRoutes()
	done := make(chan bool, 1)

	go func() {
		port, _ := strconv.Atoi(port)
		err := post.App().Listen(fmt.Sprintf(":%d", port))
		if err != nil {
			panic(fmt.Sprintf("http %s error: %s", name, err))
		}
	}()

	go gracefulShutdown(name, post, done)

	<-done
	log.Printf("%s graceful shutdown complete.", name)
}
