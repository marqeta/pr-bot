package prbot

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

type Server struct {
	*http.Server
	config *Config
}

// NewServer creates and configures a server
func NewServer(cfg *Config, routes http.Handler) *Server {

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      routes,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	return &Server{
		Server: &srv,
		config: cfg,
	}
}

func (srv *Server) Start() {
	if srv.config.Server.DisableTLS {
		srv.startHTTP()
	} else {
		srv.startHTTPS()
	}
}

func (srv *Server) startHTTP() {
	log.Info().Msg("Starting HTTP Server....")
	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Err(err).Msg("error starting server")
		os.Exit(1)
	}
}

func (srv *Server) startHTTPS() {
	log.Info().Msg("Starting HTTPS Server....")
	certFile := srv.config.Server.TLS.CertFile
	KeyFile := srv.config.Server.TLS.KeyFile
	err := srv.ListenAndServeTLS(certFile, KeyFile)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Err(err).Msg("error starting server")
		os.Exit(1)
	}
}

func (srv *Server) WaitForGracefulShutdown() {

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-quit
	log.Info().Str("Reason", sig.String()).Msg("Server is shutting down.")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		log.Err(err).Msg("Could not gracefully shutdown the server")
		os.Exit(1)
	}
	log.Info().Msg("Server stopped")
}
