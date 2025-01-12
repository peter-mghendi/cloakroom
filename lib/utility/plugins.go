package utility

import (
	"cloakroom/lib"
	"context"
	"fmt"
	"github.com/vbauerster/mpb/v8"
	"os"
	"path/filepath"
)

// Restore downloads a specified plugin from the given host to the local wardrobe directory.
// If a file already exists and force is false, it skips downloading. If force is true, it overwrites.
// The plugin's hash (if provided) is used for optional verification.
func Restore(
	ctx context.Context,
	host string,
	wardrobe string,
	key string,
	plugin lib.Plugin,
	force bool,
	progress *mpb.Progress,
) error {
	source := fmt.Sprintf("https://%s/%s/releases/download/%s/%s", host, key, plugin.Tag, plugin.Artifact)
	destination := filepath.Join(wardrobe, plugin.Artifact)

	if _, err := os.Stat(destination); err == nil {
		if force {
			fmt.Printf("[INFO] Removing existing file: %s\n", destination)
			if err := os.RemoveAll(destination); err != nil {
				return fmt.Errorf("failed to remove existing file %s: %w", destination, err)
			}
		} else {
			fmt.Printf("[SKIP] Plugin already exists: %s (use --force to overwrite)\n", destination)
			return nil
		}
	}

	if err := Download(ctx, progress, source, destination, plugin.Hash, 3); err != nil {
		return fmt.Errorf("downloading %s -> %s: %w", key, destination, err)
	}

	fmt.Printf("[OK] Downloaded %s -> %s\n", key, destination)
	return nil
}
