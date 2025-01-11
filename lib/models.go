package lib

// Manifest represents the top-level structure of the cloakroom config
type Manifest struct {
	Version string            `json:"version"`
	Host    string            `json:"host"`
	Plugins map[string]Plugin `json:"plugins"`
}

// Plugin represents the config for each plugin denoted by a "user/repo" key
type Plugin struct {
	Name     string `json:"name"`
	Release  string `json:"release"`
	Artifact string `json:"artifact"`
	Hash     string `json:"hash"`
}
