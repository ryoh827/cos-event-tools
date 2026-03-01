package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func main() {
	csvPath, autoSelected, err := resolveCSVPath(os.Args[1:], ".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		printUsage(os.Stderr)
		os.Exit(1)
	}

	if autoSelected {
		fmt.Printf("using csv file: %s\n", csvPath)
	}

	if err := run(csvPath); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage(w io.Writer) {
	fmt.Fprintf(w, "usage: %s [csv-file]\n", filepath.Base(os.Args[0]))
	fmt.Fprintln(w, "if [csv-file] is omitted, exactly one CSV in current directory is used automatically")
}

func resolveCSVPath(args []string, dir string) (string, bool, error) {
	if len(args) > 1 {
		return "", false, fmt.Errorf("too many arguments")
	}
	if len(args) == 1 {
		return args[0], false, nil
	}

	csvFiles, err := findCSVFiles(dir)
	if err != nil {
		return "", false, fmt.Errorf("read current directory: %w", err)
	}
	if len(csvFiles) == 0 {
		return "", false, fmt.Errorf("csv file is required: no CSV files found in current directory")
	}
	if len(csvFiles) > 1 {
		return "", false, fmt.Errorf(
			"csv file is required: multiple CSV files found in current directory: %s",
			strings.Join(csvFiles, ", "),
		)
	}

	return csvFiles[0], true, nil
}

func findCSVFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	csvFiles := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.EqualFold(filepath.Ext(entry.Name()), ".csv") {
			csvFiles = append(csvFiles, entry.Name())
		}
	}

	sort.Strings(csvFiles)
	return csvFiles, nil
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
			continue
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
