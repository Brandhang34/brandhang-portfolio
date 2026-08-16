package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

// defaultPort matches the port EXPOSEd by the Dockerfile.
const defaultPort = 8080

type Server struct {
	port int
}

func NewServer() *http.Server {
	// A missing or unparseable PORT used to fall through to 0, which binds a
	// random port and makes the container unreachable on its published port.
	port := defaultPort
	if raw := os.Getenv("PORT"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 0 || parsed > 65535 {
			log.Printf("invalid PORT %q, falling back to %d", raw, defaultPort)
		} else {
			port = parsed
		}
	}

	NewServer := &Server{
		port: port,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
