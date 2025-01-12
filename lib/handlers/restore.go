package handlers

import (
	"cloakroom/lib"
	"cloakroom/lib/utility"
	"context"
	"fmt"
	"github.com/vbauerster/mpb/v8"
	"sync"
)

// Restore iterates through each plugin, downloading it if it does not exist.
func Restore(manifest *lib.Manifest, wardrobe string, clean bool, force bool) error {
	ctx := context.Background()
	progress := mpb.New()

	if clean {
		fmt.Printf("[INFO] Cleaning wardrobe directory: %s\n", wardrobe)
		if err := utility.Clean(wardrobe); err != nil {
			return fmt.Errorf("failed to clean wardrobe directory: %w", err)
		}
	}

	var group sync.WaitGroup
	errs := make(chan error, len(manifest.Plugins))

	for key, plugin := range manifest.Plugins {
		group.Add(1)
		go func(key string, plugin lib.Plugin) {
			defer group.Done()
			errs <- utility.Restore(ctx, manifest.Host, wardrobe, key, plugin, force, progress)
		}(key, plugin)
	}

	group.Wait()
	close(errs)

	if len(errs) > 0 {
		return <-errs
	}
	return nil
}
