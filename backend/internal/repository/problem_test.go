package repository

import (
	"embed"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/maeshinshin/sqlearn/backend/internal/domain"
)

//go:embed testdata/*.yaml
var testFS embed.FS

type errorRootFS struct{}

func (errorRootFS) Open(name string) (fs.File, error) {
	return nil, fs.ErrNotExist
}

type errorFileReadFS struct {
	fs.FS
}

func (e errorFileReadFS) ReadFile(name string) ([]byte, error) {
	return nil, fs.ErrInvalid
}

func (e errorFileReadFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(e.FS, name)
}

func TestNewProblemRepository(t *testing.T) {
	subFS, err := fs.Sub(testFS, "testdata")
	if err != nil {
		t.Fatalf("failed to create sub fs: %v", err)
	}

	tests := []struct {
		name    string
		fs      fs.FS
		wantErr bool
		check   func(t *testing.T, repo *ProblemRepository)
	}{
		{
			name:    "Success: Multiple YAML files should be cached",
			fs:      subFS,
			wantErr: false,
			check: func(t *testing.T, repo *ProblemRepository) {
				if len(repo.cache) != 2 {
					t.Errorf("expected cache size 2, got %d", len(repo.cache))
				}

				prob1, ok := repo.cache[1]
				if !ok {
					t.Fatal("expected problem ID 1 to be cached")
				}
				if prob1.Title != "アクティブユーザーの抽出" || prob1.IsOrderMatters != false {
					t.Errorf("problem 1 data mismatch")
				}

				prob2, ok := repo.cache[2]
				if !ok {
					t.Fatal("expected problem ID 2 to be cached")
				}
				if prob2.Title != "高スコア順の並び替え" || prob2.IsOrderMatters != true {
					t.Errorf("problem 2 data mismatch")
				}
			},
		},
		{
			name: "Success: Ignores non-yaml files and directories",
			fs: fstest.MapFS{
				"valid.yaml": &fstest.MapFile{
					Data: []byte("id: 10\ntitle: 'Valid'"),
				},
				"readme.txt": &fstest.MapFile{
					Data: []byte("this should be ignored"),
				},
				"subdir": &fstest.MapFile{
					Mode: fs.ModeDir,
				},
			},
			wantErr: false,
			check: func(t *testing.T, repo *ProblemRepository) {
				if len(repo.cache) != 1 {
					t.Errorf("expected cache size 1, got %d", len(repo.cache))
				}
			},
		},
		{
			name:    "Error: fs.ReadDir fails",
			fs:      errorRootFS{},
			wantErr: true,
			check:   nil,
		},
		{
			name:    "Error: fs.ReadFile fails",
			fs:      errorFileReadFS{FS: subFS},
			wantErr: true,
			check:   nil,
		},
		{
			name: "Error: yaml.Unmarshal fails",
			fs: fstest.MapFS{
				"invalid.yaml": &fstest.MapFile{
					Data: []byte("invalid_yaml: : yaml: string"),
				},
			},
			wantErr: true,
			check:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := NewProblemRepository(tt.fs)

			// エラーの期待値が一致しているか確認
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProblemRepository() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// エラーが期待通り発生した（あるいは発生しなかった）後、
			// 追加の検証ロジック(check)が定義されていれば実行する
			if tt.check != nil && !tt.wantErr {
				tt.check(t, repo)
			}
		})
	}
}

func TestProblemRepository_GetByID(t *testing.T) {
	repo := &ProblemRepository{
		cache: map[int32]*domain.Problem{
			1: {ID: 1, Title: "Problem 1"},
			2: {ID: 2, Title: "Problem 2"},
		},
	}

	tests := []struct {
		name    string
		id      int32
		want    string
		wantErr bool
	}{
		{
			name:    "Returns problem 1 when ID 1 is provided",
			id:      1,
			want:    "Problem 1",
			wantErr: false,
		},
		{
			name:    "Returns problem 2 when ID 2 is provided",
			id:      2,
			want:    "Problem 2",
			wantErr: false,
		},
		{
			name:    "Returns an error when an invalid ID is provided",
			id:      999,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetByID(tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got.Title != tt.want {
				t.Errorf("GetByID() got = %v, want %v", got.Title, tt.want)
			}
		})
	}
}
