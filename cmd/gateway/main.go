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
	"github.com/kasragay/backend/internal/server/gateway"

	_ "github.com/joho/godotenv/autoload"
)

var (
	tlsCrtFile = os.Getenv("TLS_CRT_FILE")
	tlsKeyFile = os.Getenv("TLS_KEY_FILE")
	port       = os.Getenv("PORT")
)

func init() {
	if tlsCrtFile == "" {
		tlsCrtFile = "/etc/ssl/crt.pem"
	}
	if tlsKeyFile == "" {
		tlsKeyFile = "/etc/ssl/key.pem"
	}
	if port == "" {
		port = "8081"
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
	name := "gateway"
	gateway := gateway.New()
	gateway.RegisterRoutes()
	done := make(chan bool, 1)

	go func() {
		port, _ := strconv.Atoi(port)
		err := gateway.App().ListenTLS(fmt.Sprintf(":%d", port), tlsCrtFile, tlsKeyFile)
		if err != nil {
			panic(fmt.Sprintf("https %s error: %s", name, err))
		}
	}()

	go gracefulShutdown(name, gateway, done)

	<-done
	log.Printf("%s graceful shutdown complete.", name)
}
