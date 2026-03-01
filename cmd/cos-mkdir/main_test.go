package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveCSVPath(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		setup    func(t *testing.T, dir string)
		wantPath string
		wantAuto bool
		wantErr  string
	}{
		{
			name:     "with argument",
			args:     []string{"events.csv"},
			wantPath: "events.csv",
			wantAuto: false,
		},
		{
			name:    "too many arguments",
			args:    []string{"a.csv", "b.csv"},
			wantErr: "too many arguments",
		},
		{
			name: "auto select single csv",
			setup: func(t *testing.T, dir string) {
				t.Helper()
				if err := touchFile(filepath.Join(dir, "events.csv")); err != nil {
					t.Fatalf("touch file: %v", err)
				}
			},
			wantPath: "events.csv",
			wantAuto: true,
		},
		{
			name: "auto select uppercase ext",
			setup: func(t *testing.T, dir string) {
				t.Helper()
				if err := touchFile(filepath.Join(dir, "events.CSV")); err != nil {
					t.Fatalf("touch file: %v", err)
				}
			},
			wantPath: "events.CSV",
			wantAuto: true,
		},
		{
			name: "no csv",
			setup: func(t *testing.T, dir string) {
				t.Helper()
				if err := touchFile(filepath.Join(dir, "notes.txt")); err != nil {
					t.Fatalf("touch file: %v", err)
				}
			},
			wantErr: "csv file is required: no CSV files found in current directory",
		},
		{
			name: "multiple csv",
			setup: func(t *testing.T, dir string) {
				t.Helper()
				for _, name := range []string{"z.csv", "a.CSV"} {
					if err := touchFile(filepath.Join(dir, name)); err != nil {
						t.Fatalf("touch file %q: %v", name, err)
					}
				}
			},
			wantErr: "csv file is required: multiple CSV files found in current directory: a.CSV, z.csv",
		},
		{
			name: "ignores directory named csv",
			setup: func(t *testing.T, dir string) {
				t.Helper()
				if err := os.Mkdir(filepath.Join(dir, "data.csv"), 0o755); err != nil {
					t.Fatalf("mkdir: %v", err)
				}
				if err := touchFile(filepath.Join(dir, "memo.txt")); err != nil {
					t.Fatalf("touch file: %v", err)
				}
			},
			wantErr: "csv file is required: no CSV files found in current directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if tt.setup != nil {
				tt.setup(t, dir)
			}

			gotPath, gotAuto, err := resolveCSVPath(tt.args, dir)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("got error %q, want %q", err.Error(), tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if gotPath != tt.wantPath {
				t.Fatalf("got path %q, want %q", gotPath, tt.wantPath)
			}
			if gotAuto != tt.wantAuto {
				t.Fatalf("got auto %v, want %v", gotAuto, tt.wantAuto)
			}
		})
	}
}

func touchFile(path string) error {
	return os.WriteFile(path, []byte{}, 0o644)
}
