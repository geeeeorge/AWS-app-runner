package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"

	"github.com/geeeeorge/AWS-app-runner/src/app/handler"

	log "github.com/sirupsen/logrus"
)

// Server represents server
type Server struct {
	DB          *sqlx.DB
	Host        string
	Port        int
	ServerReady chan<- interface{}
	echo        *echo.Echo
	shutdown    chan interface{}
}

// NewServer returns new Server object
func NewServer(port int, host string, db *sqlx.DB, serverReady chan<- interface{}) *Server {
	return &Server{
		DB:          db,
		Host:        host,
		Port:        port,
		ServerReady: serverReady,
	}
}

func (s *Server) setup() {
	e := echo.New()
	h := handler.New()

	e.GET("/api/_healthz", h.GetApiHealthz)
	s.echo = e

	viper.SetEnvPrefix("AWS")
	viper.AutomaticEnv()
	s.echo.Use(
		middleware.Recover(),
		middleware.Logger(),
		middleware.RequestID(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowCredentials: true,
		}),
	)
}

// GetAddress return server address
func (s *Server) GetAddress() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// Start starts server
func (s *Server) Start() {
	s.setup()

	go func() {
		if err := s.echo.Start(s.GetAddress()); err != nil && err != http.ErrServerClosed {
			log.Error(err.Error())
			log.Info("shutting down the server")
		}
	}()

	if s.ServerReady != nil {
		s.ServerReady <- struct{}{}
	}

	s.shutdown = make(chan interface{}, 1)
	defer close(s.shutdown)
	quit := make(chan os.Signal, 1)
	defer close(quit)
	signal.Notify(quit, os.Interrupt)
	select {
	case sig := <-quit:
		log.Info("received: ", sig)
	case <-s.shutdown:
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.echo.Shutdown(ctx); err != nil {
		log.Fatal("failed to graceful shutdown the server: ", err)
	}
}

// Shutdown shutdowns the server
func (s *Server) Shutdown() {
	s.shutdown <- struct{}{}
}
