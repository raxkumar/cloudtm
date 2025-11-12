package helper

import (
	"os"
	"path/filepath"
	"strings"
)

// CopyDirectory copies a directory recursively, excluding specified paths and file patterns
func CopyDirectory(src, dst string, excludeDirs []string, excludeFiles []string, excludePatterns []string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// Skip excluded directories
		for _, exclude := range excludeDirs {
			// Check if this is the exact directory or a path within it
			if relPath == exclude || strings.HasPrefix(relPath, exclude+string(filepath.Separator)) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Skip excluded files (exact match)
		if !info.IsDir() {
			fileName := filepath.Base(path)
			for _, excludeFile := range excludeFiles {
				if fileName == excludeFile {
					return nil
				}
			}

			// Skip files matching patterns (e.g., *.log, *.tmp)
			for _, pattern := range excludePatterns {
				matched, err := filepath.Match(pattern, fileName)
				if err == nil && matched {
					return nil
				}
			}
		}

		// Create destination path
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		return CopyFile(path, dstPath)
	})
}
