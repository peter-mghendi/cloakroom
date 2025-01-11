package handlers

import (
	"cloakroom/lib"
	"cloakroom/lib/utility"
	"context"
	"errors"
	"fmt"
	"github.com/vbauerster/mpb/v8"
	"os"
	"path/filepath"
	"strings"
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
			errs <- get(ctx, manifest.Host, wardrobe, key, plugin, force, progress)
		}(key, plugin)
	}

	group.Wait()
	close(errs)

	if len(errs) > 0 {
		return <-errs
	}
	return nil
}

func get(ctx context.Context, host string, wardrobe string, key string, plugin lib.Plugin, force bool, progress *mpb.Progress) error {
	source, err := url(host, key, plugin)
	if err != nil {
		fmt.Printf("[ERROR] building download URL for %s: %v\n", key, err)
	}

	destination := filepath.Join(wardrobe, plugin.Artifact)
	if _, err := os.Stat(destination); err == nil {
		if force {
			fmt.Printf("[INFO] Removing existing file: %s\n", destination)
			if err := os.RemoveAll(destination); err != nil {
				return fmt.Errorf("failed to remove existing file %s: %w", destination, err)
			}
		} else {

		}

	}

	if err := utility.Download(ctx, progress, source, destination, plugin.Hash, 3); err != nil {
		return fmt.Errorf("downloading %s -> %s: %w", key, destination, err)
	}

	fmt.Printf("[OK] Downloaded %s -> %s\n", key, destination)
	return nil
}

// url returns a direct link for the JAR file from GitHub releases.
func url(host string, key string, plugin lib.Plugin) (string, error) {
	parts := strings.Split(key, "/")
	if len(parts) != 2 {
		return "", errors.New("key must be in form owner/repo")
	}
	owner := parts[0]
	repo := parts[1]

	base := fmt.Sprintf("https://%s/%s/%s/releases", host, owner, repo)

	// e.g. https://github.com/owner/repo/releases/download/v1.2.0/my-plugin.jar
	return fmt.Sprintf("%s/download/%s/%s", base, plugin.Tag, plugin.Artifact), nil
}
