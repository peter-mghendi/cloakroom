package lib

// Manifest represents the top-level structure of the cloakroom manifest
type Manifest struct {
	Version string            `mapstructure:"version"`
	Host    string            `mapstructure:"host"`
	Plugins map[string]Plugin `mapstructure:"plugins"`
}

// Plugin represents the configuration for each plugin denoted by a "user/repo" key
type Plugin struct {
	Tag      string  `mapstructure:"tag"`
	Artifact string  `mapstructure:"artifact"`
	Hash     *string `mapstructure:"hash"`
}
