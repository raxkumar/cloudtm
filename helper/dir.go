package helper

import (
	"os"
	"path/filepath"
	"strings"
)

// CopyDirectory copies a directory recursively, excluding specified paths
func CopyDirectory(src, dst string, excludeDirs []string) error {
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
			if strings.HasPrefix(relPath, exclude) || relPath == exclude {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
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
