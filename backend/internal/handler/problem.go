package handler

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"

	"connectrpc.com/connect"

	pb "github.com/maeshinshin/sqlearn/backend/gen/problem/v1"
	"github.com/maeshinshin/sqlearn/backend/gen/problem/v1/problemv1connect"
	"github.com/maeshinshin/sqlearn/backend/internal/repository"
)

type ProblemHandler struct {
	repo *repository.ProblemRepository
}

var _ problemv1connect.ProblemServiceHandler = (*ProblemHandler)(nil)

func NewProblemHandler(fileSystem fs.FS) (*ProblemHandler, error) {
	slog.Info("Initializing ProblemHandler")
	if fileSystem == nil {
		slog.Warn("No file system provided for ProblemHandler, using os.DirFS(\"../../problems\") as default")
		fileSystem = fs.FS(os.DirFS("../../problems"))
	}

	repo, err := repository.NewProblemRepository(fileSystem)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ProblemRepository: %w", err)
	}

	return &ProblemHandler{repo: repo}, nil
}

func (h *ProblemHandler) GetProblem(
	ctx context.Context,
	req *pb.GetProblemRequest,
) (
	*pb.GetProblemResponse,
	error,
) {
	slog.Info("Received GetProblem request", slog.Int("id", int(req.Id)))

	prob, err := h.repo.GetByID(req.Id)
	if err != nil {
		slog.Error("Failed to get problem by ID", "id", req.Id, "error", err)
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	res := &pb.GetProblemResponse{
		Id:                 prob.ID,
		Title:              prob.Title,
		Description:        prob.Description,
		SetupSql:           prob.SetupSQL,
		ExpectedResultJson: prob.ExpectedResultJSON,
		IsOrderMatters:     prob.IsOrderMatters,
	}
	return res, nil
}

func (h *ProblemHandler) GetAnswer(
	ctx context.Context,
	req *pb.GetAnswerRequest,
) (
	*pb.GetAnswerResponse,
	error,
) {
	slog.Info("Received GetAnswer request", slog.Int("id", int(req.Id)))
	prob, err := h.repo.GetByID(req.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	res := &pb.GetAnswerResponse{
		AnswerSql: prob.AnswerSQL,
	}

	return res, nil
}
