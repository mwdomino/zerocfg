package util

import (
	"strings"
)

func ShortenPath(path string) string {
	if path == "" {
		return ""
	}

	// Сохраняем начальный префикс
	prefix := ""
	if strings.HasPrefix(path, "/") {
		prefix = "/"
		path = path[1:]
	} else if strings.HasPrefix(path, "./") {
		prefix = "./"
		path = path[2:]
	} else if strings.HasPrefix(path, ".") {
		prefix = "."
		path = path[1:]
	}

	parts := strings.Split(path, "/")

	file := parts[len(parts)-1]
	dirs := parts[:len(parts)-1]

	var shortDirs []string
	for _, part := range dirs {
		if part != "" {
			shortDirs = append(shortDirs, string(part[0]))
		}
	}

	shortPath := strings.Join(shortDirs, "/")
	if shortPath != "" {
		shortPath = shortPath + "/" + file
	} else {
		shortPath = file
	}

	return prefix + shortPath
}
