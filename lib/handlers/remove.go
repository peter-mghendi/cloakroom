package handlers

import (
	"cloakroom/lib"
	"cloakroom/lib/utility"
	"fmt"
	"path/filepath"
)

// Remove removes a plugin from the manifest and optionally deletes its files.
func Remove(manifest *lib.Manifest, artifact string, wardrobe string, purge bool) error {
	plugin, exists := manifest.Plugins[artifact]
	if !exists {
		return fmt.Errorf("plugin %s not found in the manifest", artifact)
	}

	delete(manifest.Plugins, artifact)
	fmt.Printf("[INFO] Removed plugin from manifest: %s\n", artifact)

	if purge {
		destination := filepath.Join(wardrobe, plugin.Artifact)
		err := utility.Remove(destination)
		if err != nil {
			return fmt.Errorf("failed to purge plugin files for %s: %w", artifact, err)
		}
		fmt.Printf("[INFO] Successfully purged plugin files: %s\n", destination)
	}

	return nil
}
