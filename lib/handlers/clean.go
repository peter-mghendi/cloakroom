package handlers

import (
	"cloakroom/lib/utility"
	"fmt"
)

// Clean removes all contents from the wardrobe directory.
// It ensures the directory is emptied before restoring or other operations.
func Clean(wardrobe string) error {
	fmt.Printf("[INFO] Cleaning wardrobe directory: %s\n", wardrobe)

	if err := utility.Clean(wardrobe); err != nil {
		return fmt.Errorf("failed to clean wardrobe directory: %w", err)
	}

	fmt.Println("[INFO] Wardrobe directory cleaned successfully.")
	return nil
}
