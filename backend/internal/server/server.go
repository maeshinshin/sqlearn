package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/maeshinshin/sqlearn/backend/gen/problem/v1/problemv1connect"
	"github.com/maeshinshin/sqlearn/backend/internal/handler"
)

type Server struct {
	port    string
	mux     *http.ServeMux
	cors    *cors.Cors
	isDebug bool
}

type Option func(*Server)

func WithPort(port string) Option {
	slog.Info("Setting server port", "port", port)
	return func(s *Server) {
		s.port = port
	}
}

func WithDebug(debug bool) Option {
	slog.Info("Setting debug mode", "debug", debug)

	return func(s *Server) {
		s.isDebug = debug
	}
}

func WithHandler(h *handler.ProblemHandler) Option {
	slog.Info("Adding ProblemHandler to server", "handler", fmt.Sprintf("%T", h))
	return func(s *Server) {
		path, connectHandler := problemv1connect.NewProblemServiceHandler(h)
		s.mux.Handle(path, connectHandler)
	}
}

func WithCORS(cors *cors.Cors) Option {
	slog.Info("Setting custom CORS configuration")
	return func(s *Server) {
		s.cors = cors
	}
}

func NewServer(opts ...Option) *Server {
	s := &Server{
		port: "8080",
		mux:  http.NewServeMux(),
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.cors == nil {
		s.cors = cors.New(
			cors.Options{
				AllowedOrigins: func() []string {
					if s.isDebug {
						return []string{"*"}
					}
					return []string{"http://localhost:5173"}
				}(),
				AllowedMethods: []string{"GET", "POST", "OPTIONS"},
				AllowedHeaders: []string{"Accept", "Accept-Encoding", "Content-Type", "Connect-Protocol-Version"},
			},
		)
	}

	return s
}

func (s *Server) Start() error {
	slog.Info("Starting server on port " + s.port)

	handlerWithCORS := s.cors.Handler(s.mux)
	serverHandler := h2c.NewHandler(handlerWithCORS, &http2.Server{})

	err := http.ListenAndServe(":"+s.port, serverHandler)
	if err != nil {
		slog.Error("Server stopped", "error", err)
	}
	return err
}
