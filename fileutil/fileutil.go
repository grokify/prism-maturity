// Package fileutil provides file path utilities.
package fileutil

import "path/filepath"

// ReplaceExtension replaces the file extension with a new one.
// If the filename has no extension, the new extension is appended.
func ReplaceExtension(filename, newExt string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return filename + newExt
	}
	return filename[:len(filename)-len(ext)] + newExt
}
