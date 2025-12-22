package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCreatesZipForDirectories(t *testing.T) {
	root := t.TempDir()

	dir := filepath.Join(root, "event1")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatalf("failed to create event dir: %v", err)
	}

	filePath := filepath.Join(dir, "photo.txt")
	if err := os.WriteFile(filePath, []byte("hello"), 0o644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	var out bytes.Buffer
	if err := run(root, &out, strings.NewReader(""), false); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	zipPath := dir + ".zip"
	if _, err := os.Stat(zipPath); err != nil {
		t.Fatalf("zip was not created: %v", err)
	}

	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatalf("failed to open zip: %v", err)
	}
	t.Cleanup(func() { reader.Close() })

	names := map[string]bool{}
	for _, f := range reader.File {
		names[f.Name] = true
	}

	if !names["photo.txt"] {
		t.Fatalf("zip missing file: %v", names)
	}

	if got := out.String(); got == "" {
		t.Fatalf("expected output message, got empty string")
	}
}

func TestRunOverwritesWithForce(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "event1")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatalf("failed to create event dir: %v", err)
	}

	filePath := filepath.Join(dir, "photo.txt")
	if err := os.WriteFile(filePath, []byte("old"), 0o644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	existingZip := dir + ".zip"
	if err := zipDirectory(dir, existingZip); err != nil {
		t.Fatalf("failed to create initial zip: %v", err)
	}

	if err := os.WriteFile(filePath, []byte("new"), 0o644); err != nil {
		t.Fatalf("failed to update file: %v", err)
	}

	var out bytes.Buffer
	if err := run(root, &out, strings.NewReader(""), true); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	content, err := readZipFile(existingZip, "photo.txt")
	if err != nil {
		t.Fatalf("failed to read zip: %v", err)
	}

	if string(content) != "new" {
		t.Fatalf("expected overwritten content, got %q", content)
	}

	if !bytes.Contains(out.Bytes(), []byte("overwrote")) {
		t.Fatalf("expected overwrite message, got %q", out.String())
	}
}

func TestRunDefaultPromptDecline(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "event1")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatalf("failed to create event dir: %v", err)
	}

	filePath := filepath.Join(dir, "photo.txt")
	if err := os.WriteFile(filePath, []byte("old"), 0o644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	existingZip := dir + ".zip"
	if err := zipDirectory(dir, existingZip); err != nil {
		t.Fatalf("failed to create initial zip: %v", err)
	}

	if err := os.WriteFile(filePath, []byte("new"), 0o644); err != nil {
		t.Fatalf("failed to update file: %v", err)
	}

	var out bytes.Buffer
	if err := run(root, &out, strings.NewReader("n\n"), false); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	content, err := readZipFile(existingZip, "photo.txt")
	if err != nil {
		t.Fatalf("failed to read zip: %v", err)
	}

	if string(content) != "old" {
		t.Fatalf("expected original content, got %q", content)
	}

	if !bytes.Contains(out.Bytes(), []byte("user declined")) {
		t.Fatalf("expected decline message, got %q", out.String())
	}
}

func TestRunDefaultPromptAccept(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "event1")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatalf("failed to create event dir: %v", err)
	}

	filePath := filepath.Join(dir, "photo.txt")
	if err := os.WriteFile(filePath, []byte("old"), 0o644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	existingZip := dir + ".zip"
	if err := zipDirectory(dir, existingZip); err != nil {
		t.Fatalf("failed to create initial zip: %v", err)
	}

	if err := os.WriteFile(filePath, []byte("updated"), 0o644); err != nil {
		t.Fatalf("failed to update file: %v", err)
	}

	var out bytes.Buffer
	if err := run(root, &out, strings.NewReader("y\n"), false); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	content, err := readZipFile(existingZip, "photo.txt")
	if err != nil {
		t.Fatalf("failed to read zip: %v", err)
	}

	if string(content) != "updated" {
		t.Fatalf("expected updated content, got %q", content)
	}

	if !bytes.Contains(out.Bytes(), []byte("overwrote")) {
		t.Fatalf("expected overwrite message, got %q", out.String())
	}
}

func readZipFile(path, name string) ([]byte, error) {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	for _, f := range reader.File {
		if f.Name != name {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()

		return io.ReadAll(rc)
	}

	return nil, fmt.Errorf("file %s not found in zip", name)
}
