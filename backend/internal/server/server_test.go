package server

import (
	"testing"

	"github.com/rs/cors"

	"github.com/maeshinshin/sqlearn/backend/internal/handler"
)

func TestNewServer_Options(t *testing.T) {
	tests := []struct {
		name     string
		opts     []Option
		validate func(*testing.T, *Server)
	}{
		{
			name: "Default values (Production mode)",
			opts: nil,
			validate: func(t *testing.T, s *Server) {
				if s.port != "8080" {
					t.Errorf("expected port 8080, got %s", s.port)
				}
				if s.isDebug {
					t.Error("expected default isDebug to be false")
				}
				if s.cors == nil {
					t.Error("expected default CORS to be initialized")
				}
				if s.mux == nil {
					t.Error("expected default mux to be initialized")
				}
			},
		},
		{
			name: "Debug mode",
			opts: []Option{WithDebug(true)},
			validate: func(t *testing.T, s *Server) {
				if !s.isDebug {
					t.Error("expected isDebug to be true")
				}
			},
		},
		{
			name: "Custom Port",
			opts: []Option{WithPort("9090")},
			validate: func(t *testing.T, s *Server) {
				if s.port != "9090" {
					t.Errorf("expected port 9090, got %s", s.port)
				}
			},
		},
		{
			name: "WithHandler registration",
			opts: []Option{WithHandler(&handler.ProblemHandler{})},
			validate: func(t *testing.T, s *Server) {
				if s.mux == nil {
					t.Fatal("expected mux to be initialized")
				}
			},
		},
		{
			name: "Custom CORS configuration",
			opts: []Option{
				WithCORS(func() *cors.Cors {
					return cors.New(cors.Options{
						AllowedOrigins: []string{"http://example.com"},
					})
				}(),
				),
			},
			validate: func(t *testing.T, s *Server) {
				if s.cors == nil {
					t.Fatal("expected CORS to be initialized")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServer(tt.opts...)
			tt.validate(t, s)
		})
	}
}

func TestServer_Start_Errors(t *testing.T) {
	tests := []struct {
		name    string
		port    string
		wantErr bool
	}{
		{
			name:    "Invalid port number",
			port:    "-1",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServer(WithPort(tt.port))
			err := s.Start()
			if (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
