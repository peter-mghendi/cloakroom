package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/vbauerster/mpb/v8"
	"path/filepath"
	"strings"
	"sync"
)

// restore iterates through each plugin, cloning or pulling from Git
func restore(manifest *Manifest) error {
	ctx := context.Background()
	progress := mpb.New()

	var group sync.WaitGroup
	errs := make(chan error, len(manifest.Plugins))

	for key, plugin := range manifest.Plugins {
		group.Add(1)
		go func(rk string, plugin Plugin) {
			defer group.Done()

			// Build the final download URL from our schema.
			downloadURL, err := url(manifest.Host, rk, plugin.Release, plugin.Artifact)
			if err != nil {
				fmt.Printf("[ERROR] building download URL for %s: %v\n", rk, err)
				return
			}

			// Destination file is plugin.Path + "/" + plugin.Artifact
			destFile := filepath.Join(plugin.Path, plugin.Artifact)

			// Attempt the "advanced" download with concurrency, partial file writes, retries, etc.
			if err := download(ctx, progress, downloadURL, destFile, plugin.Hash, 3); err != nil {
				fmt.Printf("[ERROR] downloading %s -> %s: %v\n", rk, destFile, err)
			} else {
				fmt.Printf("[OK] Downloaded %s -> %s\n", rk, destFile)
			}
		}(key, plugin)
	}

	// Wait for all goroutines to finish
	group.Wait()
	close(errs)

	// If any errs occurred, return the first one (or aggregate them)
	if len(errs) > 0 {
		return <-errs
	}
	return nil
}

// url returns a direct link for the JAR file from GitHub releases.
func url(host, repoKey, releaseTag, artifact string) (string, error) {
	parts := strings.Split(repoKey, "/")
	if len(parts) < 2 {
		return "", errors.New("repoKey must be in form owner/repo")
	}
	owner := parts[0]
	repo := parts[1]

	base := fmt.Sprintf("https://%s/%s/%s/releases", host, owner, repo)

	if releaseTag == "latest" {
		// e.g. https://github.com/owner/repo/releases/latest/download/my-plugin.jar
		return fmt.Sprintf("%s/latest/download/%s", base, artifact), nil
	}
	// e.g. https://github.com/owner/repo/releases/download/v1.2.0/my-plugin.jar
	return fmt.Sprintf("%s/download/%s/%s", base, releaseTag, artifact), nil
}
