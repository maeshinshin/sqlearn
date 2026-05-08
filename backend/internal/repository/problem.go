package repository

import (
	"fmt"
	"io/fs"
	"log/slog"
	"strings"

	yaml "github.com/yaml/go-yaml"

	"github.com/maeshinshin/sqlearn/backend/internal/domain"
)

type ProblemRepository struct {
	cache map[int32]*domain.Problem
}

func NewProblemRepository(fileSystem fs.FS) (*ProblemRepository, error) {
	cache := make(map[int32]*domain.Problem)

	entries, err := fs.ReadDir(fileSystem, ".")
	if err != nil {
		slog.Error("failed to read problems directory", "error", err)
		return nil, fmt.Errorf("failed to read problems directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		fileBytes, err := fs.ReadFile(fileSystem, entry.Name())
		if err != nil {
			slog.Error("failed to read problem file", "file", entry.Name(), "error", err)
			return nil, fmt.Errorf("failed to read problem file %s: %w", entry.Name(), err)
		}

		var prob domain.Problem
		if err := yaml.Unmarshal(fileBytes, &prob); err != nil {
			slog.Error("failed to unmarshal problem YAML", "file", entry.Name(), "error", err)
			return nil, fmt.Errorf("failed to unmarshal problem YAML from file %s: %w", entry.Name(), err)
		}

		cache[prob.ID] = &prob
	}

	return &ProblemRepository{cache: cache}, nil
}

func (r *ProblemRepository) GetByID(id int32) (*domain.Problem, error) {
	prob, ok := r.cache[id]
	if !ok {
		return nil, fmt.Errorf("problem with ID %d not found", id)
	}
	return prob, nil
}
