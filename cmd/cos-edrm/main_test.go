package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRenameEditedFiles(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	files := []string{
		"IMG123-Edit.CR2",
		"nested/IMG999-Edit.JPG",
		"nested/IMG888-Edit",
		"nested/keep.txt",
		"nested/keep-Editable.txt",
	}

	for _, name := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", path, err)
		}
		if err := os.WriteFile(path, []byte("dummy"), 0o644); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}

	renamed, err := renameEditedFiles(root)
	if err != nil {
		t.Fatalf("renameEditedFiles returned error: %v", err)
	}

	if renamed != 3 {
		t.Fatalf("want 3 renamed files, got %d", renamed)
	}

	renamedFiles := []string{
		"IMG123.CR2",
		"nested/IMG999.JPG",
		"nested/IMG888",
	}
	for _, file := range renamedFiles {
		if _, err := os.Stat(filepath.Join(root, file)); err != nil {
			t.Fatalf("expected %s to exist, stat err: %v", file, err)
		}
	}

	originalNames := []string{
		"IMG123-Edit.CR2",
		"nested/IMG999-Edit.JPG",
		"nested/IMG888-Edit",
	}
	for _, file := range originalNames {
		if _, err := os.Stat(filepath.Join(root, file)); !os.IsNotExist(err) {
			t.Fatalf("expected %s to be renamed, stat err: %v", file, err)
		}
	}

	keptFiles := []string{
		"nested/keep.txt",
		"nested/keep-Editable.txt",
	}
	for _, file := range keptFiles {
		if _, err := os.Stat(filepath.Join(root, file)); err != nil {
			t.Fatalf("expected %s to remain, stat err: %v", file, err)
		}
	}
}

func TestRenameEditedFilesErrors(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	if _, err := renameEditedFiles(filepath.Join(root, "missing")); err == nil {
		t.Fatal("expected error for missing directory")
	}

	filePath := filepath.Join(root, "sample-Edit.jpg")
	if err := os.WriteFile(filePath, []byte("dummy"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	if _, err := renameEditedFiles(filePath); err == nil {
		t.Fatal("expected error when root is not a directory")
	}
}

func TestRenameEditedFilesConflict(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	original := filepath.Join(root, "photo.jpg")
	edited := filepath.Join(root, "photo-Edit.jpg")

	if err := os.WriteFile(original, []byte("orig"), 0o644); err != nil {
		t.Fatalf("write original: %v", err)
	}
	if err := os.WriteFile(edited, []byte("edit"), 0o644); err != nil {
		t.Fatalf("write edited: %v", err)
	}

	if _, err := renameEditedFiles(root); err == nil {
		t.Fatal("expected conflict error when target exists")
	}

	if _, err := os.Stat(edited); err != nil {
		t.Fatalf("expected %s to remain after error: %v", edited, err)
	}
}
