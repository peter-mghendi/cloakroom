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

			source, err := url(manifest.Host, key, plugin.Release, plugin.Artifact)
			if err != nil {
				fmt.Printf("[ERROR] building download URL for %s: %v\n", key, err)
				return
			}

			destination := filepath.Join(wardrobe, plugin.Artifact)
			if _, err := os.Stat(destination); err == nil {
				if force {
					fmt.Printf("[INFO] Removing existing file: %s\n", destination)
					if err := os.RemoveAll(destination); err != nil {
						errs <- fmt.Errorf("failed to remove existing file %s: %w", destination, err)
						return
					}
				} else {
					fmt.Printf("[SKIP] Plugin already exists: %s (use --force to overwrite)\n", destination)
					return
				}
			}

			if err := utility.Download(ctx, progress, source, destination, plugin.Hash, 3); err != nil {
				fmt.Printf("[ERROR] downloading %s -> %s: %v\n", key, destination, err)
			} else {
				fmt.Printf("[OK] Downloaded %s -> %s\n", key, destination)
			}
		}(key, plugin)
	}

	group.Wait()
	close(errs)

	if len(errs) > 0 {
		return <-errs
	}
	return nil
}

// url returns a direct link for the JAR file from GitHub releases.
func url(host, key, tag, artifact string) (string, error) {
	parts := strings.Split(key, "/")
	if len(parts) != 2 {
		return "", errors.New("key must be in form owner/repo")
	}
	owner := parts[0]
	repo := parts[1]

	base := fmt.Sprintf("https://%s/%s/%s/releases", host, owner, repo)

	if tag == "latest" {
		// e.g. https://github.com/owner/repo/releases/latest/download/my-plugin.jar
		return fmt.Sprintf("%s/latest/download/%s", base, artifact), nil
	}
	// e.g. https://github.com/owner/repo/releases/download/v1.2.0/my-plugin.jar
	return fmt.Sprintf("%s/download/%s/%s", base, tag, artifact), nil
}
