//go:build !windows
// +build !windows

package util

import (
	"path/filepath"
	"strings"
)

func IsHiddenFile(path string) (bool, error) {
	return strings.HasPrefix(filepath.Base(path), "."), nil
}
