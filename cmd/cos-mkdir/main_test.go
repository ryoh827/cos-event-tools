package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveCSVPath_WithArgument(t *testing.T) {
	got, auto, err := resolveCSVPath([]string{"events.csv"}, t.TempDir())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got != "events.csv" {
		t.Fatalf("got %q, want %q", got, "events.csv")
	}
	if auto {
		t.Fatalf("expected auto=false for explicit argument")
	}
}

func TestResolveCSVPath_TooManyArguments(t *testing.T) {
	_, _, err := resolveCSVPath([]string{"a.csv", "b.csv"}, t.TempDir())
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != "too many arguments" {
		t.Fatalf("got %q", err.Error())
	}
}

func TestResolveCSVPath_AutoSelectSingleCSV(t *testing.T) {
	dir := t.TempDir()
	if err := touchFile(filepath.Join(dir, "events.csv")); err != nil {
		t.Fatalf("touch file: %v", err)
	}

	got, auto, err := resolveCSVPath(nil, dir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got != "events.csv" {
		t.Fatalf("got %q, want %q", got, "events.csv")
	}
	if !auto {
		t.Fatalf("expected auto=true")
	}
}

func TestResolveCSVPath_AutoSelectUppercaseExt(t *testing.T) {
	dir := t.TempDir()
	if err := touchFile(filepath.Join(dir, "events.CSV")); err != nil {
		t.Fatalf("touch file: %v", err)
	}

	got, auto, err := resolveCSVPath(nil, dir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got != "events.CSV" {
		t.Fatalf("got %q, want %q", got, "events.CSV")
	}
	if !auto {
		t.Fatalf("expected auto=true")
	}
}

func TestResolveCSVPath_NoCSV(t *testing.T) {
	dir := t.TempDir()
	if err := touchFile(filepath.Join(dir, "notes.txt")); err != nil {
		t.Fatalf("touch file: %v", err)
	}

	_, _, err := resolveCSVPath(nil, dir)
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != "csv file is required: no CSV files found in current directory" {
		t.Fatalf("got %q", err.Error())
	}
}

func TestResolveCSVPath_MultipleCSV(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"z.csv", "a.CSV"} {
		if err := touchFile(filepath.Join(dir, name)); err != nil {
			t.Fatalf("touch file %q: %v", name, err)
		}
	}

	_, _, err := resolveCSVPath(nil, dir)
	if err == nil {
		t.Fatalf("expected error")
	}
	want := "csv file is required: multiple CSV files found in current directory: a.CSV, z.csv"
	if err.Error() != want {
		t.Fatalf("got %q, want %q", err.Error(), want)
	}
}

func TestResolveCSVPath_IgnoresDirectoryNamedCSV(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "data.csv"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := touchFile(filepath.Join(dir, "memo.txt")); err != nil {
		t.Fatalf("touch file: %v", err)
	}

	_, _, err := resolveCSVPath(nil, dir)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "no CSV files found in current directory") {
		t.Fatalf("got %q", err.Error())
	}
}

func touchFile(path string) error {
	return os.WriteFile(path, []byte{}, 0o644)
}
