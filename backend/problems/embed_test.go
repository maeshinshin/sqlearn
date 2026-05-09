package problems

import (
	"io/fs"
	"testing"
)

func TestFS(t *testing.T) {
	entries, err := fs.ReadDir(FS, ".")
	if err != nil {
		t.Fatalf("failed to read embed dir: %v", err)
	}

	foundYaml := false
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			if len(name) > 5 && (name[len(name)-5:] == ".yaml" || name[len(name)-4:] == ".yml") {
				foundYaml = true

				f, err := FS.Open(name)
				if err != nil {
					t.Errorf("failed to open file %s: %v", name, err)
					continue
				}
				f.Close()
			}
		}
	}

	if !foundYaml {
		t.Error("no yaml files found in embed FS")
	}
}
