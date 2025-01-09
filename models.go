package main

// Manifest represents the top-level structure of the cloakroom config
type Manifest struct {
	Version string            `json:"version"`
	Host    string            `json:"host"`
	Path    string            `json:"path"`
	Plugins map[string]Plugin `json:"plugins"`
}

// Plugin represents the config for each plugin denoted by a "user/repo" key
type Plugin struct {
	Name     string `json:"name"`
	Release  string `json:"release"`
	Artifact string `json:"artifact"`
	Path     string `json:"path"`
	Hash     string `json:"hash"`
}
