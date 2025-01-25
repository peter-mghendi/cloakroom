package cmd

import (
	"cloakroom/lib"
	"cloakroom/lib/utility"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var manifest string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   utility.Cloakroom,
	Short: "Minimal plugin manager for Keycloak.",
	Long: `Cloakroom is a CLI tool for managing Keycloak plugins across environments.

It allows you to:
- Initialize a manifest for your plugin configuration.
- Add, remove, or restore plugins based on the manifest and lock file.
- Clean plugin directories and ensure consistency across multiple environments.

Examples:
  # Restore plugins from the manifest
  cloakroom restore

  # Add a new plugin to the manifest
  cloakroom add user/repo --release v1.2.0 --artifact plugin.jar

  # Clean plugin directory
  cloakroom clean
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(configure)
	rootCmd.PersistentFlags().StringVar(
		&manifest,
		"manifest",
		"",
		`path to a manifest file (supported formats: .hcl, .ini, .json, .toml, .yaml).
If unspecified, cloakroom will look for the following files in the current directory:
  - cloakroom.hcl
  - cloakroom.ini
  - cloakroom.json
  - cloakroom.toml
  - cloakroom.yaml

Only one manifest is supported at a time—an error occurs if multiple valid files are found.`,
	)
}

// configure reads in a manifest file and ENV variables.
func configure() {
	if manifest != "" {
		viper.SetConfigFile(manifest)
	} else {
		options := detect()
		switch len(options) {
		case 0:
			break
		case 1:
			viper.AddConfigPath(".")
			viper.SetConfigType(options[0])
			viper.SetConfigName(utility.Cloakroom)
			viper.SetConfigFile(filename(options[0]))
			break
		default:
			message := "[ERROR] Found multiple manifests:\n%s\n\nCloakroom does not support multiple manifests at once."
			_, _ = fmt.Printf(message, strings.Join(utility.Map(options, filename), "\n- "))
		}
	}

	viper.SetEnvPrefix(utility.Cloakroom)
	viper.AutomaticEnv()

	viper.SetDefault("plugins", make(map[string]lib.Plugin))

	err := viper.ReadInConfig()
	if err == nil {
		_, _ = fmt.Fprintf(os.Stdout, "[INFO] Using manifest: %s\n", viper.ConfigFileUsed())
		return
	}

	var configFileNotFoundError viper.ConfigFileNotFoundError
	if errors.As(err, &configFileNotFoundError) {
		cmd, _, _ := rootCmd.Find(os.Args[1:])
		if cmd != nil && cmd.Name() == "init" {
			_, _ = fmt.Println("[INFO] No manifest found. This is expected for 'init' command.")
			return
		}

		_, _ = fmt.Println("[ERROR] Manifest not found. Run 'cloakroom init' to create one.")
		os.Exit(1)
	}

	cobra.CheckErr(err)
}

// detect looks for any valid manifest files
func detect() []string {
	formats := []string{"hcl", "ini", "json", "toml", "yaml"}
	detected := make([]string, 0, len(formats))

	for _, format := range formats {
		_, err := os.Stat(filename(format))
		if err == nil {
			detected = append(detected, format)
		}
	}

	return detected
}

// filename returns a manifest filename given a format
func filename(format string) string {
	return fmt.Sprintf("%s.%s", utility.Cloakroom, format)
}
