package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

var defaults = map[string]string{
	"host": "github.com",
	"path": "/opt/keycloak/plugins",
}

// load opens the specified JSON file and decodes it into a Manifest struct
func load(filename string) (*Manifest, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	var manifest Manifest
	if err := json.NewDecoder(file).Decode(&manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

// resolve sets fallback values for host, path, plugin name/path, etc.
func resolve(m *Manifest) {
	if m.Host == "" {
		m.Host = defaults["host"]
	}
	if m.Path == "" {
		m.Path = defaults["path"]
	}

	for key, plugin := range m.Plugins {
		// If plugin.Name is empty, default to the last segment of "user/repo"
		if plugin.Name == "" {
			parts := strings.Split(key, "/")
			plugin.Name = parts[len(parts)-1]
		}

		// If plugin.Path is empty, default to manifest.Path + "/" + plugin.Name
		if plugin.Path == "" {
			plugin.Path = filepath.Join(m.Path, plugin.Name)
		}

		// Update the plugin object back into the map
		m.Plugins[key] = plugin
	}
}
