package utility

import (
	"fmt"
	"os"
	"path/filepath"
)

// Clean removes all contents from the provided directory without deleting the directory itself.
// If the directory doesn't exist, it returns an error.
func Clean(dir string) error {
	// Check if the directory exists
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", dir)
	}
	if !info.IsDir() {
		return fmt.Errorf("provided path is not a directory: %s", dir)
	}

	// Read all entries in the directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory contents: %w", err)
	}

	// Remove each entry
	for _, entry := range entries {
		entryPath := filepath.Join(dir, entry.Name())

		// Remove directory contents recursively or files directly
		if entry.IsDir() {
			err = os.RemoveAll(entryPath)
		} else {
			err = os.Remove(entryPath)
		}

		if err != nil {
			return fmt.Errorf("failed to remove %s: %w", entryPath, err)
		}
	}

	return nil
}

// Remove removes a single file from the filesystem.
func Remove(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("[WARN] File does not exist: %s\n", path)
		return nil
	}

	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("failed to delete file: %s, %w", path, err)
	}

	return nil
}
