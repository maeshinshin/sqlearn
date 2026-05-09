package handler

import (
	"context"
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"

	"connectrpc.com/connect"

	pb "github.com/maeshinshin/sqlearn/backend/gen/problem/v1"
)

func TestNewProblemHandler(t *testing.T) {
	validFS := fstest.MapFS{
		"1.yaml": &fstest.MapFile{Data: []byte("id: 1\nanswer_sql: 'SELECT * FROM users;'")},
	}
	invalidFS := fstest.MapFS{
		"invalid.yaml": &fstest.MapFile{Data: []byte("invalid_yaml: : string")},
	}

	tests := []struct {
		name    string
		fs      fs.FS
		wantErr bool
	}{
		{
			name:    "Success: initialized with valid FS",
			fs:      validFS,
			wantErr: false,
		},
		{
			name:    "Error: fails to initialize repository",
			fs:      invalidFS,
			wantErr: true,
		},
		{
			name:    "Fallback: fileSystem is nil",
			fs:      nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewProblemHandler(tt.fs)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProblemHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProblemHandler_GetAnswer(t *testing.T) {
	validFS := fstest.MapFS{
		"1.yaml": &fstest.MapFile{Data: []byte("id: 1\nanswer_sql: 'SELECT * FROM users;'")},
	}
	h, err := NewProblemHandler(validFS)
	if err != nil {
		t.Fatalf("failed to setup handler: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name    string
		req     *pb.GetAnswerRequest
		wantErr bool
		check   func(t *testing.T, res *pb.GetAnswerResponse, err error)
	}{
		{
			name:    "Success: returns answer for valid ID",
			req:     &pb.GetAnswerRequest{Id: 1},
			wantErr: false,
			check: func(t *testing.T, res *pb.GetAnswerResponse, err error) {
				if res == nil {
					t.Fatal("expected response, got nil")
				}
				if res.AnswerSql != "SELECT * FROM users;" {
					t.Errorf("expected 'SELECT * FROM users;', got '%s'", res.AnswerSql)
				}
			},
		},
		{
			name:    "Error: returns CodeNotFound for invalid ID",
			req:     &pb.GetAnswerRequest{Id: 999},
			wantErr: true,
			check: func(t *testing.T, res *pb.GetAnswerResponse, err error) {
				if res != nil {
					t.Errorf("expected nil response, got %v", res)
				}

				var connectErr *connect.Error
				if ok := errors.As(err, &connectErr); ok {
					if connectErr.Code() != connect.CodeNotFound {
						t.Errorf("expected CodeNotFound, got %v", connectErr.Code())
					}
				} else {
					t.Errorf("expected connect.Error, got %T", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := h.GetAnswer(ctx, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAnswer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.check != nil {
				tt.check(t, res, err)
			}
		})
	}
}

func TestProblemHandler_GetProblem(t *testing.T) {
	validFS := fstest.MapFS{
		"1.yaml": &fstest.MapFile{Data: []byte(`id: 1
title: "Test Problem"
description: "This is a test description"
setup_sql: "CREATE TABLE test;"
expected_result_json: '[{"id": 1}]'
is_order_matters: true
`)},
	}
	h, err := NewProblemHandler(validFS)
	if err != nil {
		t.Fatalf("failed to setup handler: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name    string
		req     *pb.GetProblemRequest
		wantErr bool
		check   func(t *testing.T, res *pb.GetProblemResponse, err error)
	}{
		{
			name:    "Success: returns problem for valid ID",
			req:     &pb.GetProblemRequest{Id: 1},
			wantErr: false,
			check: func(t *testing.T, res *pb.GetProblemResponse, err error) {
				if res == nil {
					t.Fatal("expected response, got nil")
				}
				if res.Id != 1 {
					t.Errorf("expected Id 1, got %d", res.Id)
				}
				if res.Title != "Test Problem" {
					t.Errorf("expected Title 'Test Problem', got '%s'", res.Title)
				}
				if res.Description != "This is a test description" {
					t.Errorf("expected Description, got '%s'", res.Description)
				}
				if res.SetupSql != "CREATE TABLE test;" {
					t.Errorf("expected SetupSql, got '%s'", res.SetupSql)
				}
				if res.ExpectedResultJson != `[{"id": 1}]` {
					t.Errorf("expected ExpectedResultJson, got '%s'", res.ExpectedResultJson)
				}
				if res.IsOrderMatters != true {
					t.Errorf("expected IsOrderMatters true, got %v", res.IsOrderMatters)
				}
			},
		},
		{
			name:    "Error: returns CodeNotFound for invalid ID",
			req:     &pb.GetProblemRequest{Id: 999},
			wantErr: true,
			check: func(t *testing.T, res *pb.GetProblemResponse, err error) {
				if res != nil {
					t.Errorf("expected nil response, got %v", res)
				}

				var connectErr *connect.Error
				if ok := errors.As(err, &connectErr); ok {
					if connectErr.Code() != connect.CodeNotFound {
						t.Errorf("expected CodeNotFound, got %v", connectErr.Code())
					}
				} else {
					t.Errorf("expected connect.Error, got %T", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := h.GetProblem(ctx, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetProblem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.check != nil {
				tt.check(t, res, err)
			}
		})
	}
}
