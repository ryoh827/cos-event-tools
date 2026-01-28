package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("usage: %s [root-directory]", filepath.Base(os.Args[0]))
	}

	root := "."
	if len(args) == 1 {
		root = args[0]
	}

	renamed, err := renameEditedFiles(root)
	if err != nil {
		return err
	}

	fmt.Printf("renamed %d file(s)\n", renamed)
	return nil
}

func renameEditedFiles(root string) (int, error) {
	info, err := os.Stat(root)
	if err != nil {
		return 0, fmt.Errorf("stat %s: %w", root, err)
	}
	if !info.IsDir() {
		return 0, fmt.Errorf("%s is not a directory", root)
	}

	renamed := 0
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() {
			return nil
		}

		newName, ok := trimmedEditName(d.Name())
		if !ok {
			return nil
		}

		dir := filepath.Dir(path)
		targetPath := filepath.Join(dir, newName)
		if _, err := os.Stat(targetPath); err == nil {
			return fmt.Errorf("target already exists: %s", targetPath)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("stat %s: %w", targetPath, err)
		}

		if err := os.Rename(path, targetPath); err != nil {
			return fmt.Errorf("rename %s: %w", path, err)
		}

		fmt.Printf("renamed %s -> %s\n", path, targetPath)
		renamed++
		return nil
	})
	if err != nil {
		return renamed, err
	}

	return renamed, nil
}

func trimmedEditName(name string) (string, bool) {
	idx := strings.LastIndex(name, "-Edit")
	if idx == -1 {
		return name, false
	}

	suffix := name[idx+len("-Edit"):]
	if suffix != "" && !strings.HasPrefix(suffix, ".") {
		return name, false
	}

	return name[:idx] + suffix, true
}
