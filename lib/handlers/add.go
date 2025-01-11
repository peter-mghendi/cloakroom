package handlers

import (
	"cloakroom/lib"
	"context"
	"fmt"
	"github.com/vbauerster/mpb/v8"
)

// Add adds a plugin to the manifest and optionally downloads it if --fetch is true.
func Add(manifest *lib.Manifest, plugin lib.Plugin, key string, wardrobe string, fetch bool) error {
	if _, exists := manifest.Plugins[key]; exists {
		return fmt.Errorf("plugin %s already exists in the manifest", plugin.Artifact)
	}

	manifest.Plugins[key] = plugin
	fmt.Printf("[INFO] Added plugin to manifest: %s (release: %s, artifact: %s)\n",
		plugin.Artifact, plugin.Tag, plugin.Artifact)

	if fetch {
		ctx := context.Background()
		progress := mpb.New()

		return get(ctx, manifest.Host, wardrobe, key, plugin, false, progress)
	}

	return nil
}
