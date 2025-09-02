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

	"github.com/geo-afk/Online-Doctor-Appointment/app/server"
	_ "github.com/joho/godotenv/autoload"
)

func gracefulShutdown(fiberServer *server.FiberServer, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := fiberServer.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

//	@title			CD-Appointment
//	@version		1.0.0
//	@description	This is a client-doctor/s appointment project to do booking, consultation, etc

// contact.name Geovanni Stewart
// contact.url github.com/geo-afk/Online-Doctor-Appointment
// contact.email geo.a.atewart21@gmail.com

// @securityDefinitions.apiKey bearerToken
// @in header
// @name Authorization

//	@license.name	MIT
//	@license.url	github.com/geo-afk/Online-Doctor-Appointment/blob/main/LICENSE

// @host		localhost:3400
// @BasePath
func main() {

	s := server.New()
	s.RegisterFiberRoutes()

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	go func() {
		port, _ := strconv.Atoi(os.Getenv("PORT"))
		err := s.Listen(fmt.Sprintf(":%d", port))
		if err != nil {
			panic(fmt.Sprintf("http server error: %s", err))
		}
	}()

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(s, done)

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")

}
