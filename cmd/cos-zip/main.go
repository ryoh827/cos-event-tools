package main

import (
	"archive/zip"
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	force := fs.Bool("f", false, "overwrite existing zip archives")
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "usage: %s [flags] <root-directory>\n", filepath.Base(os.Args[0]))
		fs.PrintDefaults()
	}

	fs.Parse(os.Args[1:])

	args := fs.Args()
	if len(args) != 1 {
		fs.Usage()
		os.Exit(1)
	}

	if err := run(args[0], os.Stdout, os.Stdin, *force); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(root string, out io.Writer, in io.Reader, force bool) error {
	info, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("stat root: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("root %s is not a directory", root)
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return fmt.Errorf("read root: %w", err)
	}

	created := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirPath := filepath.Join(root, entry.Name())
		zipPath := dirPath + ".zip"

		exists := false
		if _, err := os.Stat(zipPath); err == nil {
			exists = true
		} else if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("stat zip for %s: %w", entry.Name(), err)
		}

		if exists && !force {
			ok, err := confirmOverwrite(out, in, zipPath)
			if err != nil {
				return fmt.Errorf("confirm overwrite for %s: %w", entry.Name(), err)
			}
			if !ok {
				fmt.Fprintf(out, "skip %s (user declined)\n", filepath.Base(zipPath))
				continue
			}
		}

		if err := zipDirectory(dirPath, zipPath); err != nil {
			return fmt.Errorf("zip %s: %w", entry.Name(), err)
		}

		action := "created"
		if exists {
			action = "overwrote"
		}
		fmt.Fprintf(out, "%s %s\n", action, filepath.Base(zipPath))
		created++
	}

	if created == 0 {
		fmt.Fprintln(out, "no directories to zip")
	}

	return nil
}

func confirmOverwrite(out io.Writer, in io.Reader, zipPath string) (bool, error) {
	fmt.Fprintf(out, "overwrite %s? [y/N]: ", filepath.Base(zipPath))

	reader := bufio.NewReader(in)
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return false, err
	}

	answer := strings.TrimSpace(strings.ToLower(line))
	return answer == "y" || answer == "yes", nil
}

func zipDirectory(dirPath, zipPath string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("create zip: %w", err)
	}
	defer zipFile.Close()

	writer := zip.NewWriter(zipFile)
	defer writer.Close()

	return filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(rel)

		if d.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		w, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := io.Copy(w, file); err != nil {
			return fmt.Errorf("copy %s: %w", path, err)
		}

		return nil
	})
}
