package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <csv-file>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	if err := run(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(csvPath string) error {
	file, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("open csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	line := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read csv line %d: %w", line+1, err)
		}

		line++
		if line == 1 {
			continue // skip header row
		}

		for i := range record {
			record[i] = strings.TrimSpace(record[i])
		}

		if len(record) < 3 {
			fmt.Fprintf(os.Stderr, "skip line %d: need at least 3 columns\n", line)
			continue
		}

		dateStr, err := parseDate(record[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip line %d: %v\n", line, err)
			continue
		}

		event := sanitizeSegment(record[1])
		name := appendHonorific(record[2])
		if event == "" || name == "" {
			fmt.Fprintf(os.Stderr, "skip line %d: empty event or name\n", line)
			continue
		}

		dirName := fmt.Sprintf("%s_%s_%s_りょう", dateStr, event, name)
		if err := os.Mkdir(dirName, 0o755); err != nil {
			if os.IsExist(err) {
				fmt.Printf("skip %s (already exists)\n", dirName)
				continue
			}
			return fmt.Errorf("create %s: %w", dirName, err)
		}

		fmt.Printf("created %s\n", dirName)
	}

	return nil
}

func parseDate(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("empty date value")
	}

	layouts := []string{
		"2006/1/2",
		"2006-1-2",
		"2006.1.2",
		"20060102",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t.Format("20060102"), nil
		}
	}

	return "", fmt.Errorf("invalid date %q", value)
}

func sanitizeSegment(value string) string {
	value = strings.TrimSpace(value)
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	return replacer.Replace(value)
}

func appendHonorific(value string) string {
	sanitized := sanitizeSegment(value)
	if sanitized == "" {
		return ""
	}
	return sanitized + "さん"
}
