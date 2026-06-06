package commons

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Hs *http.Server
}

func NewServer(router *http.ServeMux) (*Server, error) {
	host := os.Getenv("APP_HOST")
	port := os.Getenv("APP_PORT")
	if IsRunningInDockerContainer() {
		// internal port is always 80; see dockerfile for port mapping
		host = ""
		port = "80"
	}
	srv := http.Server{
		Addr:    host + ":" + port,
		Handler: router,
	}
	return &Server{Hs: &srv}, nil
}

func (s *Server) RunServer(
	ctx context.Context,
	shutdownTimeout time.Duration,
) error {
	// see [Graceful Shutdown in Go: Key Patterns you need to know!](https://www.youtube.com/watch?v=UPVSeZXBTxI)
	serverErr := make(chan error, 1)
	go func() {
		slog.Debug("server::RunServer() - Started listening")
		if err := s.Hs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
		close(serverErr)
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	select {
	case err := <-serverErr:
		return err
	case <-stop:
		slog.Debug("server::RunServer() - Shutdown signal received, stopping server")
	case <-ctx.Done():
		slog.Debug("server::RunServer() - Context cancelled, stopping server")
	}
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		shutdownTimeout)
	defer cancel()
	if err := s.Hs.Shutdown(shutdownCtx); err != nil {
		if closeErr := s.Hs.Close(); closeErr != nil {
			return errors.Join(err, closeErr)
		}
		return err
	}
	slog.Debug("server::RunServer() - Server stopped gracefully")
	return nil
}

func IsRunningInDockerContainer() bool {
	// docker creates a .dockerenv file at the root
	// of the directory tree inside the container.
	// if this file exists then the viewer is running
	// from inside a container so return true
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}
