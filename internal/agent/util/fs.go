package util

import (
	"bufio"
	"os"
	"path/filepath"
)

func ReadLines(path string, offset int, limit int) ([]string, bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, false, err
	}
	defer file.Close()

	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 2000
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		idx++
		if idx < offset {
			continue
		}
		if idx >= offset && len(lines) < limit {
			lines = append(lines, scanner.Text())
		} else if len(lines) >= limit {
			return lines, true, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, false, err
	}
	return lines, false, nil
}

func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

func ResolvePath(baseDir, rel string) (string, error) {
	if filepath.IsAbs(rel) {
		return filepath.Clean(rel), nil
	}
	return filepath.Clean(filepath.Join(baseDir, rel)), nil
}
