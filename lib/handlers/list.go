package handlers

import (
	"cloakroom/lib"
	"fmt"
)

// List outputs all plugins defined in the manifest.
func List(manifest *lib.Manifest) error {
	plugins := manifest.Plugins
	if len(plugins) == 0 {
		fmt.Println("[INFO] No plugins defined in the manifest.")
		return nil
	}

	fmt.Println("[INFO] Plugins in the manifest:")
	for repoKey, plugin := range plugins {
		fmt.Printf("  * %s\n", repoKey)
		fmt.Printf("    - tag:      %s\n", plugin.Tag)
		fmt.Printf("    - artifact: %s\n", plugin.Artifact)
		if plugin.Hash != nil {
			fmt.Printf("    - hash:     %s\n", *plugin.Hash)
		}
		fmt.Println()
	}

	return nil
}
