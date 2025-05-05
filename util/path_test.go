package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimPathByFirstLetters(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple nested path",
			input:    "asdf/bababa/cc/file.yaml",
			expected: "a/b/c/file.yaml",
		},
		{
			name:     "path with leading slash",
			input:    "/var/log/nginx/access.log",
			expected: "/v/l/n/access.log",
		},
		{
			name:     "path with leading dot and slash",
			input:    "./config/test/user/data.json",
			expected: "./c/t/u/data.json",
		},
		{
			name:     "path with leading dot only",
			input:    ".hidden/folder/file",
			expected: ".h/f/file",
		},
		{
			name:     "single file",
			input:    "file.txt",
			expected: "file.txt",
		},
		{
			name:     "empty path",
			input:    "",
			expected: "",
		},
		{
			name:     "just slashes",
			input:    "/",
			expected: "/",
		},
		{
			name:     "path with empty directories",
			input:    "a//b///c/file.txt",
			expected: "a/b/c/file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ShortenPath(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
