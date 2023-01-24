package main

import (
	"os"
	"strings"
)

// generateFailedPath inserts an indicator in to originalPath to indicate failure.
// The intention is that moving the file located at originalPath to the path returned by generateFailedPath would indicate
// to a user that processing for that file has failed without needing to inspect the file or command output.
func generateFailedPath(originalPath string) string {
	return appendToFileName(originalPath, ".failed")
}

// appendtoFilename appends a string suffix to the filename, before any extension, defined by a dot.
// For example, appending suffix "-suffix" to path "filename.ext" will result in "filename-suffix.ext".
func appendToFileName(path string, suffix string) string {
	originalExt := fullExt(path)
	insertionIndex := len(path) - len(originalExt)
	return strings.Join([]string{
		path[:insertionIndex],
		suffix,
		path[insertionIndex:],
	}, "")
}

// fullExt returns the full file name extension used by the path.
// The full extension is the suffix beginning after the first dot in the final element of the path.
// An empty string is returned if there is no dot.
func fullExt(path string) string {
	extIndex := -1
	for i := len(path) - 1; i >= 0 && !os.IsPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			extIndex = i
		}
	}
	if extIndex != -1 {
		return path[extIndex:]
	} else {
		return ""
	}
}
