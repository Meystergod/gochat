package httpserver

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

type ServerDeps struct {
	Host string
	Port string
}

type Server struct {
	host       string
	port       string
	echoServer *echo.Echo
}

func NewDefaultServer(deps *ServerDeps) *Server {
	echoServer := echo.New()

	echoServer.Use(middleware.Recover())
	echoServer.Debug = true
	echoServer.DisableHTTP2 = true
	echoServer.HideBanner = true
	echoServer.HidePort = true

	s := &Server{
		host:       deps.Host,
		port:       deps.Port,
		echoServer: echoServer,
	}

	s.echoServer = echoServer

	return s
}

func (s *Server) Server() *echo.Echo {
	return s.echoServer
}

func (s *Server) Start() error {
	if err := s.echoServer.Start(fmt.Sprintf("%s:%s", s.host, s.port)); err != nil {
		return errors.Wrap(err, "starting echo server")
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.echoServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "shutdown echo server")
	}

	return nil
}
